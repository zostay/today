package ref

import (
	"fmt"
	"math/rand"
)

// RandomCanonical returns a random book of the Bible.
func RandomCanonical() *Book {
	return &Canonical[rand.Int()%len(Canonical)] //nolint:gosec // weak random is fine here
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

func RandomPassageFromExtract(b *Resolved) []Verse {
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
}

type RandomReferenceOption func(*randomOpts)

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
func Random(opt ...RandomReferenceOption) (string, error) {
	o := &randomOpts{}
	for _, f := range opt {
		f(o)
	}

	var (
		b  *Book
		vs []Verse
	)
	// if o.category != "" {
	// 	exs, err := LookupCategory(o.category)
	// 	if err != nil {
	// 		return "", err
	// 	}
	//
	// 	// lazy way to weight the books by the number of verses they have
	// 	bag := make([]Resolved, 0, len(exs))
	// 	for i := range exs {
	// 		for range exs[i].Verses() {
	// 			bag = append(bag, exs[i])
	// 		}
	// 	}
	//
	// 	be := bag[rand.Int()%len(bag)] //nolint:gosec // weak random is fine here
	// 	b = be.Book
	// 	vs = RandomPassageFromExtract(&be)
	// } else {
	if o.book != "" {
		ex, err := Lookup(Canonical, o.book+" 1:1ffb", "")
		if err != nil {
			return "", err
		}

		b = ex.Ref.Book
	} else {
		b = RandomCanonical()
	}

	vs = RandomPassage(b)
	//}

	v1, v2 := vs[0], vs[len(vs)-1]

	if len(vs) > 1 {
		return fmt.Sprintf("%s %s-%s", b.Name, v1.Ref(), v2.Ref()), nil
	}

	return fmt.Sprintf("%s %s", b.Name, v1.Ref()), nil
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
