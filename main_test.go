package main

import (
	"embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed templates decks
var testContent embed.FS

func Test_FindDecks_embedFS(t *testing.T) {

	de, _ := testContent.ReadDir(decksDir)
	deckCount := len(de)

	decks := FindDecks(testContent)
	assert.Equal(t, deckCount, len(decks), "this should have been %d decks but was %d", deckCount, len(decks))

}

func Test_FindDecks_diskFS(t *testing.T) {

	de, _ := testContent.ReadDir(decksDir)
	deckCount := len(de)

	diskFS := os.DirFS(".")
	decks := FindDecks(diskFS)
	assert.Equal(t, deckCount, len(decks), "this should have been %d decks but was %d", deckCount, len(decks))

}

func Test_RandomCard(t *testing.T) {
	deck := NewDeck("testDeck")
	deck.numCards = 1
	deck.cardNames = append(deck.cardNames, "foo")

	card := ChooseRandomCard(deck)
	assert.NotNil(t, card)
	assert.Equal(t, "foo", card)
}
