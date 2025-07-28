package main

import (
	"embed"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
)

// TODO move this to a new package?
//
//go:embed decks
var embedFS embed.FS

//go:embed templates
var templateFS embed.FS

//go:embed css
var cssFS embed.FS

const decksDir string = "decks"

type Deck struct {
	Name      string
	numCards  int // this is redundant, but meh.
	cardNames []string
}

// for login, i need a DB of users and the decks they have access to
// that'll give a JWT or equivalent that contains the access list
// can remove most of the code for running one's own server; it'll run in a controlled server environment.
// on start connect to the db
// make an interceptor to look for auth
// if you have none go to login
// if you do, does it authorize the deck you want?
// OR if you're listing decks, list only the ones you have
// does this mean rewriting in a new framework that makes things easier?

var Decks = make(map[string]*Deck)

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
	if _, err := os.Stat(decksDir); !os.IsNotExist(err) {
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

	Decks = FindDecks(decksFS)

	var deckNames []string
	for _, d := range Decks {
		deckNames = append(deckNames, d.Name)
	}
	log.Println("Known decks:", deckNames)

	r := mux.NewRouter()

	r.PathPrefix("/" + decksDir + "/").Handler(http.FileServer(http.FS(decksFS)))
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
	card := ChooseUndrawnCard(deck, drawnCards)
	
	var imgTag string
	var message string
	
	if card != "" {
		// Mark card as drawn
		drawnCards[card] = true
		session.Values[sessionKey] = drawnCards
		
		log.Printf("Drew card: %s, Total drawn: %d/%d", card, len(drawnCards), deck.numCards)
		
		// Save session
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Error saving session: %v", err)
		}
		
		imgTag = fmt.Sprintf(`<img id="cardImg" src="/decks/%s/%s"/>`, deckName, card)
		
		// Check if deck is exhausted
		if len(drawnCards) >= deck.numCards {
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
		Decks map[string]*Deck
	}{
		Decks: Decks,
	})
	if err != nil {
		log.Println(err)
	}
}

func NewDeck(deckName string) *Deck {
	d := Deck{Name: deckName}
	return &d
}

func FindDecks(rootFS fs.FS) map[string]*Deck {

	decks := make(map[string]*Deck)

	rootEntries, err := fs.ReadDir(rootFS, decksDir)
	if err != nil {
		log.Println(err)
	}
	for _, entry := range rootEntries {
		if entry.IsDir() {
			// Can this be done without loading all the card files? Is that actually happening or only appearing in debug because of the debugger?
			deckEntries := getDeckEntries(rootFS, entry)

			deck := NewDeck(entry.Name())
			deck.numCards = len(deckEntries)
			for _, cardEntry := range deckEntries {
				deck.cardNames = append(deck.cardNames, cardEntry.Name())
			}
			decks[deck.Name] = deck
		}

	}
	return decks
}

func getDeckEntries(rootFS fs.FS, deckEntry fs.DirEntry) []fs.DirEntry {
	deckEntries, err := fs.ReadDir(rootFS, path.Join(decksDir, deckEntry.Name()))
	if err != nil {
		panic(err)
	}
	return deckEntries
}

func ChooseUndrawnCard(deck *Deck, drawnCards map[string]bool) string {
	var availableCards []string
	for _, card := range deck.cardNames {
		if !drawnCards[card] {
			availableCards = append(availableCards, card)
		}
	}
	
	if len(availableCards) == 0 {
		return ""
	}
	
	return availableCards[rand.Intn(len(availableCards))]
}

func ChooseRandomCard(deck *Deck) string {
	return deck.cardNames[rand.Intn(deck.numCards)]
}
