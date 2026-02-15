package ref

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
	Singular  string
	Ordinal   int
	Accepts   []string
}

// Book will return the Book with the exact given name.
func (c *Canon) Book(in string, opt ...ResolveOption) (*Book, error) {
	name := in

	opts := makeResolveOpts(opt)
	if opts.Abbreviations != nil {
		var err error
		name, err = opts.Abbreviations.BookName(in)
		if err != nil {
			return nil, err
		}
	}

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
				return nil, fmt.Errorf("failed to lookup ref %q: %w", refs[i], err)
			}
			ps = append(ps, p)
		}
		return ps, nil
	}
	return nil, nil
}

type resolveOpts struct {
	Abbreviations *BookAbbreviations
	Singular      bool
}

type ResolveOption func(*resolveOpts)

// WithAbbrevations will allow *Ref methods to use the given BookAbbreviations to
// resolve book names.
func WithAbbreviations(abbrs *BookAbbreviations) ResolveOption {
	return func(o *resolveOpts) {
		o.Abbreviations = abbrs
	}
}

// WithoutAbbrevations will allow *Ref methods to use no abbreviations object
// during name/abbreviation resolution. (Othrwise, these methods will default to
// ref.Abbreviations.)
func WithoutAbbreviations() ResolveOption {
	return func(o *resolveOpts) {
		o.Abbreviations = nil
	}
}

// AsSingleChapter will prefer the singular form of the book name when resolving
// references. (This is special casing for Psalms.)
func AsSingleChapter() ResolveOption {
	return func(o *resolveOpts) {
		o.Singular = true
	}
}

func makeResolveOpts(opts []ResolveOption) *resolveOpts {
	o := &resolveOpts{
		Abbreviations: Abbreviations,
	}
	for i := range opts {
		opts[i](o)
	}
	return o
}

// Clone returns a copy of the Canon.
func (c *Canon) Clone() *Canon {
	newC := Canon{
		Name:       c.Name,
		Books:      make([]Book, len(c.Books)),
		Categories: make(map[string][]string, len(c.Categories)),
	}

	for i := range c.Books {
		newC.Books[i] = c.Books[i].Clone()
	}

	for k, v := range c.Categories {
		newC.Categories[k] = make([]string, len(v))
		copy(newC.Categories[k], v)
	}

	return &newC
}

// Resolve turns an absolute reference into a slice of Resolved references or
// returns an error if the references do not match this Canon.
func (c *Canon) Resolve(ref Absolute, opt ...ResolveOption) ([]Resolved, error) {
	if err := ref.Validate(); err != nil {
		return nil, err
	}

	opts := makeResolveOpts(opt)

	switch r := ref.(type) {
	case *Multiple:
		return c.resolveMultiple(r, opts)
	case *Proper:
		return c.resolveProper(r, opts)
	case *Resolved:
		return []Resolved{*r}, nil
	}
	return nil, fmt.Errorf("unknown reference type: %T", ref)
}

func (c *Canon) resolveBook(in string, opts *resolveOpts) (*Book, error) {
	if opts.Abbreviations != nil {
		name, err := opts.Abbreviations.BookName(in)
		if err != nil {
			return nil, err
		}

		return c.Book(name)
	}

	return c.Book(in)
}

func (c *Canon) resolveMultiple(m *Multiple, opts *resolveOpts) ([]Resolved, error) {
	var rs []Resolved
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

func (c *Canon) resolveProper(p *Proper, opts *resolveOpts) ([]Resolved, error) {
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

func ensureVerseMatchesBook(b *Book, v Verse) (Verse, bool, error) {
	wholeChapter := false

	// enforce N to CV
	if !b.JustVerse {
		if nv, isN := v.(N); isN {
			wholeChapter = true

			// if we have a chapter-only reference, we need to find the first
			// verse of the chapter, which might not be 1.
			found := false
			for i := range b.Verses {
				if b.Verses[i].(CV).Chapter == nv.Number {
					v = CV{Chapter: nv.Number, Verse: b.Verses[i].(CV).Verse}
					found = true
					break
				}
			}

			if !found {
				v = CV{Chapter: nv.Number, Verse: 1}
			}
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
		}, &resolveOpts{})
		if err != nil {
			return nil, err
		}

		rs = append(rs, thisRs...)
	}

	return rs, nil
}

func (b Book) Clone() Book {
	newB := Book{
		Name:      b.Name,
		Verses:    make([]Verse, len(b.Verses)),
		JustVerse: b.JustVerse,
	}
	copy(newB.Verses, b.Verses)
	return newB
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

// NLetterAbbreviation returns an N-letter abbreviation for the given book name.
// It searches the Accepts slice for the first abbreviation with exactly N letters
// (ignoring spaces, numbers, and periods). If no such abbreviation is found, it
// truncates the book name to N letters. If withPeriod is true, a period is appended.
//
// For numbered books (e.g., "1 John"), the number prefix is preserved and only the
// book name portion is abbreviated.
func (b *BookAbbreviations) NLetterAbbreviation(name string, n int, withPeriod bool) (string, error) {
	for _, abbr := range b.Abbreviations {
		if abbr.Name != name {
			continue
		}

		// Extract number prefix if present (e.g., "1" from "1 John")
		var prefix string
		bookName := name
		if abbr.Ordinal > 0 {
			// Find the space after the number
			spaceIdx := strings.Index(name, " ")
			if spaceIdx > 0 {
				prefix = name[:spaceIdx+1] // Include the space
				bookName = name[spaceIdx+1:]
			}
		}

		// Search for N-letter abbreviation in Accepts
		for _, accept := range abbr.Accepts {
			// Remove number prefix from accept string for comparison
			acceptName := accept
			if abbr.Ordinal > 0 {
				// Remove leading digits and roman numerals
				acceptName = strings.TrimLeft(accept, "0123456789IⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩstndrdthFirstSecondThird")
			}

			// Count only letters (ignore spaces, periods, numbers)
			letterCount := 0
			for _, ch := range acceptName {
				if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
					letterCount++
				}
			}

			if letterCount == n {
				result := prefix + acceptName
				if withPeriod && !strings.HasSuffix(result, ".") {
					result += "."
				}
				return result, nil
			}
		}

		// Fallback: truncate book name to N letters
		letters := make([]rune, 0, n)
		for _, ch := range bookName {
			if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				letters = append(letters, ch)
				if len(letters) == n {
					break
				}
			}
		}

		result := prefix + string(letters)
		if withPeriod {
			result += "."
		}
		return result, nil
	}

	return "", fmt.Errorf("%w: %s", ErrNotFound, name)
}
