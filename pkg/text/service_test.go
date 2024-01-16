package text_test

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

const fjn41 = `Beloved, do not believe every spirit, but test the spirits to see whether they are from God, for many false prophets have gone out into the world.`

type testResolver struct {
	lastRef *ref.Resolved
}

func (t *testResolver) Verse(ref *ref.Resolved) (string, error) {
	t.lastRef = ref
	return fjn41, nil
}

func (t *testResolver) VerseHTML(ref *ref.Resolved) (template.HTML, error) {
	t.lastRef = ref
	return fjn41, nil
}

var _ text.Resolver = (*testResolver)(nil)

func TestService(t *testing.T) {
	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := ref.Canonical.Book("1 John")
	require.NoError(t, err)
	require.NotNil(t, b)

	txt, err := svc.Verse("1 John 4:1")
	assert.NoError(t, err)
	assert.Equal(t, fjn41, txt)
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)

	htxt, err := svc.VerseHTML("1 John 4:1")
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(fjn41), htxt)
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)
}
