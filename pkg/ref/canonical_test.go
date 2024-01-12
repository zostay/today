package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/ref"
)

func TestCanonical(t *testing.T) {
	t.Parallel()

	assert.Len(t, ref.Canonical, 66)
}

func TestCanonicalBook(t *testing.T) {
	t.Parallel()

	g, err := ref.Canonical.Book("Genesis")
	assert.NoError(t, err)
	assert.Equal(t, "Genesis", g.Name)
	assert.False(t, g.JustVerse)

	first := g.Verses[0]
	assert.Equal(t, (&ref.CV{Chapter: 1, Verse: 1}).Ref(), first.Ref())

	last := g.Verses[len(g.Verses)-1]
	assert.Equal(t, (&ref.CV{Chapter: 50, Verse: 26}).Ref(), last.Ref())
}
