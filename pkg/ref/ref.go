package ref

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zostay/go-std/slices"
)

const (
	FurtherFollowing = -1
	ToEndOfChapter   = -2
	ToEndOfBook      = -3

	FollowingNotation        = "ff"
	FollowingChapterNotation = "ffc"
	FollowingBookNotation    = "ffb"
)

type Following int

const (
	FollowingNone Following = iota
	FollowingRemainingChapter
	FollowingRemainingBook
)

func ValidFollowing(n int) bool {
	return n == FurtherFollowing || n == ToEndOfChapter || n == ToEndOfBook
}

// Ref is any kind of Bible reference.
type Ref interface {
	// Ref returns the string representation of the reference.
	Ref() string

	// Validate returns whether the reference is valid.
	Validate() error
}

// Relative is any kind of reference that does not specify a Book.
type Relative interface {
	Ref

	// InBook turns this relative reference into a proper reference for the
	// given book.
	InBook(*Book) *Proper
}

// Verse is a reference to a verse of the Bible relative to a book.
type Verse interface {
	Relative
	Before(Verse) bool
	Equal(Verse) bool
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

// Validate returns true iff the chapter and verse are both positive.
func (v *ChapterVerse) Validate() error {
	if v.chapter <= 0 {
		return fmt.Errorf("chapter must be positive")
	}
	if v.verse <= 0 {
		return fmt.Errorf("verse must be positive")
	}
	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *ChapterVerse) InBook(b *Book) *Proper {
	return NewProper(b, v)
}

func (v *ChapterVerse) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case *ChapterVerse:
		return v.chapter < sv.chapter ||
			(v.chapter == sv.chapter && v.verse < sv.verse)
	case *JustVerse:
		return v.verse < sv.verse
	}
	panic("unable to validate unknown verse type")
}

func (v *ChapterVerse) Equal(ov Verse) bool {
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

// Validate returns true iff the verse is positive.
func (v *JustVerse) Validate() error {
	if v.verse <= 0 {
		return fmt.Errorf("verse must be positive")
	}
	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *JustVerse) InBook(b *Book) *Proper {
	return NewProper(b, v)
}

func (v *JustVerse) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case *ChapterVerse:
		return false
	case *JustVerse:
		return v.verse < sv.verse
	}
	panic("unable to validate unknown verse type")
}

func (v *JustVerse) Equal(ov Verse) bool {
	return v.verse == ov.(*JustVerse).verse
}

// AndFollowing takes a ref.Verse and attaches a notation indicating that the
// reference continues onward. Normally, this notation can be simply and
// informally "ff" which is a sort of generic reference to more verses after
// those listed. This could be understood to mean until the end of some pericope
// (which is itself a vaguely defined concept).
//
// I have formally extended it into three forms. The first is "ff" which we
// interpret to mean "and following". To that I add two formal forms: "ffc" and
// "ffb". The first means "on to the end of the chapter" and the second means "on
// to the end of the book".
type AndFollowing struct {
	Verse
	Following
}

func NewAndFollowing(verse Verse, following Following) *AndFollowing {
	return &AndFollowing{
		Verse:     verse,
		Following: following,
	}
}

func (v *AndFollowing) Ref() string {
	switch v.Following {
	case FollowingNone:
		return v.Verse.Ref()
	case FollowingRemainingChapter:
		return v.Verse.Ref() + FollowingChapterNotation
	case FollowingRemainingBook:
		return v.Verse.Ref() + FollowingBookNotation
	default:
		return v.Verse.Ref() + FollowingNotation
	}
}

func (v *AndFollowing) Validate() error {
	if !ValidFollowing(int(v.Following)) {
		return fmt.Errorf("invalid following notation: %d", v.Following)
	}
	return v.Verse.Validate()
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *AndFollowing) InBook(b *Book) *Proper {
	return NewProper(b, v)
}

// Range is a reference to a range of verses relative to a Book. It is formed of
// two ref.Verse references, which are the inclusive bounds of a relative Bible
// reference. The First ref.Verse may either be a ref.ChapterVerse or
// ref.JustVerse. The Last ref.Verse must be a ref.JustVerse if the first is
// ref.JustVerse. It may be a ref.JustVerse if the First is a ref.ChapterVerse.
// In either case, the given ref.Verse in Last must be greater than First.
type Range struct {
	First Verse
	Last  Verse
}

