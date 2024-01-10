package ref

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("scripture reference not found")
	ErrMismatch = errors.New("scripture reference mismatch (chapter/verse vs. verse)")
	ErrReversed = errors.New("scripture reference is reversed")
)

func LookupBook(name string) (*Resolved, error) {
	for i := range Canonical {
		b := &Canonical[i]
		if b.Name == name {
			return &Resolved{
				Book:  b,
				First: b.Verses[0],
				Last:  b.Verses[len(b.Verses)-1],
			}, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
}

func MustLookupBook(name string) *Resolved {
	b, err := LookupBook(name)
	if err != nil {
		panic(err)
	}
	return b
}

func LookupBookExtract(name, first, last string) (*Resolved, error) {
	b, err := LookupBook(name)
	if err != nil {
		return nil, err
	}

	opt := &parseVerseOpts{
		expectedRefType: expectEither,
		allowWildcard:   false,
	}
	if b.Book.JustVerse {
		opt.expectedRefType = expectJustVerse
	} else {
		opt.expectedRefType = expectChapterAndVerse
	}

	fv, err := parseVerseRef(first, opt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse first ref for extract: %w", err)
	}

	opt.allowWildcard = true
	lv, err := parseVerseRef(last, opt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse last ref for extract: %w", err)
	}

	fvi := 0
	for i, verse := range b.Verses() {
		if verse.Equal(fv) {
			fvi = i
			break
		}
	}

	if fvi == 0 {
		return nil, fmt.Errorf("%w: %s %s", ErrNotFound, name, first)
	}

	b.First = fv

	for i := fvi; i < len(b.Verses()); i++ {
		if b.Verses()[i].Equal(lv) {
			b.First = fv
			b.Last = lv
			return b, nil
		}
	}

	if !fv.Before(lv) {
		return nil, fmt.Errorf("%w: %s is before %s", ErrReversed, last, first)
	}

	return nil, fmt.Errorf("%w: %s %s not found", ErrNotFound, name, last)
}

func MustLookupBookExtract(name, first, last string) *Resolved {
	b, err := LookupBookExtract(name, first, last)
	if err != nil {
		panic(err)
	}
	return b
}

//func LookupCategory(name string) ([]Resolved, error) {
//	exs, ok := Categories[name]
//	if !ok {
//		return nil, ErrNotFound
//	}
//	return exs, nil
//}
//
//func MustLookupCategory(name string) []Resolved {
//	b, err := LookupCategory(name)
//	if err != nil {
//		panic(err)
//	}
//	return b
//}
