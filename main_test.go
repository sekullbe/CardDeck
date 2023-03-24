package main

import (
	"embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed templates decks
var testContent embed.FS

// this test is not great; if you add more decks in your local environment it will fail.
func Test_FindDecks(t *testing.T) {

	decks := FindDecks(testContent)
	assert.Equal(t, 2, len(decks), "this should have been 2 decks but was %d", len(decks))

}

func Test_FindDecks_diskFS(t *testing.T) {

	diskFS := os.DirFS(".")
	decks := FindDecks(diskFS)
	assert.Equal(t, 2, len(decks), "this should have been 2 decks but was %d", len(decks))

}

func Test_RandomCard(t *testing.T) {
	deck := NewDeck("testDeck")
	deck.numCards = 1
	deck.cardNames = append(deck.cardNames, "foo")

	card := ChooseRandomCard(deck)
	assert.NotNil(t, card)
	assert.Equal(t, "foo", card)
}
