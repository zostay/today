package ref

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zostay/go-std/slices"
)

const (
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

type ValidationError struct {
	Cause error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.Cause)
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

func invalid(msg string, args ...any) error {
	err := fmt.Errorf(msg, args...)
	return &ValidationError{Cause: err}
}

func unravelInvalid(err error) error {
	if err == nil {
		return nil
	}
	if verr, isInvalid := err.(*ValidationError); isInvalid {
		return verr.Cause
	}
	return err
}

func ValidFollowing(n Following) bool {
	return n == FollowingNone ||
		n == FollowingRemainingChapter ||
		n == FollowingRemainingBook
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
	InBook(string) *Proper
}

// Verse is a reference to a verse of the Bible relative to a book.
type Verse interface {
	Relative
	Before(Verse) bool
	Equal(Verse) bool
}

// Absolute is any kind of reference that specifies a Book.
type Absolute interface {
	Ref

	// Names returns the names of the books referenced.
	Names() []string
}

// CV is a reference to a specific chapter and verse for books with
// both chapters and Verses.
type CV struct {
	Chapter int
	Verse   int
}

// Ref returns the Chapter:Verse reference string.
func (v *CV) Ref() string {
	return strconv.Itoa(v.Chapter) + ":" + strconv.Itoa(v.Verse)
}

// Validate returns true iff the chapter and verse are both positive.
func (v *CV) Validate() error {
	if v.Chapter <= 0 {
		return invalid("chapter must be positive")
	}
	if v.Verse <= 0 {
		return invalid("verse must be positive")
	}
	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *CV) InBook(b string) *Proper {
	return NewProper(b, v)
}

func (v *CV) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case *CV:
		return v.Chapter < sv.Chapter ||
			(v.Chapter == sv.Chapter && v.Verse < sv.Verse)
	case *V:
		return v.Verse < sv.Verse
	}
	panic("unable to validate unknown verse type")
}

func (v *CV) Equal(ov Verse) bool {
	return v.Chapter == ov.(*CV).Chapter && v.Verse == ov.(*CV).Verse
}

// V is a reference to a specific verse for books without chapters.
type V struct {
	Verse int
}

// Ref returns the Verse reference string (no Chapter:).
func (v *V) Ref() string {
	return strconv.Itoa(v.Verse)
}

// Validate returns true iff the verse is positive.
func (v *V) Validate() error {
	if v.Verse <= 0 {
		return invalid("verse must be positive")
	}
	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *V) InBook(b string) *Proper {
	return NewProper(b, v)
}

func (v *V) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case *CV:
		return false
	case *V:
		return v.Verse < sv.Verse
	}
	panic("unable to validate unknown verse type")
}

func (v *V) Equal(ov Verse) bool {
	return v.Verse == ov.(*V).Verse
}

// Single is a relative reference to a single verse. It wraps a single verse.
type Single struct {
	Verse Verse
}

func NewSingle(verse Verse) *Single {
	return &Single{
		Verse: verse,
	}
}

func (v *Single) Ref() string {
	return v.Verse.Ref()
}

func (v *Single) Validate() error {
	return v.Verse.Validate()
}

func (v *Single) InBook(b string) *Proper {
	return NewProper(b, v)
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
	Verse     Verse
	Following Following
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
	if !ValidFollowing(v.Following) {
		return invalid("invalid following notation")
	}
	return v.Verse.Validate()
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *AndFollowing) InBook(b string) *Proper {
	return NewProper(b, v)
}

// Range is a reference to a range of verses relative to a Book. It is formed of
// two ref.Verse references, which are the inclusive bounds of a relative Bible
// reference. The First ref.Verse may either be a ref.CV or
// ref.V. The Last ref.Verse must be a ref.V if the first is
// ref.V. It may be a ref.V if the First is a ref.CV.
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
		return invalid("range is incorrect: both first and last are required")
	}

	var errs []error
	if err := r.First.Validate(); err != nil {
		errs = append(errs, unravelInvalid(err))
	}
	if err := r.Last.Validate(); err != nil {
		errs = append(errs, unravelInvalid(err))
	}

	if len(errs) > 0 {
		return invalid("range is incorrect: %w", errors.Join(errs...))
	}

	_, isJvFirst := r.First.(*V)
	_, isJvLast := r.Last.(*V)
	if isJvFirst && !isJvLast {
		return invalid("range is incorrect: first is verse-only but last is not")
	}

	if !r.First.Before(r.Last) {
		return invalid("range is incorrect: first must be before last")
	}

	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (r *Range) InBook(b string) *Proper {
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
	Refs []Relative
}

