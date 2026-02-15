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

func TestCanon_Resolve_Proper_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Gn.",
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
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_Single_WholeChapter_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book:  "Is.",
			Verse: &ref.Single{Verse: ref.N{Number: 33}},
		},
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_AndFollowingChapter_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isa.",
			Verse: &ref.AndFollowing{
				Verse:     ref.CV{Chapter: 33, Verse: 1},
				Following: ref.FollowingRemainingChapter,
			},
		},
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_AndFollowingBook_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isai",
			Verse: &ref.AndFollowing{
				Verse:     ref.CV{Chapter: 33, Verse: 1},
				Following: ref.FollowingRemainingBook,
			},
		},
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_Range_WholeChapter_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Proper{
			Book: "Isaia",
			Verse: &ref.Range{
				First: ref.N{Number: 24},
				Last:  ref.N{Number: 27},
			},
		},
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_Multiple_Simple_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Multiple{
			Refs: []ref.Ref{
				&ref.Proper{
					Book: "Ge.",
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
					Book: "Ex.",
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
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestCanon_Resolve_Multiple_Relative_Abbr(t *testing.T) {
	t.Parallel()

	rs, err := ref.Canonical.Resolve(
		&ref.Multiple{
			Refs: []ref.Ref{
				&ref.Proper{
					Book: "Ge.",
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
		ref.WithAbbreviations(ref.Abbreviations),
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

func TestBookAbbreviations_BookName(t *testing.T) {
	t.Parallel()

	abbrs := ref.BookAbbreviations{
		Abbreviations: []ref.BookAbbreviation{
			{
				Name:      "Jonah",
				Preferred: "Jonah",
				Accepts: []string{
					"Jonah",
					"Jnh",
				},
			},
			{
				Name:      "John",
				Preferred: "John",
				Accepts: []string{
					"John",
					"Jhn",
					"Jn",
				},
			},
		},
	}

	_, err := abbrs.BookName("J")
	assert.Error(t, err)
	assert.Equal(t, &ref.MultipleMatchError{
		Input: "J",
		Matches: []string{
			"John",
			"Jonah",
		},
	}, err)

	name, err := abbrs.BookName("Jn")
	assert.NoError(t, err)
	assert.Equal(t, "John", name)

	name, err = abbrs.BookName("Jnh")
	assert.NoError(t, err)
	assert.Equal(t, "Jonah", name)

	name, err = abbrs.BookName("Jhn")
	assert.NoError(t, err)
	assert.Equal(t, "John", name)

	name, err = abbrs.BookName("Joh")
	assert.NoError(t, err)
	assert.Equal(t, "John", name)

	name, err = abbrs.BookName("John")
	assert.NoError(t, err)
	assert.Equal(t, "John", name)

	name, err = abbrs.BookName("Jon")
	assert.NoError(t, err)
	assert.Equal(t, "Jonah", name)

	name, err = abbrs.BookName("Jona")
	assert.NoError(t, err)
	assert.Equal(t, "Jonah", name)

	name, err = abbrs.BookName("Jonah")
	assert.NoError(t, err)
	assert.Equal(t, "Jonah", name)

	_, err = abbrs.BookName("Johna")
	assert.ErrorIs(t, err, ref.ErrNotFound)
}

func TestBookAbbreviations_SingularName(t *testing.T) {
	t.Parallel()

	sname, err := ref.Abbreviations.SingularName("John")
	assert.NoError(t, err)
	assert.Equal(t, "John", sname)

	sname, err = ref.Abbreviations.SingularName("Psalms")
	assert.NoError(t, err)
	assert.Equal(t, "Psalm", sname)

	_, err = ref.Abbreviations.SingularName("Psalm")
	assert.ErrorIs(t, err, ref.ErrNotFound)
}

func TestBookAbbreviations_PreferredAbbreviation(t *testing.T) {
	t.Parallel()

	abbrs := ref.BookAbbreviations{
		Abbreviations: []ref.BookAbbreviation{
			{
				Name:      "Genesis",
				Preferred: "Gen.",
				Accepts: []string{
					"Genesis",
					"Gn",
				},
			},
			{
				Name:      "Jonah",
				Preferred: "Jonah",
				Accepts: []string{
					"Jonah",
					"Jnh",
				},
			},
			{
				Name:      "John",
				Preferred: "John",
				Accepts: []string{
					"John",
					"Jhn",
					"Jn",
				},
			},
		},
	}

	abbr, err := abbrs.PreferredAbbreviation("Genesis")
	assert.NoError(t, err)
	assert.Equal(t, "Gen.", abbr)

	abbr, err = abbrs.PreferredAbbreviation("Jonah")
	assert.NoError(t, err)
	assert.Equal(t, "Jonah", abbr)

	abbr, err = abbrs.PreferredAbbreviation("John")
	assert.NoError(t, err)
	assert.Equal(t, "John", abbr)

	_, err = abbrs.PreferredAbbreviation("Jn")
	assert.ErrorIs(t, err, ref.ErrNotFound)
}

func TestBookAbbreviations_NLetterAbbreviation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bookName   string
		n          int
		withPeriod bool
		want       string
		wantErr    bool
	}{
		// 2-letter abbreviations
		{"Genesis 2letter", "Genesis", 2, false, "Gn", false},
		{"Genesis 2letter.", "Genesis", 2, true, "Gn.", false},
		{"Exodus 2letter", "Exodus", 2, false, "Ex", false},
		{"Leviticus 2letter", "Leviticus", 2, false, "Lv", false},
		{"Numbers 2letter", "Numbers", 2, false, "Nm", false},
		{"Deuteronomy 2letter", "Deuteronomy", 2, false, "Dt", false},
		{"Joshua 2letter", "Joshua", 2, false, "Jo", false},
		{"Judges 2letter", "Judges", 2, false, "Jg", false},
		{"Ruth 2letter", "Ruth", 2, false, "Ru", false},
		{"Psalms 2letter", "Psalms", 2, false, "Ps", false},
		{"Romans 2letter", "Romans", 2, false, "Rm", false},

		// 3-letter abbreviations
		{"Genesis 3letter", "Genesis", 3, false, "Gen", false},
		{"Genesis 3letter.", "Genesis", 3, true, "Gen.", false},
		{"Exodus 3letter", "Exodus", 3, false, "Exo", false},
		{"Leviticus 3letter", "Leviticus", 3, false, "Lev", false},
		{"Deuteronomy 3letter", "Deuteronomy", 3, false, "Deu", false},
		{"Joshua 3letter", "Joshua", 3, false, "Jsh", false},
		{"Psalms 3letter", "Psalms", 3, false, "Psm", false},

		// Numbered books - 2 letter
		{"1 Samuel 2letter", "1 Samuel", 2, false, "1 Sm", false},
		{"1 Samuel 2letter.", "1 Samuel", 2, true, "1 Sm.", false},
		{"2 Samuel 2letter", "2 Samuel", 2, false, "2 Sm", false},
		{"1 Kings 2letter", "1 Kings", 2, false, "1 Ki", false},
		{"2 Kings 2letter", "2 Kings", 2, false, "2 Ki", false},
		{"1 Chronicles 2letter", "1 Chronicles", 2, false, "1 Ch", false},
		{"2 Chronicles 2letter", "2 Chronicles", 2, false, "2 Ch", false},
		{"1 Corinthians 2letter", "1 Corinthians", 2, false, "1 Co", false},
		{"2 Corinthians 2letter", "2 Corinthians", 2, false, "2 Co", false},
		{"1 Thessalonians 2letter", "1 Thessalonians", 2, false, "1 Th", false},
		{"2 Thessalonians 2letter", "2 Thessalonians", 2, false, "2 Th", false},
		{"1 Timothy 2letter", "1 Timothy", 2, false, "1 Ti", false},
		{"2 Timothy 2letter", "2 Timothy", 2, false, "2 Ti", false},
		{"1 Peter 2letter", "1 Peter", 2, false, "1 Pt", false},
		{"2 Peter 2letter", "2 Peter", 2, false, "2 Pt", false},
		{"1 John 2letter", "1 John", 2, false, "1 Jn", false},
		{"2 John 2letter", "2 John", 2, false, "2 Jn", false},
		{"3 John 2letter", "3 John", 2, false, "3 Jn", false},

		// Numbered books - 3 letter
		{"1 Samuel 3letter", "1 Samuel", 3, false, "1 Sam", false},
		{"1 Samuel 3letter.", "1 Samuel", 3, true, "1 Sam.", false},
		{"2 Samuel 3letter", "2 Samuel", 3, false, "2 Sam", false},
		{"1 Kings 3letter", "1 Kings", 3, false, "1 Kgs", false},
		{"2 Kings 3letter", "2 Kings", 3, false, "2 Kgs", false},
		{"1 Chronicles 3letter", "1 Chronicles", 3, false, "1 Chr", false},
		{"2 Chronicles 3letter", "2 Chronicles", 3, false, "2 Chr", false},
		{"1 Corinthians 3letter", "1 Corinthians", 3, false, "1 Cor", false},
		{"2 Corinthians 3letter", "2 Corinthians", 3, false, "2 Cor", false},
		{"1 John 3letter", "1 John", 3, false, "1 Jhn", false},
		{"2 John 3letter", "2 John", 3, false, "2 Jhn", false},
		{"3 John 3letter", "3 John", 3, false, "3 Jhn", false},

		// Single-chapter books
		{"Obadiah 2letter", "Obadiah", 2, false, "Ob", false},
		{"Obadiah 3letter", "Obadiah", 3, false, "Oba", false},
		{"Philemon 2letter", "Philemon", 2, false, "Pm", false},
		{"Philemon 3letter", "Philemon", 3, false, "Phm", false},
		{"Jude 2letter", "Jude", 2, false, "Jd", false},
		{"Jude 3letter", "Jude", 3, false, "Jud", false},

		// Books that use available abbreviations
		{"Matthew 2letter", "Matthew", 2, false, "Mt", false},
		{"Matthew 3letter fallback", "Matthew", 3, false, "Mat", false},

		// Error cases
		{"Unknown book", "UnknownBook", 2, false, "", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ref.Abbreviations.NLetterAbbreviation(tt.bookName, tt.n, tt.withPeriod)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ref.ErrNotFound)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
