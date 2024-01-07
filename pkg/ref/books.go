package ref

import (
	"strconv"
)

const Final = -1

// Book is a book of the Bible. We use this with a global map to do client-side
// verification of book names, chapter, and verse references.
type Book struct {
	Name      string
	JustVerse bool
	Verses    []VerseRef
}

type WildcardType int

const (
	WildcardNone WildcardType = iota
	WildcardChapter
	WildcardVerse
)

// VerseRef is an object that can return a Bible reference. Since not all books
// of the Bible have both chapters and Verses, this allows us to handle both.
type VerseRef interface {
	Ref() string
	Wildcard() WildcardType
	Before(VerseRef) bool
	Equal(VerseRef) bool
}

// ChapterVerse is a reference to a specific chapter and verse for books with
// both chapters and Verses.
type ChapterVerse struct {
	chapter int
	verse   int
}

func NewChapterVerse(chapter, verse int) *ChapterVerse {
	return &ChapterVerse{
		chapter: chapter,
		verse:   verse,
	}
}

// Ref returns the Chapter:Verse reference string.
func (v *ChapterVerse) Ref() string {
	return strconv.Itoa(v.chapter) + ":" + strconv.Itoa(v.verse)
}

func (v *ChapterVerse) Wildcard() WildcardType {
	if v.chapter == -1 {
		return WildcardChapter
	}
	if v.verse == -1 {
		return WildcardVerse
	}
	return WildcardNone
}

func (v *ChapterVerse) Before(ov VerseRef) bool {
	return v.chapter < ov.(*ChapterVerse).chapter || (v.chapter == ov.(*ChapterVerse).chapter && v.verse < ov.(*ChapterVerse).verse)
}

func (v *ChapterVerse) Equal(ov VerseRef) bool {
	return v.chapter == ov.(*ChapterVerse).chapter && v.verse == ov.(*ChapterVerse).verse
}

// JustVerse is a reference to a specific verse for books without chapters.
type JustVerse struct {
	verse int
}

func NewJustVerse(verse int) *JustVerse {
	return &JustVerse{
		verse: verse,
	}
}

// Ref returns the Verse reference string (no Chapter:).
func (v *JustVerse) Ref() string {
	return strconv.Itoa(v.verse)
}

func (v *JustVerse) Wildcard() WildcardType {
	if v.verse == -1 {
		return WildcardVerse
	}
	return WildcardNone
}

func (v *JustVerse) Before(ov VerseRef) bool {
	return v.verse < ov.(*JustVerse).verse
}

func (v *JustVerse) Equal(ov VerseRef) bool {
	return v.verse == ov.(*JustVerse).verse
}

type BookExtract struct {
	*Book
	First VerseRef
	Last  VerseRef
}

func (e *BookExtract) FullRef() string {
	return e.Book.Name + " " + e.First.Ref() + "-" + e.Last.Ref()
}

func (e *BookExtract) Verses() []VerseRef {
	verses := make([]VerseRef, 0, len(e.Book.Verses))
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