func (r *Related) Ref() string {
	return strings.Join(
		slices.Map(r.Refs, func(r Relative) string {
			return r.Ref()
		}), ", ")
}

func (r *Related) Validate() error {
	if len(r.Refs) == 0 {
		return invalid("related list of references is incorrect: no references")
	}

	isJV := !strings.Contains(r.Refs[0].Ref(), ":")
	for _, ref := range r.Refs {
		if ref == nil {
			return invalid("related list of references is incorrect: nil reference")
		}

		if err := ref.Validate(); err != nil {
			return invalid("related list of references is incorrect: %w", unravelInvalid(err))
		}

		if isJV && strings.Contains(ref.Ref(), ":") {
			return invalid("related list of references is incorrect: related reference list starts with verse-only reference, but contains a chapter-verse reference")
		}
	}

	return nil
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (r *Related) InBook(b string) *Proper {
	return NewProper(b, r)
}

// Proper is some kind of reference that also includes the book that the verses
// or chapters and verses are relative to.
type Proper struct {
	Book  string
	Verse Relative
}

func NewProper(book string, verse Relative) *Proper {
	return &Proper{
		Book:  book,
		Verse: verse,
	}
}

func (p *Proper) Ref() string {
	return fmt.Sprintf("%s %s", p.Book, p.Verse.Ref())
}

func (p *Proper) Validate() error {
	if p.Book == "" {
		return fmt.Errorf("book is required")
	}
	if p.Verse == nil {
		return fmt.Errorf("verse is required")
	}

	return p.Verse.Validate()
}

func (p *Proper) Names() []string {
	return []string{p.Book}
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
		return invalid("multiple list of references is incorrect: no references")
	}

	if _, isProper := m.Refs[0].(*Proper); !isProper {
		return invalid("multiple list of references is incorrect: first reference must be a proper reference")
	}

	for i := range m.Refs {
		switch m.Refs[i].(type) {
		case Relative:
		case *Proper:
		default:
			return fmt.Errorf("multiple list of references is incorrect: only relative or proper references are permitted in multiple reference lists")
		}

		if err := m.Refs[i].Validate(); err != nil {
			return fmt.Errorf("multiple list of references is incorrect: %w", unravelInvalid(err))
		}
	}

	return nil
}

func (m *Multiple) Names() []string {
	var names []string
	for i := range m.Refs {
		if p, isProper := m.Refs[i].(*Proper); isProper {
			names = append(names, p.Book)
		}
	}
	return names
}

// Resolved is a normalized reference to a single range of verses in a single
// book, which may have a length of one. Both Verse references are inclusive and
// must match the verse type of the book. (I.e., if the book has chapters, then
// both First and Last must be ref.CV references.)
type Resolved struct {
	Book  *Book
	First Verse
	Last  Verse
}

func (r *Resolved) Ref() string {
	if r.First.Equal(r.Last) {
		return fmt.Sprintf("%s %s", r.Book.Name, r.First.Ref())
	}
	return fmt.Sprintf("%s %s-%s", r.Book.Name, r.First.Ref(), r.Last.Ref())
}

func (r *Resolved) Validate() error {
	if r.Book == nil {
		return fmt.Errorf("book is required")
	}
	if r.First == nil {
		return fmt.Errorf("first reference is required")
	}
	if r.Last == nil {
		return fmt.Errorf("last reference is required")
	}

	if err := r.First.Validate(); err != nil {
		return invalid("first reference is invalid: %w", unravelInvalid(err))
	}
	if err := r.Last.Validate(); err != nil {
		return invalid("last reference is invalid: %w", unravelInvalid(err))
	}

	if r.Last.Before(r.First) {
		return fmt.Errorf("first reference must be before or equal to last reference")
	}

	return nil
}

func (r *Resolved) Names() []string {
	return []string{r.Book.Name}
}

func (r *Resolved) Verses() []Verse {
	verses := make([]Verse, 0, len(r.Book.Verses))
	started := false
	for _, verse := range r.Book.Verses {
		if verse.Equal(r.First) {
			started = true
		}
		if started {
			verses = append(verses, verse)
		}
		if verse.Equal(r.Last) {
			break
		}
	}
	return verses
}

var _ Verse = (*CV)(nil)
var _ Verse = (*V)(nil)
var _ Relative = (*Single)(nil)
var _ Relative = (*AndFollowing)(nil)
var _ Relative = (*Range)(nil)
var _ Relative = (*Related)(nil)
var _ Absolute = (*Proper)(nil)
var _ Absolute = (*Multiple)(nil)
var _ Absolute = (*Resolved)(nil)
