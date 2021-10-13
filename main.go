package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
)

// TODO move this to a new package?
//go:embed static
var embedFS embed.FS
//go:embed templates/*
var templateFS embed.FS
//go:embed css
var cssFS embed.FS

const cardsDir string = "static"

type Deck struct {
	Name     string
	numCards int  // this is redundant, but meh.
	cardNames []string
}

var Decks = make(map[string]*Deck)

func main() {

	useOS := len(os.Args) > 1 && os.Args[1] == "live"

	var staticFS fs.FS
	if useOS {
		staticFS = os.DirFS(".")
	} else {
		staticFS = embedFS
	}

	Decks = FindDecks(staticFS)

	// TODO store the cards a user has seen in a session?

	// load the cards whether they are in the OS FS or the embedded one
	http.Handle("/" + cardsDir + "/", http.FileServer(http.FS(staticFS)))
	// CSS is also static, but separated out so it works regardless of other static files
	http.Handle("/css/", http.FileServer(http.FS(cssFS)))
	// everything else is the main template
	http.HandleFunc("/", serveTemplate)
	log.Println(http.ListenAndServe(":8888", nil))
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	desiredDeck := r.Form.Get("deck")
	var card string
	if desiredDeck != "" {
		card = ChooseRandomCard(Decks[desiredDeck])
	}

	// would prefer to pass this in but a handler has a fixed signature
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Println(err)
	}
	err = tmpl.ExecuteTemplate(w, "layout", struct {
		Card string
		Deck string
		Decks map[string]*Deck
	}{
		Card: card,
		Deck: desiredDeck,
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

	rootEntries, err := fs.ReadDir(rootFS, cardsDir)
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
	deckEntries, err := fs.ReadDir(rootFS, path.Join(cardsDir, deckEntry.Name()))
	if err != nil {
		panic(err)
	}
	return deckEntries
}

func ChooseRandomCard (deck *Deck) string {
	return deck.cardNames[rand.Intn(deck.numCards)]
}


