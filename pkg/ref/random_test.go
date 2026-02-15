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

	for i := range ref.Canonical.Books {
		b := &ref.Canonical.Books[i]

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

func TestRandomWithoutBooks(t *testing.T) {
	t.Parallel()

	r, err := ref.Random(
		ref.ExcludeReferences(
			"Matthew", "Mark", "Luke", "John", "Acts",
			"Romans", "1 Corinthians", "2 Corinthians",
			"Galatians", "Ephesians", "Philippians",
			"Colossians", "1 Thessalonians", "2 Thessalonians",
			"1 Timothy", "2 Timothy", "Titus", "Philemon",
			"Hebrews", "James", "1 Peter", "2 Peter", "1 John", "2 John", "3 John", "Jude",
			"Revelation",
		),
	)

	otBooks := []string{
		"Genesis", "Exodus", "Leviticus", "Numbers", "Deuteronomy",
		"Joshua", "Judges", "Ruth", "1 Samuel", "2 Samuel", "1 Kings", "2 Kings",
		"1 Chronicles", "2 Chronicles", "Ezra", "Nehemiah", "Esther", "Job", "Psalms",
		"Proverbs", "Ecclesiastes", "Song of Solomon", "Isaiah", "Jeremiah",
		"Lamentations", "Ezekiel", "Daniel", "Hosea", "Joel", "Amos", "Obadiah",
		"Jonah", "Micah", "Nahum", "Habakkuk", "Zephaniah", "Haggai", "Zechariah",
		"Malachi",
	}

	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.NoError(t, r.Validate())
	assert.Contains(t, otBooks, r.Book.Name)
}

func TestRandomWithoutBooksFromCategory(t *testing.T) {
	t.Parallel()

	r, err := ref.Random(
		ref.ExcludeReferences(
			"Matthew", "Mark",
		),
		ref.FromCategory("Gospels"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.NoError(t, r.Validate())
	assert.NotEqual(t, "Matthew", r.Book.Name)
	assert.NotEqual(t, "Mark", r.Book.Name)
}

func TestCanon_Random_Filtered_BugFix_Colossians(t *testing.T) {
	t.Parallel()

	filters := []string{
		"Colossians 1:15-2:5",
		"Colossians 1:1-14",
		"Colossians 1:13-20",
		"Colossians 3:12-4:5",
		"Colossians 1:10-1:15",
	}

	c, err := ref.Canonical.Filtered(filters...)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	b, err := c.Book("Colossians")
	assert.NoError(t, err)
	assert.NotNil(t, b)

	expected := []ref.Verse{}
	for i := 6; i <= 23; i++ {
		expected = append(expected, ref.CV{Chapter: 2, Verse: i})
	}
	for i := 1; i < 12; i++ {
		expected = append(expected, ref.CV{Chapter: 3, Verse: i})
	}
	for i := 6; i <= 18; i++ {
		expected = append(expected, ref.CV{Chapter: 4, Verse: i})
	}

	assert.Equal(t, expected, b.Verses)

	r, err := ref.Random(
		ref.FromCanon(c),
		ref.FromBook("Colossians"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, r)
}
