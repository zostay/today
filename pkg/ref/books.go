package ref

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound  = errors.New("scripture reference not found")
	ErrWideRange = errors.New("last verse is after the end of the book")
)

// Book is a book of the Bible. We use this with a global map to do client-side
// verification of book names, chapter, and verse references.
type Book struct {
	Name      string
	JustVerse bool
	Verses    []Verse
}

// Canon is a list of books.
type Canon []Book

func (c Canon) Book(name string) (*Book, error) {
	for i := range c {
		b := &c[i]
		if b.Name == name {
			return b, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
}

// Resolve turns an absolute reference into a slice of Resolved references or
// returns an error if the references do not match this Canon.
func (c Canon) Resolve(ref Absolute) ([]Resolved, error) {
	if err := ref.Validate(); err != nil {
		return nil, err
	}

	switch r := ref.(type) {
	case *Multiple:
		return c.resolveMultiple(r)
	case *Proper:
		return c.resolveProper(r)
	case *Resolved:
		return []Resolved{*r}, nil
	}
	return nil, fmt.Errorf("unknown reference type: %T", ref)
}

func (c Canon) resolveMultiple(m *Multiple) ([]Resolved, error) {
	var rs []Resolved
	var b *Book
	for i := range m.Refs {
		switch r := m.Refs[i].(type) {
		case *Proper:
			var err error
			b, err = c.Book(r.Book)
			if err != nil {
				return nil, err
			}

			thisRs, err := c.resolveProper(r)
			if err != nil {
				return nil, err
			}
			rs = append(rs, thisRs...)
		case Relative:
			thisRs, err := c.resolveRelative(b, r)
			if err != nil {
				return nil, err
			}
			rs = append(rs, thisRs...)
		}
	}

	return rs, nil
}

func (c Canon) resolveRelative(b *Book, r Relative) ([]Resolved, error) {
	switch r := r.(type) {
	case *AndFollowing:
		return c.resolveAndFollowing(b, r)
	case *Range:
		return c.resolveRange(b, r)
	case *Related:
		return c.resolveRelated(b, r)
	}
	return nil, fmt.Errorf("unknown reference type: %T", r)
}

func (c Canon) resolveProper(p *Proper) ([]Resolved, error) {
	b, err := c.Book(p.Book)
	if err != nil {
		return nil, err
	}

	switch r := p.Verse.(type) {
	case *Single:
		return c.resolveSingle(b, r)
	case *AndFollowing:
		return c.resolveAndFollowing(b, r)
	case *Range:
		return c.resolveRange(b, r)
	case *Related:
		return c.resolveRelated(b, r)
	}
	return nil, fmt.Errorf("unknown reference type: %T", p.Verse)
}

func (c Canon) resolveSingle(b *Book, s *Single) ([]Resolved, error) {
	return []Resolved{
		{
			Book:  b,
			First: s.Verse,
			Last:  s.Verse,
		},
	}, nil
}

func (c Canon) resolveAndFollowing(
	b *Book,
	a *AndFollowing,
) ([]Resolved, error) {
	switch a.Following {
	case FollowingNone:
		return c.resolveSingle(b, &Single{Verse: a.Verse})
	case FollowingRemainingChapter:
		return c.resolveAndFollowingChapter(b, a.Verse)
	case FollowingRemainingBook:
		return []Resolved{
			{
				Book:  b,
				First: a.Verse,
				Last:  b.Verses[len(b.Verses)-1],
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown following type: %d", a.Following)
}
func (c Canon) resolveAndFollowingChapter(
	b *Book,
	v Verse,
) ([]Resolved, error) {
	ss, err := c.resolveSingle(b, &Single{Verse: v})
	if err != nil {
		return nil, err
	}

	_, hasCv := b.Verses[0].(*CV)
	_, expectCv := ss[0].Last.(*CV)
	if hasCv && !expectCv {
		return nil, errors.New("expected a chapter-and-verse reference, but got verse-only")
	} else if !hasCv && expectCv {
		return nil, errors.New("expected a just-verse reference, but got chapter-and-verse")
	}

	lastVerse := ss[0].Last
	started := false
	for i := range b.Verses {
		if b.Verses[i].Equal(lastVerse) {
			started = true
		} else if !started {
			continue
		}

		if cv, isCv := b.Verses[i].(*CV); started && isCv {
			if cv.Chapter == lastVerse.(*CV).Chapter {
				lastVerse = b.Verses[i]
				continue
			}
			break
		} else {
			lastVerse = b.Verses[len(b.Verses)-1]
			break
		}
	}

	if !started {
		return nil, ErrNotFound
	}

	return []Resolved{
		{
			Book:  b,
			First: v,
			Last:  lastVerse,
		},
	}, nil
}

func (c Canon) resolveRange(
	b *Book,
	r *Range,
) ([]Resolved, error) {
	hasFirst := b.Contains(r.First)
	if !hasFirst {
		return nil, ErrNotFound
	}

	hasLast := b.Contains(r.Last.RelativeTo(r.First))
	if !hasLast {
		return nil, ErrNotFound
	}

	return []Resolved{
		{
			Book:  b,
			First: r.First,
			Last:  r.Last,
		},
	}, nil
}

func (c Canon) resolveRelated(
	b *Book,
	r *Related,
) ([]Resolved, error) {
	var rs []Resolved
	for i := range r.Refs {
		thisRs, err := c.resolveProper(&Proper{
			Book:  b.Name,
			Verse: r.Refs[i],
		})
		if err != nil {
			return nil, err
		}

		rs = append(rs, thisRs...)
	}

	return rs, nil
}

// Contains returns true if the given verse is in the book.
func (b Book) Contains(v Verse) bool {
	for i := range b.Verses {
		if b.Verses[i].Equal(v) {
			return true
		}
	}
	return false
}

// Sub returns a subset of verses in a book. If the given verses are not found,
// this will return nil with an error.
func (b Book) Sub(first, last Verse) ([]Verse, error) {
	var (
		vs        []Verse
		lastFound bool
	)
	for i := range b.Verses {
		if b.Verses[i].Equal(first) || len(vs) > 0 {
			vs = append(vs, b.Verses[i])
		}
		if b.Verses[i].Equal(last) {
			lastFound = true
			break
		}
	}

	if !lastFound {
		return nil, ErrWideRange
	}

	if len(vs) == 0 {
		return nil, ErrNotFound
	}

	return vs, nil
}
