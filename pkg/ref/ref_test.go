package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
)

func TestCV(t *testing.T) {
	t.Parallel()

	cv := ref.CV{Chapter: 12, Verse: 4}
	assert.Equal(t, "12:4", cv.Ref())
	assert.NoError(t, cv.Validate())
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Single{
			Verse: ref.CV{Chapter: 12, Verse: 4},
		},
	}, cv.InBook("Genesis"))

	assert.True(t, cv.Before(ref.CV{Chapter: 13, Verse: 4}))
	assert.True(t, cv.Before(ref.CV{Chapter: 12, Verse: 5}))
	assert.False(t, cv.Before(ref.CV{Chapter: 11, Verse: 4}))
	assert.False(t, cv.Before(ref.CV{Chapter: 12, Verse: 4}))

	assert.True(t, cv.Before(ref.N{Number: 5}))
	assert.False(t, cv.Before(ref.N{Number: 4}))

	assert.True(t, cv.Equal(ref.CV{Chapter: 12, Verse: 4}))
	assert.False(t, cv.Equal(ref.CV{Chapter: 13, Verse: 4}))
	assert.False(t, cv.Equal(ref.CV{Chapter: 12, Verse: 5}))

	assert.True(t, cv.Equal(ref.N{Number: 4}))
	assert.False(t, cv.Equal(ref.N{Number: 5}))

	assert.Equal(t, ref.CV{Chapter: 12, Verse: 4}, cv.RelativeTo(ref.CV{Chapter: 1, Verse: 1}))
	assert.Equal(t, ref.N{Number: 12}, cv.RelativeTo(ref.N{Number: 1}))

	assert.Error(t, ref.CV{Verse: 1}.Validate())
	assert.Error(t, ref.CV{Chapter: 1}.Validate())
}

func TestN(t *testing.T) {
	t.Parallel()

	n := ref.N{Number: 12}
	assert.Equal(t, "12", n.Ref())
	assert.NoError(t, n.Validate())
	assert.Equal(t, &ref.Proper{
		Book: "Genesis",
		Verse: &ref.Single{
			Verse: ref.N{Number: 12},
		},
	}, n.InBook("Genesis"))

	assert.True(t, n.Before(ref.N{Number: 13}))
	assert.False(t, n.Before(ref.N{Number: 12}))

	assert.True(t, n.Before(ref.CV{Chapter: 13, Verse: 1}))
	assert.False(t, n.Before(ref.CV{Chapter: 12, Verse: 1}))

	assert.True(t, n.Equal(ref.N{Number: 12}))
	assert.False(t, n.Equal(ref.N{Number: 13}))

	assert.True(t, n.Equal(ref.CV{Chapter: 12, Verse: 1}))
	assert.True(t, n.Equal(ref.CV{Chapter: 12, Verse: 2}))
	assert.False(t, n.Equal(ref.CV{Chapter: 13, Verse: 1}))

	assert.Equal(t, ref.N{Number: 12}, n.RelativeTo(ref.N{Number: 1}))
	assert.Equal(t, ref.CV{Chapter: 42, Verse: 12}, n.RelativeTo(ref.CV{Chapter: 42, Verse: 6}))

	assert.Error(t, ref.N{}.Validate())
}

func TestSingle(t *testing.T) {
	t.Parallel()

	s := &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}}
	assert.Equal(t, "12:4", s.Ref())
	assert.NoError(t, s.Validate())

	assert.Equal(t, &ref.Proper{
		Book:  "Genesis",
		Verse: s,
	}, s.InBook("Genesis"))

	assert.Error(t, (&ref.Single{}).Validate())
	assert.Error(t, (&ref.Single{Verse: ref.CV{}}).Validate())
}

func TestAndFollowing(t *testing.T) {
	t.Parallel()

	cf := &ref.AndFollowing{
		Verse:     ref.CV{Chapter: 12, Verse: 4},
		Following: ref.FollowingRemainingChapter,
	}

	bf := &ref.AndFollowing{
		Verse:     ref.CV{Chapter: 12, Verse: 4},
		Following: ref.FollowingRemainingBook,
	}

	assert.Equal(t, "12:4ff", cf.Ref())
	assert.Equal(t, "12:4ffb", bf.Ref())

	assert.NoError(t, cf.Validate())
	assert.NoError(t, bf.Validate())

	assert.Error(t, (&ref.AndFollowing{
		Verse:     ref.CV{Chapter: 12, Verse: 4},
		Following: ref.Following(42),
	}).Validate())
	assert.Error(t, (&ref.AndFollowing{
		Following: ref.FollowingRemainingChapter,
	}).Validate())
	assert.Error(t, (&ref.AndFollowing{
		Verse:     ref.CV{},
		Following: ref.FollowingRemainingChapter,
	}).Validate())

	assert.Equal(t, &ref.Proper{
		Book:  "Genesis",
		Verse: cf,
	}, cf.InBook("Genesis"))
}

