package main

import (
	"embed"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sekullbe/carddeck/pkg/decks"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed decks
var embedFS embed.FS

//go:embed templates
var templateFS embed.FS

//go:embed css
var cssFS embed.FS

// for login, i need a DB of users and the decks they have access to
// that'll give a JWT or equivalent that contains the access list
// can remove most of the code for running one's own server; it'll run in a controlled server environment.
// on start connect to the db
// make an interceptor to look for auth
// if you have none go to login
// if you do, does it authorize the deck you want?
// OR if you're listing decks, list only the ones you have
// does this mean rewriting in a new framework that makes things easier?

var Decks = make(map[string]*decks.Deck)

// Session store for tracking drawn cards
var store = sessions.NewCookieStore([]byte("your-secret-key-change-in-production"))

func init() {
	// Register map[string]bool type with gob for proper serialization
	gob.Register(map[string]bool{})
}

func main() {

	port := flag.String("port", "8888", "default http port")
	useBuiltinDecks := flag.Bool("builtin", false, "use built-in decks")
	flag.Parse()

	localDecksExist := false
	if _, err := os.Stat("decks"); !os.IsNotExist(err) {
		localDecksExist = true
	}

	var decksFS fs.FS
	// If local decks exist and we're not asked to use builtin decks anyway, use the local decks
	if localDecksExist && !*useBuiltinDecks {
		decksFS = os.DirFS(".")
		log.Println("Running with local decks")
	} else {
		log.Println("Running with embedded decks")
		decksFS = embedFS
	}

	Decks = decks.FindDecks(decksFS)

	var deckNames []string
	for _, d := range Decks {
		deckNames = append(deckNames, d.Name)
	}
	log.Println("Known decks:", deckNames)

	r := mux.NewRouter()

	r.PathPrefix("/decks/").Handler(http.FileServer(http.FS(decksFS)))
	r.PathPrefix("/css/").Handler(http.FileServer(http.FS(cssFS)))
	r.HandleFunc("/card/{deck}", serveDeck)
	r.HandleFunc("/shuffle/{deck}", shuffleDeck)
	r.HandleFunc("/", serveTemplate)

	log.Println("Ready at http://localhost:" + *port)
	log.Println(http.ListenAndServe(":"+*port, r))
	//http.Handle("/", r)
}

func serveDeck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deckName := vars["deck"]
	
	if deckName == "" {
		http.Error(w, "Deck name required", http.StatusBadRequest)
		return
	}
	
	deck, exists := Decks[deckName]
	if !exists {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}
	
	// Get or create session
	session, err := store.Get(r, "card-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	
	// Get drawn cards for this deck from session
	sessionKey := fmt.Sprintf("drawn_%s", deckName)
	drawnCardsInterface, exists := session.Values[sessionKey]
	var drawnCards map[string]bool
	
	if exists {
		if drawnMap, ok := drawnCardsInterface.(map[string]bool); ok {
			drawnCards = drawnMap
		} else {
			log.Printf("Session data type assertion failed for key %s", sessionKey)
			drawnCards = make(map[string]bool)
		}
	} else {
		drawnCards = make(map[string]bool)
	}
	
	log.Printf("Deck: %s, Previously drawn cards: %v", deckName, drawnCards)
	
	// Choose a card that hasn't been drawn
	card := decks.ChooseUndrawnCard(deck, drawnCards)
	
	var imgTag string
	var message string
	
	if card != "" {
		// Mark card as drawn
		drawnCards[card] = true
		session.Values[sessionKey] = drawnCards
		
		log.Printf("Drew card: %s, Total drawn: %d/%d", card, len(drawnCards), deck.NumCards)
		
		// Save session
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Error saving session: %v", err)
		}
		
		imgTag = fmt.Sprintf(`<img id="cardImg" src="/decks/%s/%s"/>`, deckName, card)
		
		// Check if deck is exhausted
		if len(drawnCards) >= deck.NumCards {
			message = `<p style="color: orange; font-weight: bold;">All cards drawn! Click shuffle to reset.</p>`
		}
	} else {
		message = `<p style="color: red; font-weight: bold;">All cards have been drawn! Click shuffle to reset the deck.</p>`
	}
	
	response := imgTag + message
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func shuffleDeck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deckName := vars["deck"]
	
	if deckName == "" {
		http.Error(w, "Deck name required", http.StatusBadRequest)
		return
	}
	
	// Get or create session
	session, err := store.Get(r, "card-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	
	// Clear drawn cards for this deck
	sessionKey := fmt.Sprintf("drawn_%s", deckName)
	session.Values[sessionKey] = make(map[string]bool)
	
	log.Printf("Shuffled deck: %s", deckName)
	
	// Save session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Error shuffling deck", http.StatusInternalServerError)
		return
	}
	
	// Return success message
	message := `<p style="color: green; font-weight: bold;">Deck shuffled! All cards are available again.</p>`
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// would prefer to pass this in but a handler has a fixed signature
	tmpl, err := template.ParseFS(templateFS, "templates/*.gohtml")
	if err != nil {
		log.Println(err)
	}
	err = tmpl.ExecuteTemplate(w, "layout", struct {
		Decks map[string]*decks.Deck
	}{
		Decks: Decks,
	})
	if err != nil {
		log.Println(err)
	}
}
