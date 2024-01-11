package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/ref"
)

func TestParseV(t *testing.T) {
	t.Parallel()

	v, err := ref.ParseV("1")
	assert.NoError(t, err)
	assert.Equal(t, &ref.V{Verse: 1}, v)

	v, err = ref.ParseV("23")
	assert.NoError(t, err)
	assert.Equal(t, &ref.V{Verse: 23}, v)

	v, err = ref.ParseV("*")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, v)

	v, err = ref.ParseV("1:2")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, ":2", moreInputErr.Remaining)
	assert.Equal(t, &ref.V{Verse: 1}, v)
}

func TestParseCV(t *testing.T) {
	t.Parallel()

	cv, err := ref.ParseCV("1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.CV{Chapter: 1, Verse: 2}, cv)

	cv, err = ref.ParseCV("23:45")
	assert.NoError(t, err)
	assert.Equal(t, &ref.CV{Chapter: 23, Verse: 45}, cv)

	cv, err = ref.ParseCV("1:*")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, cv)

	cv, err = ref.ParseCV("23:45ff")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "ff", moreInputErr.Remaining)
	assert.Equal(t, &ref.CV{Chapter: 23, Verse: 45}, cv)

	cv, err = ref.ParseCV("1.2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.CV{Chapter: 1, Verse: 2}, cv)

	cv, err = ref.ParseCV("23.45")
	assert.NoError(t, err)
	assert.Equal(t, &ref.CV{Chapter: 23, Verse: 45}, cv)

	cv, err = ref.ParseCV("1.*")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, cv)

	cv, err = ref.ParseCV("23.45ff")
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "ff", moreInputErr.Remaining)
	assert.Equal(t, &ref.CV{Chapter: 23, Verse: 45}, cv)
}

func TestParseSingle(t *testing.T) {
	t.Parallel()

	s, err := ref.ParseSingle("1")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Single{Verse: &ref.V{Verse: 1}}, s)

	s, err = ref.ParseSingle("23")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Single{Verse: &ref.V{Verse: 23}}, s)

	s, err = ref.ParseSingle("*")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, s)

	s, err = ref.ParseSingle("1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}}, s)

	s, err = ref.ParseSingle("23:45")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Single{Verse: &ref.CV{Chapter: 23, Verse: 45}}, s)

	s, err = ref.ParseSingle("1:*")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, ":*", moreInputErr.Remaining)
	assert.Equal(t, &ref.Single{Verse: &ref.V{Verse: 1}}, s)

	s, err = ref.ParseSingle("23:45ff")
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "ff", moreInputErr.Remaining)
	assert.Equal(t, &ref.Single{Verse: &ref.CV{Chapter: 23, Verse: 45}}, s)
}

func TestParseAndFollowing(t *testing.T) {
	t.Parallel()

	af, err := ref.ParseAndFollowing("1ff")
	assert.NoError(t, err)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.V{Verse: 1}, Following: ref.FollowingRemainingChapter}, af)

	af, err = ref.ParseAndFollowing("23ff")
	assert.NoError(t, err)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.V{Verse: 23}, Following: ref.FollowingRemainingChapter}, af)

	af, err = ref.ParseAndFollowing("11")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, af)

	af, err = ref.ParseAndFollowing("1:2ff")
	assert.NoError(t, err)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.CV{Chapter: 1, Verse: 2}, Following: ref.FollowingRemainingChapter}, af)

	af, err = ref.ParseAndFollowing("23:45ff")
	assert.NoError(t, err)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.CV{Chapter: 23, Verse: 45}, Following: ref.FollowingRemainingChapter}, af)

	af, err = ref.ParseAndFollowing("1:*ff")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, af)

	af, err = ref.ParseAndFollowing("23:45ff. ")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, ". ", moreInputErr.Remaining)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.CV{Chapter: 23, Verse: 45}, Following: ref.FollowingRemainingChapter}, af)

	af, err = ref.ParseAndFollowing("23ff:45ff. ")
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, ":45ff. ", moreInputErr.Remaining)
	assert.Equal(t, &ref.AndFollowing{Verse: &ref.V{Verse: 23}, Following: ref.FollowingRemainingChapter}, af)
}

func TestParseRange(t *testing.T) {
	t.Parallel()

	r, err := ref.ParseRange("1:2-3:4")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}}, r)

	r, err = ref.ParseRange("1:2-3")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.V{Verse: 3}}, r)

	r, err = ref.ParseRange("12:34-56:78")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Range{First: &ref.CV{Chapter: 12, Verse: 34}, Last: &ref.CV{Chapter: 56, Verse: 78}}, r)

	r, err = ref.ParseRange("12:34-56")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Range{First: &ref.CV{Chapter: 12, Verse: 34}, Last: &ref.V{Verse: 56}}, r)

	r, err = ref.ParseRange("1-3:4")
	var validationErr *ref.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Nil(t, r)

	r, err = ref.ParseRange("1:2-3:4ff")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "ff", moreInputErr.Remaining)
	assert.Equal(t, &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}}, r)
}

