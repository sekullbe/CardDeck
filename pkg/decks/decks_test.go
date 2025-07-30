package decks

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindDecks_diskFS(t *testing.T) {
	// Use the filesystem from the project root
	diskFS := os.DirFS("../..")
	
	// Read the decks directory to get expected count
	dirEntries, err := fs.ReadDir(diskFS, CardsDir)
	if err != nil {
		t.Skip("Skipping test - decks directory not found")
		return
	}
	
	// Count only directories (actual decks)
	deckCount := 0
	for _, entry := range dirEntries {
		if entry.IsDir() {
			deckCount++
		}
	}

	decks := FindDecks(diskFS)
	assert.Equal(t, deckCount, len(decks), "this should have been %d decks but was %d", deckCount, len(decks))
}

func Test_RandomCard(t *testing.T) {
	deck := NewDeck("testDeck")
	deck.NumCards = 1
	deck.CardNames = append(deck.CardNames, "foo")

	card := ChooseRandomCard(deck)
	assert.NotNil(t, card)
	assert.Equal(t, "foo", card)
}
