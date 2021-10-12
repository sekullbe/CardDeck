package main

import (
	"embed"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path"
)

// TODO move this to a new package?
//go:embed static
var cardsFS embed.FS
//go:embed templates/*
var templateFS embed.FS

const cardsDir string = "static"

type Deck struct {
	Name     string
	numCards int  // this is redundant, but meh.
	cardNames []string
}

var Decks = make(map[string]*Deck)

func main() {

	Decks = FindDecks(cardsFS)

	// Handler for the card images only
	// normally put this in something like /static then strip that, but that breaks the embed fs
	http.Handle("/" + cardsDir + "/", http.FileServer(http.FS(cardsFS)))
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

func FindDecks(rootFS embed.FS) map[string]*Deck {

	 decks := make(map[string]*Deck)

	rootEntries, _ := rootFS.ReadDir(cardsDir)
	for _, entry := range rootEntries {
		if entry.IsDir() {
			// Can this be done without loading all the card files? Is that actually happening or only appearing in debug because of the debugger?
			deckEntries, err := rootFS.ReadDir(path.Join(cardsDir, entry.Name()))
			if err != nil {
				panic(err)
			}
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

// TODO store the cards a user has seen in a session?
// TODO index.html is going to have to be a template

func ChooseRandomCard (deck *Deck) string {
	return deck.cardNames[rand.Intn(deck.numCards)]
}


