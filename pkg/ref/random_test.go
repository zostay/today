package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/ref"
)

func TestRandomCanonical(t *testing.T) {
	t.Parallel()

	b := ref.RandomCanonical(bible.Protestant)
	assert.NotNil(t, b)

	found := false
	for i := range bible.Protestant.Books {
		if b == &bible.Protestant.Books[i] {
			found = true
			break
		}
	}

	assert.True(t, found)
}

func TestRandomPassage(t *testing.T) {
	t.Parallel()

	for i := range bible.Protestant.Books {
		b := &bible.Protestant.Books[i]

		p := ref.RandomPassage(b, 1, 30)
		assert.NotNil(t, p)
		assert.NotEmpty(t, p)

		for i := range p {
			assert.True(t, b.Contains(p[i]))
		}
	}
}

func TestRandomPassageFromRef(t *testing.T) {
	t.Parallel()

	for i := range bible.Protestant.Books {
		b := &bible.Protestant.Books[i]

		p := ref.RandomPassage(b, 1, 30)
		require.NotNil(t, p)
		require.NotEmpty(t, p)

		firstp := p[0]
		lastp := p[len(p)-1]

		r := &ref.Resolved{
			Book:  b,
			First: p[0],
			Last:  p[len(p)-1],
		}

		vs := ref.RandomPassageFromRef(r, 1, 30)
		require.NotNil(t, vs)
		require.NotEmpty(t, vs)

		firstv := vs[0]
		lastv := vs[len(vs)-1]

		// p range should contain v range
		assert.True(t, firstp.Equal(firstv) || firstp.Before(lastv))
		assert.True(t, firstv.Equal(lastp) || firstv.Before(lastp))
		assert.True(t, firstp.Equal(lastv) || firstp.Before(lastv))
		assert.True(t, lastv.Equal(lastp) || lastv.Before(lastp))
	}
}

func TestRandom(t *testing.T) {
	t.Parallel()

	r, err := ref.Random()
	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.NoError(t, r.Validate())

	r, err = ref.Random(ref.FromCategory("Gospels"))
	assert.NoError(t, err)
	assert.NotNil(t, r)

	// TODO Make sure it matches the category

	assert.NoError(t, r.Validate())

	r, err = ref.Random(ref.FromBook("John"))
	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.Equal(t, "John", r.Book.Name)

	assert.NoError(t, r.Validate())
}
