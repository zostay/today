package bible_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/canon"
)

func TestProtestant(t *testing.T) {
	t.Parallel()

	assert.Len(t, bible.Protestant.Books, 66)
}

func TestCanonicalBook(t *testing.T) {
	t.Parallel()

	g, err := bible.Protestant.Book("Genesis")
	assert.NoError(t, err)
	assert.Equal(t, "Genesis", g.Name)
	assert.False(t, g.JustVerse)

	first := g.Verses[0]
	assert.Equal(t, canon.V{C: 1, V: 1}, first)

	last := g.Verses[len(g.Verses)-1]
	assert.Equal(t, canon.V{C: 50, V: 26}, last)
}
