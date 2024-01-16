package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

func TestRandomCanonical(t *testing.T) {
	t.Parallel()

	b := ref.RandomCanonical(ref.Canonical)
	assert.NotNil(t, b)

	found := false
	for i := range ref.Canonical.Books {
		if b == &ref.Canonical.Books[i] {
			found = true
			break
		}
	}

	assert.True(t, found)
}

func TestRandomPassage(t *testing.T) {
	t.Parallel()

	for i := range ref.Canonical.Books {
		b := &ref.Canonical.Books[i]

		p := ref.RandomPassage(b)
		assert.NotNil(t, p)
		assert.NotEmpty(t, p)

		for i := range p {
			assert.True(t, b.Contains(p[i]))
		}
	}
}

func TestRandomPassageFromRef(t *testing.T) {
	t.Parallel()

	for i := range ref.Canonical.Books {
		b := &ref.Canonical.Books[i]

		p := ref.RandomPassage(b)
		require.NotNil(t, p)
		require.NotEmpty(t, p)

		firstp := p[0]
		lastp := p[len(p)-1]

		r := &ref.Resolved{
			Book:  b,
			First: p[0],
			Last:  p[len(p)-1],
		}

		vs := ref.RandomPassageFromRef(r)
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