func (r *Range) Ref() string {
	if r.First == r.Last {
		return r.First.Ref()
	}
	return fmt.Sprintf("%s-%s", r.First.Ref(), r.Last.Ref())
}

func (r *Range) Validate() error {
	if r.First == nil || r.Last == nil {
		return fmt.Errorf("range is incorrect: both first and last are required")
	}

	var errs []error
	if err := r.First.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := r.Last.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("range is incorrect: %w", errors.Join(errs...))
	}

	_, isJvFirst := r.First.(*JustVerse)
	_, isJvLast := r.Last.(*JustVerse)
	if isJvFirst && !isJvLast {
		return fmt.Errorf("range is incorrect: first is verse-only but last is not")
	}

	if !r.First.Before(r.Last) {
		return fmt.Errorf("range is incorrect: first must be before last")
	}

	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (r *Range) InBook(b *Book) *Proper {
	return NewProper(b, r)
}

// Related is a lite of closely-related references (vaguely defined). These are
// rendered as a comma-separated list and allows for a more compact rendering of
// a set of references.
//
// For example: 3:16-18, 19-21, 12:14, 22, 23, 24
//
// Any verse-only references in a ref.Related reference must be preceded by a
// chapter and verse reference.
type Related struct {
	Refs []Ref
}

func (r *Related) Ref() string {
	return strings.Join(
		slices.Map(r.Refs, func(r Ref) string {
			return r.Ref()
		}), ", ")
}

func (r *Related) Validate() error {
	if len(r.Refs) == 0 {
		return fmt.Errorf("related list of references is incorrect: no references")
	}

	var cv *ChapterVerse
	for _, ref := range r.Refs {
		if ref == nil {
			return fmt.Errorf("related list of references is incorrect: nil reference")
		}

		if err := ref.Validate(); err != nil {
			return fmt.Errorf("related list of references is incorrect: %w", err)
		}

		switch r := ref.(type) {
		case *ChapterVerse:
			cv = r
		case *JustVerse:
			if cv == nil {
				return fmt.Errorf("related list of references is incorrect: verse-only reference must be preceded by a chapter and verse reference")
			}
		case *Range:
			if _, isJv := r.First.(*JustVerse); isJv && cv == nil {
				return fmt.Errorf("related list of references is incorrect: verse-only range reference must be preceded by a chapter and verse reference")
			}
		default:
			return fmt.Errorf("related list of references is incorrect: only chapter and verse, verse-only, or ranges are permitted in related reference lists")
		}
	}
	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (r *Related) InBook(b *Book) *Proper {
	return NewProper(b, r)
}

// Proper is some kind of reference that also includes the book that the verses
// or chapters and verses are relative to.
type Proper struct {
	*Book
	Verse Relative
}

func NewProper(book *Book, verse Relative) *Proper {
	return &Proper{
		Book:  book,
		Verse: verse,
	}
}

func (p *Proper) Ref() string {
	return fmt.Sprintf("%s %s", p.Book.Name, p.Verse.Ref())
}

func (p *Proper) Validate() error {
	if p.Book == nil {
		return fmt.Errorf("book is required")
	}
	if p.Verse == nil {
		return fmt.Errorf("verse is required")
	}

	return p.Verse.Validate()
}

// Multiple is a list of references to verses relative to a Book of the Bible. These
// rendered as a set of references separated by semi-colon. A List may not be the
// child of another List. All references in a list must be of the same type.
//
// The first references must a ref.Proper reference. The remaining references will
// be considered relative to that.
//
// For example: Genesis 3:15-18; 5:8; 12:10ff; 14:14-23; 15:1-6
type Multiple struct {
	Refs []Ref
}

func (m *Multiple) Ref() string {
	return strings.Join(
		slices.Map(m.Refs, func(r Ref) string {
			return r.Ref()
		}), "; ")
}

func (m *Multiple) Validate() error {
	if len(m.Refs) == 0 {
		return fmt.Errorf("multiple list of references is incorrect: no references")
	}

	if _, isProper := m.Refs[0].(*Proper); !isProper {
		return fmt.Errorf("multiple list of references is incorrect: first reference must be a proper reference")
	}

	return nil
}

var _ Verse = (*ChapterVerse)(nil)
var _ Verse = (*JustVerse)(nil)
var _ Relative = (*AndFollowing)(nil)
var _ Relative = (*Range)(nil)
var _ Relative = (*Related)(nil)
var _ Ref = (*Proper)(nil)
var _ Ref = (*Multiple)(nil)