func TestParseRelated(t *testing.T) {
	r, err := ref.ParseRelated("1:2, 4, 6-8, 5:12ff, 12:34-13:56")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Related{
		Refs: []ref.Relative{
			&ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
			&ref.Single{Verse: &ref.V{Verse: 4}},
			&ref.Range{First: &ref.V{Verse: 6}, Last: &ref.V{Verse: 8}},
			&ref.AndFollowing{Verse: &ref.CV{Chapter: 5, Verse: 12}, Following: ref.FollowingRemainingChapter},
			&ref.Range{First: &ref.CV{Chapter: 12, Verse: 34}, Last: &ref.CV{Chapter: 13, Verse: 56}},
		},
	}, r)

	r, err = ref.ParseRelated("4, 6-8, 5:12ff, 12:34-13:56")
	var validationErr *ref.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.ErrorContains(t, err, "but contains a chapter-verse")
	assert.Nil(t, r)

	r, err = ref.ParseRelated("1:2, 4, 6-8, 5:12ff, 12:34-13:56; ")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "; ", moreInputErr.Remaining)
	assert.Equal(t, &ref.Related{
		Refs: []ref.Relative{
			&ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
			&ref.Single{Verse: &ref.V{Verse: 4}},
			&ref.Range{First: &ref.V{Verse: 6}, Last: &ref.V{Verse: 8}},
			&ref.AndFollowing{Verse: &ref.CV{Chapter: 5, Verse: 12}, Following: ref.FollowingRemainingChapter},
			&ref.Range{First: &ref.CV{Chapter: 12, Verse: 34}, Last: &ref.CV{Chapter: 13, Verse: 56}},
		},
	}, r)
}

func TestParseProper(t *testing.T) {
	t.Parallel()

	p, err := ref.ParseProper("Genesis 1:2-3:4")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Range{
			First: &ref.CV{Chapter: 1, Verse: 2},
			Last:  &ref.CV{Chapter: 3, Verse: 4},
		},
	}, p)

	p, err = ref.ParseProper("Genesis 1:2-3")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Range{
			First: &ref.CV{Chapter: 1, Verse: 2},
			Last:  &ref.V{Verse: 3},
		},
	}, p)

	p, err = ref.ParseProper("Genesis 12:34-56:78")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Range{
			First: &ref.CV{Chapter: 12, Verse: 34},
			Last:  &ref.CV{Chapter: 56, Verse: 78},
		},
	}, p)

	p, err = ref.ParseProper("Genesis 12:34-56")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{Book: "Genesis", Verse: &ref.Range{First: &ref.CV{Chapter: 12, Verse: 34}, Last: &ref.V{Verse: 56}}}, p)

	p, err = ref.ParseProper("Genesis 1-3:4")
	var validationErr *ref.ValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Nil(t, p)

	p, err = ref.ParseProper("Genesis 1:2, 3:4ff")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Related{
			Refs: []ref.Relative{
				&ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
				&ref.AndFollowing{Verse: &ref.CV{Chapter: 3, Verse: 4}, Following: ref.FollowingRemainingChapter},
			},
		},
	}, p)

	p, err = ref.ParseProper("Genesis 1:2, 3:4ff; ")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "; ", moreInputErr.Remaining)
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Related{
			Refs: []ref.Relative{
				&ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
				&ref.AndFollowing{Verse: &ref.CV{Chapter: 3, Verse: 4}, Following: ref.FollowingRemainingChapter},
			},
		},
	}, p)

	p, err = ref.ParseProper("G 1:2")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, p)

	p, err = ref.ParseProper("4 1:2")
	assert.ErrorIs(t, err, ref.ErrParseFail)
	assert.Nil(t, p)

	p, err = ref.ParseProper("Ge 1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book:  "Ge",
		Verse: &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
	}, p)

	p, err = ref.ParseProper("4c 1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book:  "4c",
		Verse: &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
	}, p)

	p, err = ref.ParseProper("Ge1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book:  "Ge",
		Verse: &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
	}, p)

	p, err = ref.ParseProper("Ge.1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book:  "Ge.",
		Verse: &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
	}, p)

	p, err = ref.ParseProper("Ge. 1:2")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Proper{
		Book:  "Ge.",
		Verse: &ref.Single{Verse: &ref.CV{Chapter: 1, Verse: 2}},
	}, p)
}

func TestParseMultiple(t *testing.T) {
	m, err := ref.ParseMultiple("Genesis 1:2-3:4, 5:6-7:8")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book: "Genesis",
				Verse: &ref.Related{
					Refs: []ref.Relative{
						&ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}},
						&ref.Range{First: &ref.CV{Chapter: 5, Verse: 6}, Last: &ref.CV{Chapter: 7, Verse: 8}},
					},
				},
			},
		},
	}, m)

	m, err = ref.ParseMultiple("Genesis 1:2-3:4; 5:6-7:8")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}},
			},
			&ref.Range{First: &ref.CV{Chapter: 5, Verse: 6}, Last: &ref.CV{Chapter: 7, Verse: 8}},
		},
	}, m)

	m, err = ref.ParseMultiple("Genesis 1:2-3:4; Ex. 5:6-7:8")
	assert.NoError(t, err)
	assert.Equal(t, &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}},
			},
			&ref.Proper{
				Book:  "Ex.",
				Verse: &ref.Range{First: &ref.CV{Chapter: 5, Verse: 6}, Last: &ref.CV{Chapter: 7, Verse: 8}},
			},
		},
	}, m)

	m, err = ref.ParseMultiple("Genesis 1:2-3:4; x 5:6-7:8")
	var moreInputErr *ref.MoreInputError
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "; x 5:6-7:8", moreInputErr.Remaining)
	assert.Equal(t, &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}},
			},
		},
	}, m)

	m, err = ref.ParseMultiple("Genesis 1:2-3:4; Ex. 5:6-7:8; ")
	assert.ErrorAs(t, err, &moreInputErr)
	assert.Equal(t, "; ", moreInputErr.Remaining)
	assert.Equal(t, &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Range{First: &ref.CV{Chapter: 1, Verse: 2}, Last: &ref.CV{Chapter: 3, Verse: 4}},
			},
			&ref.Proper{
				Book:  "Ex.",
				Verse: &ref.Range{First: &ref.CV{Chapter: 5, Verse: 6}, Last: &ref.CV{Chapter: 7, Verse: 8}},
			},
		},
	}, m)
}
