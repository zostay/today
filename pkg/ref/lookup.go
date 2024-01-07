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

func LookupBook(name string) (BookExtract, error) {
	for i := range Canonical {
		b := &Canonical[i]
		if b.Name == name {
			return BookExtract{
				Book:  b,
				First: b.Verses[0],
				Last:  b.Verses[len(b.Verses)-1],
			}, nil
		}
	}
	return BookExtract{}, fmt.Errorf("%w: %s", ErrNotFound, name)
}

func MustLookupBook(name string) BookExtract {
	b, err := LookupBook(name)
	if err != nil {
		panic(err)
	}
	return b
}

func LookupBookExtract(name, first, last string) (BookExtract, error) {
	b, err := LookupBook(name)
	if err != nil {
		return BookExtract{}, err
	}

	opt := &parseVerseOpts{
		expectedRefType: expectEither,
		allowWildcard:   false,
	}
	if b.JustVerse {
		opt.expectedRefType = expectJustVerse
	} else {
		opt.expectedRefType = expectChapterAndVerse
	}

	fv, err := parseVerseRef(first, opt)
	if err != nil {
		return BookExtract{}, fmt.Errorf("unable to parse first ref for extract: %w", err)
	}

	opt.allowWildcard = true
	lv, err := parseVerseRef(last, opt)
	if err != nil {
		return BookExtract{}, fmt.Errorf("unable to parse last ref for extract: %w", err)
	}

	fvi := 0
	for i, verse := range b.Verses() {
		if verse.Equal(fv) {
			fvi = i
			break
		}
	}

	if fvi == 0 {
		return BookExtract{}, fmt.Errorf("%w: %s %s", ErrNotFound, name, first)
	}

	b.First = fv

	if lv.Wildcard() != WildcardNone {
		if b.JustVerse || lv.Wildcard() == WildcardChapter {
			b.Last = b.Verses()[len(b.Verses())-1]
			return b, nil
		}

		if lv.Wildcard() == WildcardVerse {
			b.Last = b.First
			for i := fvi; i < len(b.Verses()); i++ {
				if b.Verses()[i].(*ChapterVerse).chapter != fv.(*ChapterVerse).chapter {
					break
				}
				b.Last = b.Verses()[i]
			}
			return b, nil
		}
	}

	for i := fvi; i < len(b.Verses()); i++ {
		if b.Verses()[i].Equal(lv) {
			b.First = fv
			b.Last = lv
			return b, nil
		}
	}

	if !fv.Before(lv) {
		return BookExtract{}, fmt.Errorf("%w: %s is before %s", ErrReversed, last, first)
	}

	return BookExtract{}, fmt.Errorf("%w: %s %s not found", ErrNotFound, name, last)
}

func MustLookupBookExtract(name, first, last string) BookExtract {
	b, err := LookupBookExtract(name, first, last)
	if err != nil {
		panic(err)
	}
	return b
}

func LookupCategory(name string) ([]BookExtract, error) {
	exs, ok := Categories[name]
	if !ok {
		return nil, ErrNotFound
	}
	return exs, nil
}

func MustLookupCategory(name string) []BookExtract {
	b, err := LookupCategory(name)
	if err != nil {
		panic(err)
	}
	return b
}
