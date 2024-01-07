package ref

// Book is a book of the Bible. We use this with a global map to do client-side
// verification of book names, chapter, and verse references.
type Book struct {
	Name      string
	JustVerse bool
	Verses    []Verse
}

type BookExtract struct {
	*Book
	First Verse
	Last  Verse
}

func (e *BookExtract) FullRef() string {
	return e.Book.Name + " " + e.First.Ref() + "-" + e.Last.Ref()
}

func (e *BookExtract) Verses() []Verse {
	verses := make([]Verse, 0, len(e.Book.Verses))
	started := false
	for _, verse := range e.Book.Verses {
		if verse == e.First {
			started = true
		}
		if started {
			verses = append(verses, verse)
		}
		if verse == e.Last {
			break
		}
	}

	return verses
}
