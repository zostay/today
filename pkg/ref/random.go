package ref

import (
	"math/rand"
)

// RandomCanonical returns a random book of the Bible.
func RandomCanonical() *Book {
	return &Canonical[rand.Int()%len(Canonical)] //nolint:gosec // weak random is fine here
}

// RandomPassage returns a random passage from the given book of the Bible. It
// returns a two element slice of VerseRef, the first being the start of the
// passage and the second being the end of the passage. As of this writing, the
// verse references may be up to 30 Verses apart.
func RandomPassage(b *Book) []VerseRef {
	x := rand.Int() % len(b.Verses) //nolint:gosec // weak random is fine here
	o := rand.Int() % 30            //nolint:gosec // weak random is fine here
	y := x + o
	if y >= len(b.Verses) {
		y = len(b.Verses) - 1
	}

	return b.Verses[x:y]
}

func RandomPassageFromExtract(b *BookExtract) []VerseRef {
	x := rand.Int() % len(b.Verses()) //nolint:gosec // weak random is fine here
	o := rand.Int() % 30              //nolint:gosec // weak random is fine here
	y := x + o
	if y >= len(b.Verses()) {
		y = len(b.Verses()) - 1
	}

	return b.Verses()[x:y]
}

//// RandomReference returns a random reference to a passage in the Bible in a
//// standard notation recognizable for American English speakers.
////
//// This uses RandomCanonical and RandomPassage to generate the reference.
////
//// Deprecated: Use Random() instead.
//func RandomReference() string {
//	b := RandomCanonical()
//	vs := RandomPassage(b)
//
//	v1 := vs[0]
//	v2 := vs[len(vs)-1]
//
//	if len(vs) > 1 {
//		return fmt.Sprintf("%s %s-%s", b.Name, v1.Ref(), v2.Ref())
//	} else {
//		return fmt.Sprintf("%s %s", b.Name, v1.Ref())
//	}
//}
//
//// RandomVerse returns a random verse from the Bible. This uses RandomReference
//// to select a random reference and then uses GetVerse to retrieve the text of
//// the verse.
//func RandomVerse() (string, error) {
//	ref := RandomReference()
//	return GetVerse(ref)
//}
