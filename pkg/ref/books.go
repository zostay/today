package ref

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a reference is not found in the canon.
	ErrNotFound = errors.New("scripture reference not found")

	// ErrWideRange is returned when the second reference in a range is not
	// found.
	ErrWideRange = errors.New("last verse is after the end of the book")
)

type MultipleMatchError struct {
	Input   string
	Matches []string
}

func (m *MultipleMatchError) Error() string {
	return fmt.Sprintf("input string %q has multiple matches: %v", m.Input, m.Matches)
}

// Book is a book of the Bible. We use this with a global map to do client-side
// verification of book names, chapter, and verse references.
type Book struct {
	Name      string
	JustVerse bool
	Verses    []Verse
}

// Canon is primarily a collection of books, but may include other metadata.
type Canon struct {
	Name       string
	Books      []Book
	Categories map[string][]string
}

// BookAbbreviations is configuration for book names and abbreviations according
// to a standardized scheme. This allows for configurable preferences for
// abbreviations when citing references and for more flexible parsing of
// references.
type BookAbbreviations struct {
	Abbreviations []BookAbbreviation
	root          *AbbrTree
}

// BookAbbreviation is an individual configuration of a book name, selects a
// standard abbreviation, and provides a place for recording accepted
// abbreviations when parsing book names.
type BookAbbreviation struct {
	Name      string
	Preferred string
	Accepts   []string
}

// Book will return the Book with the given name or any
func (c *Canon) Book(name string) (*Book, error) {
	for i := range c.Books {
		b := &c.Books[i]
		if b.Name == name {
			return b, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
}

// Category returns a list of Pericopes associated with that Category or nil if
// no such category is defined. Returns nil and error if there's a problem with
// the category definition.
func (c *Canon) Category(name string) ([]*Pericope, error) {
	if refs, hasCategory := c.Categories[name]; hasCategory {
		var ps []*Pericope
		for i := range refs {
			p, err := Lookup(c, refs[i], "")
			if err != nil {
				return nil, err
			}
			ps = append(ps, p)
		}
		return ps, nil
	}
	return nil, nil
}

// Resolve turns an absolute reference into a slice of Resolved references or
// returns an error if the references do not match this Canon.
func (c *Canon) Resolve(ref Absolute) ([]Resolved, error) {
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

func (c *Canon) resolveMultiple(m *Multiple) ([]Resolved, error) {
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

func (c *Canon) resolveRelative(b *Book, r Relative) ([]Resolved, error) {
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

func (c *Canon) resolveProper(p *Proper) ([]Resolved, error) {
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

func ensureVerseMatchesBook(b *Book, v Verse) (Verse, bool, error) {
	wholeChapter := false

	// enforce N to CV
	if !b.JustVerse {
		if nv, isN := v.(N); isN {
			wholeChapter = true
			v = CV{Chapter: nv.Number, Verse: 1}
		}
	}

	if _, isCV := v.(CV); b.JustVerse && isCV {
		return nil, false, errors.New("expected a verse-only reference, but got chapter-and-verse")
	}

	return v, wholeChapter, nil
}

func (c *Canon) resolveSingle(b *Book, s *Single) ([]Resolved, error) {
	v, wholeChapter, err := ensureVerseMatchesBook(b, s.Verse)
	if err != nil {
		return nil, err
	}

	if wholeChapter {
		return c.resolveAndFollowing(b, &AndFollowing{
			Verse:     v,
			Following: FollowingRemainingChapter,
		})
	}

	if !b.Contains(v) {
		return nil, ErrNotFound
	}

	return []Resolved{
		{
			Book:  b,
			First: v,
			Last:  v,
		},
	}, nil
}

func (c *Canon) resolveAndFollowing(
	b *Book,
	a *AndFollowing,
) ([]Resolved, error) {
	v, _, err := ensureVerseMatchesBook(b, a.Verse)
	if err != nil {
		return nil, err
	}

	if !b.Contains(v) {
		return nil, ErrNotFound
	}

	switch a.Following { //nolint:exhaustive // we don't need to handle all cases
	case FollowingRemainingBook:
		return []Resolved{
			{
				Book:  b,
				First: v,
				Last:  b.Verses[len(b.Verses)-1],
			},
		}, nil
	default:
		lv, err := lastVerseInChapter(b, v)
		if err != nil {
			return nil, err
		}

		return []Resolved{
			{
				Book:  b,
				First: v,
				Last:  lv,
			},
		}, nil
	}
}

func (b Book) LastVerseInChapter(
	n int,
) (int, error) {
	if b.JustVerse {
		return b.Verses[len(b.Verses)-1].(N).Number, nil
	}

	fv := CV{Chapter: n, Verse: 1}
	lv, err := lastVerseInChapter(&b, fv)

	if err != nil {
		return 0, err
	}

	return lv.(CV).Verse, nil
}

func lastVerseInChapter(
	b *Book,
	v Verse,
) (Verse, error) {
	if b.JustVerse {
		return b.Verses[len(b.Verses)-1], nil
	}

	lv := v
	started := false
	for i := range b.Verses {
		if !started {
			if b.Verses[i].Equal(lv) {
				started = true
			}
			continue
		}

		cv := b.Verses[i].(CV)
		if !(N{Number: cv.Chapter}).Equal(v) {
			break
		}

		lv = b.Verses[i]
	}

	if !started {
		return nil, ErrNotFound
	}

	return lv, nil
}

func (c *Canon) resolveRange(
	b *Book,
	r *Range,
) ([]Resolved, error) {
	first, wholeChapter, err := ensureVerseMatchesBook(b, r.First)
	if first == nil {
		return nil, err
	}

	hasFirst := b.Contains(first)
	if !hasFirst {
		return nil, ErrNotFound
	}

	var last Verse
	if wholeChapter {
		last, _, err = ensureVerseMatchesBook(b, r.Last)
		if err != nil {
			return nil, err
		}

		last, err = lastVerseInChapter(b, last)
		if err != nil {
			return nil, err
		}
	} else {
		last = r.Last.RelativeTo(first)
	}

	hasLast := b.Contains(last)
	if !hasLast {
		return nil, ErrNotFound
	}

	return []Resolved{
		{
			Book:  b,
			First: first,
			Last:  last,
		},
	}, nil
}

func (c *Canon) resolveRelated(
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

// BookName returns the book name that matches the given string. This will apply as
// liberal a match as possible against the abberviations in the configurations.
// The word is checked against all possible abbreviations.
//
// If there are no matches, this will return ErrNotFound. If there are multiple
// matches, this will return a MultipleMatchError, which can be interrogated to
// determine all book names that matched.
func (b *BookAbbreviations) BookName(in string) (string, error) {
	if b.root == nil {
		b.root = NewAbbrTree(b)
	}

	matches := b.root.Get(in)
	if len(matches) == 1 {
		for _, m := range matches {
			return m.Name, nil
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("%w: %s", ErrNotFound, in)
	}

	matchNames := make([]string, 0, len(matches))
	for _, m := range matches {
		matchNames = append(matchNames, m.Name)
	}

	return "", &MultipleMatchError{
		Input:   in,
		Matches: matchNames,
	}
}

// PreferredAbbreviation returns the preferred abbreviation for the given book
// name.
func (b *BookAbbreviations) PreferredAbbreviation(name string) (string, error) {
	for _, abbr := range b.Abbreviations {
		if abbr.Name == name {
			return abbr.Preferred, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrNotFound, name)
}