func TestRange(t *testing.T) {
	t.Parallel()

	r := &ref.Range{
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.N{Number: 6},
	}

	assert.Equal(t, "12:4-6", r.Ref())

	assert.NoError(t, r.Validate())

	assert.Error(t, (&ref.Range{First: ref.CV{Chapter: 12, Verse: 4}}).Validate())
	assert.Error(t, (&ref.Range{Last: ref.CV{Chapter: 12, Verse: 4}}).Validate())
	assert.Error(t, (&ref.Range{First: ref.CV{}, Last: ref.N{Number: 6}}).Validate())
	assert.Error(t, (&ref.Range{First: ref.CV{Chapter: 12, Verse: 4}, Last: ref.N{}}).Validate())
	assert.Error(t, (&ref.Range{First: ref.CV{}, Last: ref.N{}}).Validate())
	assert.Error(t, (&ref.Range{First: ref.N{Number: 12}, Last: ref.CV{Chapter: 12, Verse: 4}}).Validate())
	assert.Error(t, (&ref.Range{First: ref.N{Number: 12}, Last: ref.N{Number: 12}}).Validate())

	assert.Equal(t, &ref.Proper{
		Book:  "Genesis",
		Verse: r,
	}, r.InBook("Genesis"))
}

func TestRelated(t *testing.T) {
	t.Parallel()

	r := &ref.Related{
		Refs: []ref.Relative{
			&ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
			&ref.Range{First: ref.N{Number: 8}, Last: ref.N{Number: 12}},
			&ref.AndFollowing{
				Verse:     ref.N{Number: 14},
				Following: ref.FollowingRemainingChapter,
			},
		},
	}

	assert.Equal(t, "12:4, 8-12, 14ff", r.Ref())

	assert.NoError(t, r.Validate())

	assert.Error(t, (&ref.Related{}).Validate())
	assert.Error(t, (&ref.Related{Refs: []ref.Relative{}}).Validate())
	assert.Error(t, (&ref.Related{Refs: []ref.Relative{nil}}).Validate())
	assert.Error(t, (&ref.Related{
		Refs: []ref.Relative{
			&ref.Related{
				Refs: []ref.Relative{
					&ref.Single{Verse: ref.N{Number: 12}},
				},
			},
		},
	}).Validate())
	assert.Error(t, (&ref.Related{
		Refs: []ref.Relative{
			&ref.Single{},
		},
	}).Validate())
	assert.Error(t, (&ref.Related{
		Refs: []ref.Relative{
			&ref.Single{Verse: ref.N{Number: 12}},
			&ref.Single{Verse: ref.CV{Chapter: 13, Verse: 4}},
		},
	}).Validate())

	assert.Equal(t, &ref.Proper{
		Book:  "Genesis",
		Verse: r,
	}, r.InBook("Genesis"))
}

func TestProper(t *testing.T) {
	t.Parallel()

	p := &ref.Proper{
		Book:  "Genesis",
		Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
	}

	np := ref.NewProper(
		"Genesis",
		&ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
	)

	assert.Equal(t, np, p)

	assert.Equal(t, "Genesis 12:4", p.Ref())

	assert.NoError(t, p.Validate())

	assert.Error(t, (&ref.Proper{
		Book:  "",
		Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
	}).Validate())
	assert.Error(t, (&ref.Proper{
		Book: "Genesis",
	}).Validate())
	assert.Error(t, (&ref.Proper{
		Book:  "Genesis",
		Verse: &ref.Single{},
	}).Validate())

	assert.Equal(t, []string{"Genesis"}, p.Names())
	assert.True(t, p.IsSingleRange())

	assert.True(t, (&ref.Proper{
		Book:  "Genesis",
		Verse: &ref.AndFollowing{},
	}).IsSingleRange())
	assert.True(t, (&ref.Proper{
		Book:  "Genesis",
		Verse: &ref.Range{},
	}).IsSingleRange())
	assert.False(t, (&ref.Proper{
		Book:  "Genesis",
		Verse: &ref.Related{},
	}).IsSingleRange())
}

