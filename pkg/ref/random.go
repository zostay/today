package ref

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/agnivade/levenshtein"

	"github.com/zostay/today/internal/options"
	"github.com/zostay/today/pkg/bible"
)

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

type RandomReferenceOption = options.RandomReferenceOption

var (
	// FromCanon limits the random selection to a specific canon of the Bible.
	FromCanon = options.FromCanon

	// FromBook limits the random selection to a specific book of the Bible.
	FromBook = options.FromBook

	// FromCategory limits the random selection to a specific category of books.
	FromCategory = options.FromCategory

	// WithAtLeast ensures the random selection has at least the given number of verses.
	WithAtLeast = options.WithAtLeast

	// WithAtMost ensures the random selection has at most the given number of verses.
	WithAtMost = options.WithAtMost
)

// Random pulls a random reference from the Bible and returns it. You can use the
// options to help narrow down where the passages are selected from.
func Random(opt ...RandomReferenceOption) (*Resolved, error) {
	o := options.MakeRandomReferenceOpts(opt)

	var (
		b  *Book
		vs []Verse
	)

	if o.Category != "" {
		_, hasCategory := o.Canon.Categories[o.Category]
		if !hasCategory {
			var possibilities []string
			for cat := range o.Canon.Categories {
				if levenshtein.ComputeDistance(o.Category, cat) <= 4 {
					possibilities = append(possibilities, cat)
				}
			}

			return nil, &UnknownCategoryError{
				Category:      o.Category,
				Possibilities: possibilities,
			}
		}

		ps, err := o.Canon.Category(o.Category)
		if err != nil {
			return nil, err
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
		vs = RandomPassageFromRef(be.Ref, o.Min, o.Max)
	} else {
		if o.Book != "" {
			ex, err := Lookup(bible.Protestant, o.Book+" 1:1ffb", "")
			if err != nil {
				return nil, err
			}

			b = ex.Ref.Book
		} else {
			b = RandomCanonical(o.Canon)
		}

		vs = RandomPassage(b, o.Min, o.Max)
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
