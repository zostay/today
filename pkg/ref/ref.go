package ref

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zostay/go-std/slices"

	"github.com/zostay/today/internal/options"
	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/canon"
)

const (
	FollowingNotation        = "ff"
	FollowingChapterNotation = "ffc"
	FollowingBookNotation    = "ffb"
)

type Following int

const (
	FollowingRemainingChapter Following = iota
	FollowingRemainingBook
)

func ValidFollowing(n Following) bool {
	return n == FollowingRemainingChapter ||
		n == FollowingRemainingBook
}

type (
	Ref      = canon.Ref
	Verse    = canon.Verse
	Absolute = canon.Absolute
	Relative = canon.Relative
)

// CV is a reference to a specific chapter and verse for books with
// both chapters and Verses.
type CV struct {
	Chapter int
	Verse   int
}

// Ref returns the Chapter:Verse reference string.
func (v CV) Ref() string {
	return strconv.Itoa(v.Chapter) + ":" + strconv.Itoa(v.Verse)
}

// InBook turns this relative reference into an absolute reference for the given
// book.
func (v CV) InBook(b string) Absolute {
	return NewProper(b, &Single{Verse: v})
}

func (v CV) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case CV:
		return v.Chapter < sv.Chapter ||
			(v.Chapter == sv.Chapter && v.Verse < sv.Verse)
	default:
		return v.Before(sv.RelativeTo(v))
	}
}

func (v CV) Equal(ov Verse) bool {
	switch sv := ov.(type) {
	case CV:
		return v.Chapter == sv.Chapter && v.Verse == sv.Verse
	default:
		return v.Equal(sv.RelativeTo(v))
	}
}

func (v CV) RelativeTo(ov Verse) Verse {
	switch ov.(type) {
	case CV:
		return v
	case N:
		return N{Number: v.Chapter}
	}
	panic("unable to convert CV to unknown verse type")
}

// N is a reference to either a specific verse or to a chapter. It represents
// any case where a single number is used in a reference. Whether this number
// represents a verse or chapter is determined by the context.
//
// For example:
//
//	Philemon 12 - this is a verse because Philemon has no chapters
//	Isaiah 24 - this is a chapter because Isaiah has chapters
//	John 3:16-17 - the 17 is a verse, the chapter is implied by the previous chapter/verse reference in the range
//
// N represents all the of the above cases.
type N struct {
	Number int
}

// Ref returns the Verse reference string (no Chapter:).
func (n N) Ref() string {
	return strconv.Itoa(n.Number)
}

// Validate returns true iff the verse is positive.
func (n N) Validate() error {
	if n.Number <= 0 {
		return invalid("chapter or verse must be positive")
	}
	return nil
}

// InBook turns this relative reference into an absolute reference for the given
// book.
func (n N) InBook(b string) Absolute {
	return NewProper(b, &Single{Verse: n})
}

func (n N) Before(ov Verse) bool {
	switch sv := ov.(type) {
	case N:
		return n.Number < sv.Number
	default:
		return n.Before(sv.RelativeTo(n))
	}
}

func (n N) Equal(ov Verse) bool {
	switch sv := ov.(type) {
	case N:
		return n.Number == sv.Number
	default:
		return n.Equal(sv.RelativeTo(n))
	}
}

func (n N) RelativeTo(ov Verse) Verse {
	switch sv := ov.(type) {
	case N:
		return n
	case CV:
		return CV{Chapter: sv.Chapter, Verse: n.Number}
	}
	panic("unable to convert N to unknown verse type")
}

// Single is a relative reference to a single verse. It wraps a single verse.
type Single struct {
	Verse Verse
}

func (v *Single) Ref() string {
	return v.Verse.Ref()
}

func (v *Single) Validate() error {
	if v.Verse == nil {
		return errors.New("verse is required")
	}
	return v.Verse.Validate()
}

