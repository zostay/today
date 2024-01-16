package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

// const james112 = `Blessed is the man who remains steadfast under trial, for when he has stood the test he will receive the crown of life, which God has promised to those who love him.`

func TestPericope(t *testing.T) {
	t.Parallel()

	p, err := ref.Lookup(ref.Canonical, "James 1:12-14", "Testing")
	assert.NotNil(t, p)
	assert.NoError(t, err)

	b, err := ref.Canonical.Book("James")
	require.NotNil(t, b)
	require.NoError(t, err)

	assert.Equal(t, &ref.Pericope{
		Ref: &ref.Resolved{
			Book:  b,
			First: ref.CV{Chapter: 1, Verse: 12},
			Last:  ref.CV{Chapter: 1, Verse: 14},
		},
		Canon: ref.Canonical,
		Title: "Testing",
	}, p)

	v, err := p.Verses()
	assert.NotNil(t, v)
	assert.NoError(t, err)

	nextV, ok := <-v
	assert.Equal(t, ref.CV{Chapter: 1, Verse: 12}, nextV)
	assert.True(t, ok)

	nextV, ok = <-v
	assert.Equal(t, ref.CV{Chapter: 1, Verse: 13}, nextV)
	assert.True(t, ok)

	nextV, ok = <-v
	assert.Equal(t, ref.CV{Chapter: 1, Verse: 14}, nextV)
	assert.True(t, ok)

	nextV, ok = <-v
	assert.Nil(t, nextV)
	assert.False(t, ok)
}
