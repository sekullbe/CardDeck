package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
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

	// TODO store the cards a user has seen in a session?

	r := mux.NewRouter()

	r.PathPrefix("/" + decksDir + "/").Handler(http.FileServer(http.FS(decksFS)))
	r.PathPrefix("/css/").Handler(http.FileServer(http.FS(decksFS)))
	r.HandleFunc("/card/{deck}", serveDeck)
	r.HandleFunc("/", serveTemplate)

	log.Println("Ready at http://localhost:" + *port)
	log.Println(http.ListenAndServe(":"+*port, r))
	//http.Handle("/", r)
}

func serveDeck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deck := vars["deck"]
	var card string
	if deck != "" {
		card = ChooseRandomCard(Decks[deck])
	}
	var imgTag string
	if len(card) > 0 {
		imgTag = fmt.Sprintf(`<img src="/decks/%s/%s"/>`, deck, card)
	}
	_, err := w.Write([]byte(imgTag))
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

func ChooseRandomCard(deck *Deck) string {
	return deck.cardNames[rand.Intn(deck.numCards)]
}
