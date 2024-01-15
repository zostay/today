package ref

import (
	"fmt"
	"math/rand"
)

// RandomCanonical returns a random book of the Bible.
func RandomCanonical() *Book {
	return &Canonical.Books[rand.Int()%len(Canonical.Books)] //nolint:gosec // weak random is fine here
}

// RandomPassage returns a random passage from the given book of the Bible. It
// returns a two element slice of Verse, the first being the start of the
// passage and the second being the end of the passage. As of this writing, the
// verse references may be up to 30 Verses apart.
func RandomPassage(b *Book) []Verse {
	x := rand.Int() % len(b.Verses) //nolint:gosec // weak random is fine here
	o := rand.Int() % 30            //nolint:gosec // weak random is fine here
	y := x + o
	if y >= len(b.Verses) {
		y = len(b.Verses) - 1
	}

	return b.Verses[x:y]
}

func RandomPassageFromRef(b *Resolved) []Verse {
	x := rand.Int() % len(b.Verses()) //nolint:gosec // weak random is fine here
	o := rand.Int() % 30              //nolint:gosec // weak random is fine here
	y := x + o
	if y >= len(b.Verses()) {
		y = len(b.Verses()) - 1
	}

	return b.Verses()[x:y]
}

type randomOpts struct {
	category string
	book     string
	canon    *Canon
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

// Random pulls a random reference from the Bible and returns it. You can use the
// options to help narrow down where the passages are selected from.
func Random(opt ...RandomReferenceOption) (*Resolved, error) {
	o := &randomOpts{
		canon: Canonical,
	}
	for _, f := range opt {
		f(o)
	}

	var (
		b  *Book
		vs []Verse
	)

	if o.category != "" {
		_, hasCategory := o.canon.Categories[o.category]
		if !hasCategory {
			return nil, fmt.Errorf("unknown category: %s", o.category)
		}

		ps, err := o.canon.Category(o.category)
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
		vs = RandomPassageFromRef(be.Ref)
	} else {
		if o.book != "" {
			ex, err := Lookup(Canonical, o.book+" 1:1ffb", "")
			if err != nil {
				return nil, err
			}

			b = ex.Ref.Book
		} else {
			b = RandomCanonical()
		}

		vs = RandomPassage(b)
	}

	v1, v2 := vs[0], vs[len(vs)-1]

	return &Resolved{
		Book:  b,
		First: v1,
		Last:  v2,
	}, nil
}

// // RandomReference returns a random reference to a passage in the Bible in a
// // standard notation recognizable for American English speakers.
// //
// // This uses RandomCanonical and RandomPassage to generate the reference.
// //
// // Deprecated: Use Random() instead.
// func RandomReference() string {
// 	b := RandomCanonical()
// 	vs := RandomPassage(b)
//
// 	v1 := vs[0]
// 	v2 := vs[len(vs)-1]
//
// 	if len(vs) > 1 {
// 		return fmt.Sprintf("%s %s-%s", b.Name, v1.Ref(), v2.Ref())
// 	} else {
// 		return fmt.Sprintf("%s %s", b.Name, v1.Ref())
// 	}
// }
//
// // RandomVerse returns a random verse from the Bible. This uses RandomReference
// // to select a random reference and then uses GetVerse to retrieve the text of
// // the verse.
// func RandomVerse() (string, error) {
// 	ref := RandomReference()
// 	return GetVerse(ref)
// }
