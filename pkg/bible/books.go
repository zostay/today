package bible

import (
	"errors"
	"fmt"
	"sort"
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

// V is a highly simplified version of a verse. If the book is a JustVerse book,
// then this will be a single verse number in V. Otherwise, this will be a
// chapter and verse number in (C, V).
type V struct {
	C int
	V int
}

// Equal returns true if the two verses are the same.
func (v V) Equal(o V) bool {
	return v.C == o.C && v.V == o.V
}

// Ref returns a string representation of the verse.
func (v V) Ref() string {
	if v.C == 0 {
		return fmt.Sprintf("%d", v.V)
	}
	return fmt.Sprintf("%d:%d", v.C, v.V)
}

// Validate returns true iff the chapter and verse are both positive.
func (v V) Validate(justVerse bool) error {
	if !justVerse && v.C <= 0 {
		return invalid("chapter must be positive")
	}
	if v.V <= 0 {
		return invalid("verse must be positive")
	}
	return nil
}

// RelativeTo returns a new verse that is relative to the given verse. This is
// used to calculate the last verse in a range.
func (v V) RelativeTo(o V) V {
	if v.C == 0 {
		return V{
			C: o.C,
			V: v.V,
		}
	}
	return v
}

// Before returns true if the verse is before the other verse.
func (v V) Before(o V) bool {
	if v.C == o.C {
		return v.V < o.V
	}
	return v.C < o.C
}

// Book is a book of the Bible. We use this with a global map to do client-side
// verification of book names, chapter, and verse references.
type Book struct {
	Name      string
	JustVerse bool
	Verses    []V
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
	Singular  string
	Ordinal   int
	Accepts   []string
}

// Resolve turns an absolute reference into a slice of Resolved references or
// returns an error if the references do not match this Canon.
func (c *Canon) Resolve(r AbsoluteRef, opt ...ResolveOption) ([]CanonicalRef, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	opts := makeResolveOpts(opt)

	switch r := r.(type) {
	case *Multiple:
		return c.resolveMultiple(r, opts)
	case *Proper:
		return c.resolveProper(r, opts)
	case *CanonicalRef:
		return []CanonicalRef{*r}, nil
	}
	return nil, fmt.Errorf("unknown reference type: %T", r)
}


func (c *Canon) resolveMultiple(m *Multiple, opts *ResolveOptions) ([]CanonicalRef, error) {
	var rs []CanonicalRef
	var b *Book
	for i := range m.Refs {
		switch r := m.Refs[i].(type) {
		case *Proper:
			var err error
			b, err = c.resolveBook(r.Book, opts)
			if err != nil {
				return nil, err
			}

			thisRs, err := c.resolveProper(r, opts)
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

func (c *Canon) resolveRelative(b *Book, r RelativeRef) ([]CanonicalRef, error) {
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

func (c *Canon) resolveProper(p *Proper, opts *ResolveOptions) ([]CanonicalRef, error) {
	b, err := c.resolveBook(p.Book, opts)
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

func ensureVerseMatchesBook(b *Book, v VerseRef) (VerseRef, bool, error) {
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

func (c *Canon) resolveSingle(b *Book, s *Single) ([]CanonicalRef, error) {
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

	return []CanonicalRef{
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
) ([]CanonicalRef, error) {
	v, _, err := ensureVerseMatchesBook(b, a.Verse)
	if err != nil {
		return nil, err
	}

	if !b.Contains(v) {
		return nil, ErrNotFound
	}

	switch a.Following { //nolint:exhaustive // we don't need to handle all cases
	case FollowingRemainingBook:
		return []CanonicalRef{
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

		return []CanonicalRef{
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
	v VerseRef,
) (VerseRef, error) {
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
) ([]CanonicalRef, error) {
	first, wholeChapter, err := ensureVerseMatchesBook(b, r.First)
	if first == nil {
		return nil, err
	}

	hasFirst := b.Contains(first)
	if !hasFirst {
		return nil, ErrNotFound
	}

	var last VerseRef
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

	return []CanonicalRef{
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
) ([]CanonicalRef, error) {
	var rs []CanonicalRef
	for i := range r.Refs {
		thisRs, err := c.resolveProper(&Proper{
			Book:  b.Name,
			Verse: r.Refs[i],
		}, &Resolve{})
		if err != nil {
			return nil, err
		}

		rs = append(rs, thisRs...)
	}

	return rs, nil
}

// Contains returns true if the given verse is in the book.
func (b Book) Contains(v VerseRef) bool {
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

	sort.Strings(matchNames)

	return "", &MultipleMatchError{
		Input:   in,
		Matches: matchNames,
	}
}

// SingularName returns the name to use for the singular form of the book name.
// This is basically special casing for Psalms, which are quoted as Psalms 12-14
// when multiple chapters are cited, but as Psalm 12 when a single chapter is
// cited. The given name should be the full name of the book as resolved via the
// BookName method.
//
// If the given book is not found in the abbreviations, this will return
// ErrNotFound.
func (b *BookAbbreviations) SingularName(name string) (string, error) {
	for _, abbr := range b.Abbreviations {
		if abbr.Name == name {
			if abbr.Singular != "" {
				return abbr.Singular, nil
			}
			return abbr.Name, nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrNotFound, name)
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
