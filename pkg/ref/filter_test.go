package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zostay/today/pkg/ref"
)

func TestCanon_Filtered_ByBook(t *testing.T) {
	t.Parallel()

	c, err := ref.Canonical.Filtered("Matthew", "Mark")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	b, err := c.Book("Luke")
	assert.NoError(t, err)
	assert.NotNil(t, b)

	b, err = c.Book("Matthew")
	assert.ErrorIs(t, err, ref.ErrNotFound)
	assert.Nil(t, b)

	b, err = c.Book("Mark")
	assert.ErrorIs(t, err, ref.ErrNotFound)
	assert.Nil(t, b)

	ps, err := c.Category("Gospels")
	assert.NoError(t, err)

	toRefs := make([]string, len(ps))
	for i, p := range ps {
		toRefs[i], err = p.Ref.CompactRef()
		assert.NoError(t, err)
	}

	assert.Equal(t, []string{"Luke", "John", "Acts"}, toRefs)
}

func TestCanon_Filtered_ByRef(t *testing.T) {
	t.Parallel()

	c, err := ref.Canonical.Filtered("Daniel 7:1-10")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	b, err := c.Book("Daniel")
	assert.NoError(t, err)
	assert.NotNil(t, b)

	ps, err := c.Category("Apocalyptic")
	assert.NoError(t, err)

	toRefs := make([]string, len(ps))
	for i, p := range ps {
		toRefs[i], err = p.Ref.CompactRef()
		assert.NoError(t, err)
	}

	assert.Equal(t, []string{
		"Daniel 7:11-12:13",
		"Revelation",
		"Amos 7:1-9",
		"Amos 8:1-13",
		"Isaiah 24-27",
		"Isaiah 33",
		"Isaiah 55-56",
		"Jeremiah 1:11-16",
		"Ezekiel 38-39",
		"Zechariah 9-14",
		"Joel",
	}, toRefs)
}
