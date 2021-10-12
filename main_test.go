package main

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed templates cards
var testContent embed.FS

func Test_FindDecks(t *testing.T) {

	var decks map[string]*Deck
	decks = FindDecks(testContent)
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



