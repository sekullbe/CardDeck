package main

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//go:embed templates static
var testContent embed.FS

// this test works in my local environment that has four decks but those aren't mine to share so aren't distributed
func Test_FindDecks(t *testing.T) {

	var decks map[string]*Deck
	decks = FindDecks(testContent)
	assert.Equal(t, 4, len(decks), "this should have been four decks but was %d", len(decks))

}

func Test_FindDecks_diskFS(t *testing.T) {

	diskFS := os.DirFS(".")
	var decks map[string]*Deck
	decks = FindDecks(diskFS)
	assert.Equal(t, 4, len(decks), "this should have been four decks but was %d", len(decks))

}

func Test_RandomCard(t *testing.T) {
	deck := NewDeck("testDeck")
	deck.numCards = 1
	deck.cardNames = append(deck.cardNames, "foo")

	card := ChooseRandomCard(deck)
	assert.NotNil(t, card)
	assert.Equal(t, "foo", card)
}



