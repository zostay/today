package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

func TestCanon_Resolve_Proper(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Genesis",
			Verse: &ref.Range{
				First: ref.CV{
					Chapter: 1,
					Verse:   1,
				},
				Last: ref.CV{
					Chapter: 1,
					Verse:   31,
				},
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 1, Verse: 1},
			Last:  ref.CV{Chapter: 1, Verse: 31},
		},
	}, rs)
}

func TestCanon_Resolve_Single_WholeChapter(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book:  "Isaiah",
			Verse: &ref.Single{Verse: ref.N{Number: 33}},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[22],
			First: ref.CV{Chapter: 33, Verse: 1},
			Last:  ref.CV{Chapter: 33, Verse: 24},
		},
	}, rs)
}

func TestCanon_Resolve_AndFollowingChapter(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isaiah",
			Verse: &ref.AndFollowing{
				Verse:     ref.CV{Chapter: 33, Verse: 1},
				Following: ref.FollowingRemainingChapter,
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[22],
			First: ref.CV{Chapter: 33, Verse: 1},
			Last:  ref.CV{Chapter: 33, Verse: 24},
		},
	}, rs)
}

func TestCanon_Resolve_AndFollowingBook(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isaiah",
			Verse: &ref.AndFollowing{
				Verse:     ref.CV{Chapter: 33, Verse: 1},
				Following: ref.FollowingRemainingBook,
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[22],
			First: ref.CV{Chapter: 33, Verse: 1},
			Last:  ref.CV{Chapter: 66, Verse: 24},
		},
	}, rs)
}

func TestCanon_Resolve_Range_WholeChapter(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isaiah",
			Verse: &ref.Range{
				First: ref.N{Number: 24},
				Last:  ref.N{Number: 27},
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[22],
			First: ref.CV{Chapter: 24, Verse: 1},
			Last:  ref.CV{Chapter: 27, Verse: 13},
		},
	}, rs)
}

func TestCanon_Resolve_Multiple_Simple(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Multiple{
			Refs: []ref.Ref{
				&ref.Proper{
					Book: "Genesis",
					Verse: &ref.Range{
						First: ref.CV{
							Chapter: 1,
							Verse:   1,
						},
						Last: ref.CV{
							Chapter: 1,
							Verse:   31,
						},
					},
				},
				&ref.Proper{
					Book: "Exodus",
					Verse: &ref.Range{
						First: ref.CV{
							Chapter: 1,
							Verse:   1,
						},
						Last: ref.CV{
							Chapter: 1,
							Verse:   7,
						},
					},
				},
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 1, Verse: 1},
			Last:  ref.CV{Chapter: 1, Verse: 31},
		},
		{
			Book:  &ref.Canonical.Books[1],
			First: ref.CV{Chapter: 1, Verse: 1},
			Last:  ref.CV{Chapter: 1, Verse: 7},
		},
	}, rs)
}

func TestCanon_Resolve_Resolved(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Resolved{
			Book: &ref.Canonical.Books[0],
			First: ref.CV{
				Chapter: 1,
				Verse:   1,
			},
			Last: ref.CV{
				Chapter: 1,
				Verse:   31,
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 1, Verse: 1},
			Last:  ref.CV{Chapter: 1, Verse: 31},
		},
	}, rs)
}

func TestCanon_Resolve_Multiple_Relative(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Multiple{
			Refs: []ref.Ref{
				&ref.Proper{
					Book: "Genesis",
					Verse: &ref.Range{
						First: ref.CV{
							Chapter: 1,
							Verse:   1,
						},
						Last: ref.CV{
							Chapter: 1,
							Verse:   31,
						},
					},
				},
				&ref.AndFollowing{
					Verse: ref.CV{
						Chapter: 10,
						Verse:   21,
					},
					Following: ref.FollowingRemainingChapter,
				},
				&ref.Range{
					First: ref.CV{
						Chapter: 12,
						Verse:   10,
					},
					Last: ref.CV{
						Chapter: 12,
						Verse:   16,
					},
				},
				&ref.Related{
					Refs: []ref.Relative{
						&ref.Single{
							Verse: ref.CV{
								Chapter: 15,
								Verse:   1,
							},
						},
						&ref.Range{
							First: ref.CV{
								Chapter: 16,
								Verse:   11,
							},
							Last: ref.CV{
								Chapter: 16,
								Verse:   12,
							},
						},
					},
				},
			},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []ref.Resolved{
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 1, Verse: 1},
			Last:  ref.CV{Chapter: 1, Verse: 31},
		},
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 10, Verse: 21},
			Last:  ref.CV{Chapter: 10, Verse: 32},
		},
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 12, Verse: 10},
			Last:  ref.CV{Chapter: 12, Verse: 16},
		},
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 15, Verse: 1},
			Last:  ref.CV{Chapter: 15, Verse: 1},
		},
		{
			Book:  &ref.Canonical.Books[0],
			First: ref.CV{Chapter: 16, Verse: 11},
			Last:  ref.CV{Chapter: 16, Verse: 12},
		},
	}, rs)
}

func TestBook_LastVerseInChapter(t *testing.T) {
	t.Parallel()

	for _, b := range ref.Canonical.Books {
		var lastLv int
		var prevV ref.Verse
		for _, v := range b.Verses {
			if b.JustVerse {
				lv, err := b.LastVerseInChapter(v.(ref.N).Number)
				assert.NoError(t, err)

				assert.Equal(t, lv, b.Verses[len(b.Verses)-1].(ref.N).Number)

				lastLv = lv
			} else {
				lv, err := b.LastVerseInChapter(v.(ref.CV).Chapter)
				assert.NoError(t, err)

				if v.(ref.CV).Verse == 1 {
					if prevV != nil {
						assert.Equal(t, lastLv, prevV.(ref.CV).Verse)
					}

					assert.Greater(t, lv, 1)
					lastLv = lv
				}

				assert.Equal(t, lastLv, lv)
				assert.GreaterOrEqual(t, lastLv, v.(ref.CV).Verse)
			}

			prevV = v
		}

		if b.JustVerse {
			assert.Equal(t, lastLv, prevV.(ref.N).Number)
		} else {
			assert.Equal(t, lastLv, prevV.(ref.CV).Verse)
		}
	}

	gen, err := ref.Canonical.Book("Genesis")
	require.NoError(t, err)

	_, err = gen.LastVerseInChapter(60)
	assert.Error(t, err)
}
