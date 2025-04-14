package ref

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/agnivade/levenshtein"
)

type randomOpts struct {
	category string
	book     string
	canon    *Canon
	min, max int
	exclude  []string
}

type RandomReferenceOption func(*randomOpts)

func FromCanon(canon *Canon) RandomReferenceOption {
	return func(o *randomOpts) {
		o.canon = canon
	}
}

func FromBook(name string) RandomReferenceOption {
	return func(o *randomOpts) {
		o.book = name
	}
}

func FromCategory(name string) RandomReferenceOption {
	return func(o *randomOpts) {
		o.category = name
	}
}

func WithAtLeast(n uint) RandomReferenceOption {
	return func(o *randomOpts) {
		if n > math.MaxInt {
			n = math.MaxInt
		}
		o.min = int(n)
	}
}

func WithAtMost(n uint) RandomReferenceOption {
	return func(o *randomOpts) {
		if n > math.MaxInt {
			n = math.MaxInt
		}
		o.max = int(n)
	}
}

func ExcludeReferences(verses ...string) RandomReferenceOption {
	return func(o *randomOpts) {
		o.exclude = append(o.exclude, verses...)
	}
}

type UnknownCategoryError struct {
	Category      string
	Possibilities []string
}

func (o *UnknownCategoryError) Error() string {
	alternates := ""
	if o.Possibilities != nil {
		alternates = fmt.Sprintf(". Did you mean?\n\n - %s",
			strings.Join(o.Possibilities, "\n - "))
	}

	return fmt.Sprintf("unknown category: %s%s",
		o.Category,
		alternates)
}

// Random pulls a random reference from the Bible and returns it. You can use the
// options to help narrow down where the passages are selected from.
func Random(opt ...RandomReferenceOption) (*Resolved, error) {
	o := &randomOpts{
		canon: Canonical,
		min:   1,
		max:   30,
	}
	for _, f := range opt {
		f(o)
	}

	var err error
	if len(o.exclude) > 0 {
		o.canon, err = o.canon.Filtered(o.exclude...)
		if err != nil {
			return nil, fmt.Errorf("error while filtering for requested verses: %w", err)
		}
	}

	var (
		b  *Book
		vs []Verse
	)

	if o.category != "" {
		_, hasCategory := o.canon.Categories[o.category]
		if !hasCategory {
			var possibilities []string
			for cat := range o.canon.Categories {
				if levenshtein.ComputeDistance(o.category, cat) <= 4 {
					possibilities = append(possibilities, cat)
				}
			}

			return nil, &UnknownCategoryError{
				Category:      o.category,
				Possibilities: possibilities,
			}
		}

		ps, err := o.canon.Category(o.category)
		if err != nil {
			return nil, fmt.Errorf("error getting category pericopes %q: %w", o.category, err)
		}

		// lazy way to weight the books by the number of verses they have
		bag := make([]*Pericope, 0, len(ps))
		for i := range ps {
			for range ps[i].Ref.Verses() {
				bag = append(bag, ps[i])
			}
		}

		be := bag[rand.Int()%len(bag)] //nolint:gosec // weak random is fine here
		b = be.Ref.Book
		vs = RandomPassageFromRef(be.Ref, o.min, o.max)
	} else {
		if o.book != "" {
			b, err = o.canon.Book(o.book)
			if err != nil {
				return nil, fmt.Errorf("error looking up book %q: %w", o.book, err)
			}

			firstVerse := b.Verses[0]
			ex, err := Lookup(o.canon, o.book+" "+firstVerse.Ref()+"ffb", "")
			if err != nil {
				return nil, fmt.Errorf("error looking up book %q: %w", o.book, err)
			}

			b = ex.Ref.Book
		} else {
			b = RandomCanonical(o.canon)
		}

		vs = RandomPassage(b, o.min, o.max)
	}

	v1, v2 := vs[0], vs[len(vs)-1]

	return &Resolved{
		Book:  b,
		First: v1,
		Last:  v2,
	}, nil
}

// RandomCanonical returns a random book of the Bible.
func RandomCanonical(c *Canon) *Book {
	return &c.Books[rand.Int()%len(c.Books)] //nolint:gosec // weak random is fine here
}

// RandomPassage returns a random passage from the given book of the Bible. It
// returns all the verses in the canon for the selected passage. The first being
// the start of the passage and the second being the end of the passage. As of
// this writing, the verse references may be up to 30 Verses apart.
//
// You could turn the result into a ref.Resolved via:
//
//	vs := ref.RandomPassage(b, min, max)
//	r := &ref.Resolved{
//	  Book:  b,
//	  First: vs[0],
//	  Last:  vs[len(vs)-1],
//	}
//
// You must select the minimum and maximum number of verses to include in the
// passage. If you want a single verse, set both the minimum and maximum to 1.
// The values will be automatically capped to the number of verses in the book
// and automatically set to 1 if they are less than 1.
func RandomPassage(b *Book, mn, mx int) []Verse {
	return pickVerses(b.Verses, mn, mx)
}

// RandomPassageFromRef returns a random passage from the given ref.Resolved of
// the Bible. It returns a two element slice of ref.Verse, the first being the
// start of the passage and the second being the end of the passage. As of this
// writing, the verse references may be up to 30 Verses apart.
//
// You could turn the result into a ref.Resolved via:
//
//	vs := ref.RandomPassage(r1, min, max)
//	r2 := &ref.Resolved{
//	  Book:  r1.Book,
//	  First: vs[0],
//	  Last:  vs[len(vs)-1],
//	}
//
// You must select the minimum and maximum number of verses to include in the
// passage. If you want a single verse, set both the minimum and maximum to 1.
// The values will be automatically capped to the number of verses in the book
// and automatically set to 1 if they are less than 1.
func RandomPassageFromRef(b *Resolved, mn, mx int) []Verse {
	return pickVerses(b.Verses(), mn, mx)
}

func pickVerses(verses []Verse, mn, mx int) []Verse {
	// This is a little convoluted, but let me explain:
	//
	// * User selects the minimum and maximum length of the passage to return in
	//   number of verses.
	//
	// * The value of min and max must each be less than the number of verses in
	//   the book (or we cap the values to that number).
	//
	// * The value of min and max must each be at least 1 (or we lift the values
	//   to 1).
	//
	// * Picking a passage should be evently distributed across the book.
	//
	// If we were to pick a passage at random and then pick a length at random
	// (as I did in the initial implementation), this will result in a bias
	// toward passages at the end of the book. To avoid this problem we do
	// the following to ensure that the passage is evenly distributed across
	// the verses in the book.
	//
	// 1. Pick a number of verses to include in the passage.
	//
	// 2. Pick a starting verse between index 0 and the length of the book minus
	//    the number of verses in the passage.
	//
	// This guarantees every passage selection is just as probably to be selected
	// as every other possible selection.

	// cap the minimum mn such that 1 <= mn <= mx <= len(b.Verses)
	// mx = mn if mx < mn
	mn = max(1, min(mn, len(verses)))
	mx = max(1, min(mx, len(verses)))
	mx = max(mx, mn)

	// pick a passage length
	var n int
	if mn == mx {
		n = mn
	} else {
		n = rand.Int()%(mx-mn) + mn //nolint:gosec // weak random is fine here
	}

	// pick a starting verse
	var x int
	if n >= len(verses) {
		x = 0
	} else {
		x = rand.Int() % (len(verses) - n) //nolint:gosec // weak random is fine here
	}
	y := x + n

	// return the selected verses
	vs := verses[x:y]
	return vs
}