func TestMultiple(t *testing.T) {
	t.Parallel()

	m := &ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
			},
			&ref.Range{
				First: ref.CV{Chapter: 13, Verse: 1},
				Last:  ref.N{Number: 4},
			},
		},
	}

	assert.Equal(t, "Genesis 12:4; 13:1-4", m.Ref())

	assert.NoError(t, m.Validate())
	assert.Error(t, (&ref.Multiple{}).Validate())
	assert.Error(t, (&ref.Multiple{Refs: []ref.Ref{}}).Validate())
	assert.Error(t, (&ref.Multiple{
		Refs: []ref.Ref{&ref.Single{Verse: ref.N{Number: 12}}},
	}).Validate())
	assert.NoError(t, (&ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
			},
			ref.CV{Chapter: 13, Verse: 4},
		},
	}).Validate())
	assert.Error(t, (&ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
			},
			&ref.Multiple{
				Refs: []ref.Ref{
					&ref.Proper{
						Book:  "Exodus",
						Verse: &ref.Single{Verse: ref.CV{Chapter: 12, Verse: 4}},
					},
				},
			},
		},
	}).Validate())
	assert.Error(t, (&ref.Multiple{
		Refs: []ref.Ref{
			&ref.Proper{
				Book:  "Genesis",
				Verse: &ref.Single{Verse: ref.CV{}},
			},
		},
	}).Validate())

	assert.Equal(t, []string{"Genesis"}, m.Names())
	assert.Equal(t, []string{"Genesis", "Exodus"}, (&ref.Multiple{
		Refs: []ref.Ref{
			ref.NewProper("Genesis", &ref.Single{Verse: ref.N{Number: 12}}),
			ref.NewProper("Exodus", &ref.Single{Verse: ref.N{Number: 12}}),
		},
	}).Names())

	assert.False(t, m.IsSingleRange())
	assert.True(t, (&ref.Multiple{
		Refs: []ref.Ref{
			ref.NewProper("Genesis", &ref.Single{Verse: ref.N{Number: 12}}),
		},
	}).IsSingleRange())
}

func TestResolved(t *testing.T) {
	t.Parallel()

	gen, err := ref.Canonical.Book("Genesis")
	require.NotNil(t, gen)
	require.NoError(t, err)

	oba, err := ref.Canonical.Book("Obadiah")
	require.NotNil(t, oba)
	require.NoError(t, err)

	r := &ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.CV{Chapter: 12, Verse: 6},
	}

	assert.Equal(t, "Genesis 12:4-12:6", r.Ref())
	assert.Equal(t, "Genesis 12:4", (&ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.CV{Chapter: 12, Verse: 4},
	}).Ref())

	assert.NoError(t, r.Validate())
	assert.Error(t, (&ref.Resolved{
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.N{Number: 6},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book: gen,
		Last: ref.N{Number: 6},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.CV{},
		Last:  ref.N{Number: 6},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.N{},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.CV{Chapter: 11, Verse: 4},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.N{Number: 4},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  gen,
		First: ref.N{Number: 4},
		Last:  ref.CV{Chapter: 12, Verse: 4},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  oba,
		First: ref.CV{Chapter: 12, Verse: 4},
		Last:  ref.N{Number: 4},
	}).Validate())
	assert.Error(t, (&ref.Resolved{
		Book:  oba,
		First: ref.N{Number: 4},
		Last:  ref.CV{Chapter: 12, Verse: 4},
	}).Validate())
	assert.NoError(t, (&ref.Resolved{
		Book:  oba,
		First: ref.N{Number: 4},
		Last:  ref.N{Number: 4},
	}).Validate())

	assert.Equal(t, []string{"Genesis"}, r.Names())
	assert.True(t, r.IsSingleRange())

	assert.Equal(t, []ref.Verse{
		ref.CV{Chapter: 12, Verse: 4},
		ref.CV{Chapter: 12, Verse: 5},
		ref.CV{Chapter: 12, Verse: 6},
	}, r.Verses())
}
