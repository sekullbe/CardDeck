package decks

import (
	"io/fs"
	"log"
	"math/rand"
	"path"
	"strings"
)

const CardsDir string = "decks"

type Deck struct {
	Name      string
	NumCards  int
	CardNames []string
}

// NewDeck creates a new deck with the given name
func NewDeck(deckName string) *Deck {
	return &Deck{Name: deckName}
}

// FindDecks discovers all available decks in the embedded filesystem
func FindDecks(rootFS fs.FS) map[string]*Deck {
	deckEntries, err := fs.ReadDir(rootFS, CardsDir)
	if err != nil {
		log.Fatal(err)
	}

	decks := make(map[string]*Deck)
	for _, deckEntry := range deckEntries {
		if !deckEntry.IsDir() {
			continue
		}

		deckName := deckEntry.Name()
		deck := NewDeck(deckName)

		cardEntries := getDeckEntries(rootFS, deckEntry)
		deck.NumCards = len(cardEntries)

		for _, cardEntry := range cardEntries {
			deck.CardNames = append(deck.CardNames, cardEntry.Name())
		}

		decks[deckName] = deck
	}

	return decks
}

// getDeckEntries returns the card entries for a given deck
func getDeckEntries(rootFS fs.FS, deckEntry fs.DirEntry) []fs.DirEntry {
	deckPath := path.Join(CardsDir, deckEntry.Name())
	cardEntries, err := fs.ReadDir(rootFS, deckPath)
	if err != nil {
		log.Fatal(err)
	}
	return cardEntries
}

// ChooseUndrawnCard selects a random card that hasn't been drawn yet
func ChooseUndrawnCard(deck *Deck, drawnCards map[string]bool) string {
	var availableCards []string
	for _, cardName := range deck.CardNames {
		if !drawnCards[cardName] {
			availableCards = append(availableCards, cardName)
		}
	}

	if len(availableCards) == 0 {
		return ""
	}

	return availableCards[rand.Intn(len(availableCards))]
}

// ChooseRandomCard selects a random card from the deck (ignoring drawn status)
func ChooseRandomCard(deck *Deck) string {
	return deck.CardNames[rand.Intn(len(deck.CardNames))]
}

// GetCardPath returns the full path to a card file
func GetCardPath(deckName, cardName string) string {
	return path.Join(CardsDir, deckName, cardName)
}

// GetCardNameWithoutExtension removes the file extension from a card name
func GetCardNameWithoutExtension(cardName string) string {
	return strings.TrimSuffix(cardName, path.Ext(cardName))
}