// InBook turns this relative reference into an absolute reference for the given
// book.
func (v *Single) InBook(b string) Absolute {
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

func (v *AndFollowing) Ref() string {
	switch v.Following { //nolint:exhaustive // we don't need to handle all cases
	case FollowingRemainingBook:
		if v.Verse.Ref() == "1" || v.Verse.Ref() == "1:1" {
			return ""
		}
		return v.Verse.Ref() + FollowingBookNotation
	default:
		return v.Verse.Ref() + FollowingNotation
	}
}

func (v *AndFollowing) Validate() error {
	if !ValidFollowing(v.Following) {
		return invalid("invalid following notation")
	}
	if v.Verse == nil {
		return errors.New("verse is required")
	}
	return v.Verse.Validate()
}

// InBook turns this relative reference into a proper reference for the given
// book.
func (v *AndFollowing) InBook(b string) Absolute {
	return NewProper(b, v)
}

// Range is a reference to a range of verses relative to a Book. It is formed of
// two ref.Verse references, which are the inclusive bounds of a relative Bible
// reference. The First ref.Verse may either be a ref.CV or
// ref.N. The Last ref.Verse must be a ref.N if the first is
// ref.N. It may be a ref.N if the First is a ref.CV.
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

	_, isJvFirst := r.First.(N)
	_, isJvLast := r.Last.(N)
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
func (r *Range) InBook(b string) Absolute {
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

	first := true
	isJV := false
	for _, ref := range r.Refs {
		if ref == nil {
			return invalid("related list of references is incorrect: nil reference")
		}

		if _, isRelated := ref.(*Related); isRelated {
			return invalid("related list of references may not contain a nested list of related references")
		}

		if err := ref.Validate(); err != nil {
			return invalid("related list of references is incorrect: %w", unravelInvalid(err))
		}

		if first {
			isJV = !strings.Contains(r.Refs[0].Ref(), ":")
			first = false
			continue
		}

		if isJV && strings.Contains(ref.Ref(), ":") {
			return invalid("related list of references is incorrect: related reference list starts with verse-only reference, but contains a chapter-verse reference")
		}
	}

	return nil
}

// InBook turns this relative reference into an absolute reference for the given
// book.
func (r *Related) InBook(b string) Absolute {
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

func (p *Proper) FullNameRef(opt ...ResolveOption) (string, error) {
	o := options.MakeResolveOpts(opt)
	abbrs := o.Abbreviations

	if abbrs == nil {
		return p.Ref(), nil
	}

	fullName, err := abbrs.BookName(p.Book)
	if err != nil {
		return "", err
	}

	if o.Singular {
		fullName, err = abbrs.SingularName(fullName)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s %s", fullName, p.Verse.Ref()), nil
}

func (p *Proper) AbbreviatedRef(opt ...ResolveOption) (string, error) {
	o := options.MakeResolveOpts(opt)
	abbrs := o.Abbreviations

	if abbrs == nil {
		return p.Ref(), nil
	}

	fullName, err := abbrs.BookName(p.Book)
	if err != nil {
		return "", err
	}

	abbrName, err := abbrs.PreferredAbbreviation(fullName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", abbrName, p.Verse.Ref()), nil
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

func (p *Proper) IsSingleRange() bool {
	_, isSingle := p.Verse.(*Single)
	_, isRange := p.Verse.(*Range)
	_, isAndFollowing := p.Verse.(*AndFollowing)
	return isSingle || isRange || isAndFollowing
}

func (p *Proper) InCanon(c *bible.Canon, opts *bible.Resolve) ([]bible.CanonicalRef, error) {
	b, err := Book(c, p.Book, opts)
}

// Multiple is a list of references to verses relative to a Book of the Bible. These
// rendered as a set of references separated by semi-colon. A List may not be the
// child of another List. All references in a list must be of the same type.
//
// The first references must a ref.Proper reference. The remaining references will
// be considered relative to that.
//
// For example:
//
//	Genesis 3:15-18; 5:8; 12:10ff; 14:14-23; 15:1-6
type Multiple struct {
	Refs []Ref
}

func (m *Multiple) Ref() string {
	return strings.Join(
		slices.Map(m.Refs, func(r Ref) string {
			return r.Ref()
		}), "; ")
}

func (m *Multiple) FullNameRef(opt ...ResolveOption) (string, error) {
	out := &bytes.Buffer{}
	needSemi := false
	for _, ref := range m.Refs {
		if needSemi {
			fmt.Fprint(out, "; ")
		}

		if abs, isAbs := ref.(Absolute); isAbs {
			ref, err := abs.FullNameRef(opt...)
			if err != nil {
				return "", err
			}

			fmt.Fprint(out, ref)
			continue
		}

		fmt.Fprint(out, ref.Ref())
	}
	return out.String(), nil
}

func (m *Multiple) AbbreviatedRef(opt ...ResolveOption) (string, error) {
	out := &bytes.Buffer{}
	needSemi := false
	for _, ref := range m.Refs {
		if needSemi {
			fmt.Fprint(out, "; ")
		}

		if abs, isAbs := ref.(Absolute); isAbs {
			ref, err := abs.AbbreviatedRef(opt...)
			if err != nil {
				return "", err
			}

			fmt.Fprint(out, ref)
			continue
		}

		fmt.Fprint(out, ref.Ref())
	}
	return out.String(), nil
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

func (m *Multiple) IsSingleRange() bool {
	if len(m.Refs) != 1 {
		return false
	}

	if p, isProper := m.Refs[0].(*Proper); isProper {
		_, isSingle := p.Verse.(*Single)
		_, isRange := p.Verse.(*Range)
		_, isAndFollowing := p.Verse.(*AndFollowing)
		return isSingle || isRange || isAndFollowing
	}

	return false
}

func (m *Multiple) InCanon(c *bible.Canon, opts *bible.Resolve) ([]bible.CanonicalRef, error) {
	var rs []bible.CanonicalRef
	var b *bible.Book
	for i := range m.Refs {
		switch r := m.Refs[i].(type) {
		case Absolute:
			refs, err := r.InCanon(c, opts)
			if err != nil {
				return nil, err
			}

			rs = append(rs, refs...)
		default:
			thisRs, err := r.InBook(b.Name).InCanon(c, opts)
			if err != nil {
				return nil, err
			}

			rs = append(rs, thisRs...)
		}
	}

	return rs, nil
}

var _ Verse = CV{}
var _ Verse = N{}
var _ Relative = (*Single)(nil)
var _ Relative = (*AndFollowing)(nil)
var _ Relative = (*Range)(nil)
var _ Relative = (*Related)(nil)
var _ Absolute = (*Proper)(nil)
var _ Absolute = (*Multiple)(nil)
var _ Absolute = (*Resolved)(nil)
