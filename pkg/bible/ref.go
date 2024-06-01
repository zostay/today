package bible

import "fmt"

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
	if verr, isInvalid := err.(*ValidationError); isInvalid { //nolint:errorlint // I want to unwrap only the topmost here...
		return verr.Cause
	}
	return err
}

// Ref is any kind of Bible reference.
type Ref interface {
	// Ref returns the string representation of the reference.
	Ref() string

	// Validate returns whether the reference is valid.
	Validate() error
}

// RelativeRef is any kind of reference that does not specify a Book.
type RelativeRef interface {
	Ref

	// InBook turns this relative reference into an absolute reference for the
	// given book.
	InBook(string, *Resolve) AbsoluteRef
}

// VerseRef is a reference to a verse of the Bible relative to a book.
type VerseRef interface {
	RelativeRef

	// Before returns true if this verse is before the given verse. If the given
	// verse is not the same type as this verse, it must be converted to this
	// type using RelativeTo first.
	Before(VerseRef) bool

	// Equal returns true if this verse is equal to the given verse. If the given
	// verse is not the same type as this verse, it must be converted to this
	// type using RelativeTo first.
	Equal(VerseRef) bool

	// RelativeTo returns the given verse converted to this type of verse using
	// this verse as context.
	RelativeTo(VerseRef) VerseRef
}

// AbsoluteRef is any kind of reference that specifies a Book.
type AbsoluteRef interface {
	Ref

	// Names returns the names of the books referenced.
	Names() []string

	// AbbreviatedRef returns the reference, but with the book name using a
	// standard abbreviation. The WithAbbreviations option may be passed to
	// control the abbreviations to be used. This method does not ensure that
	// the verse can be resolved and represents a valid reference. The argument
	// is the options object to use, which should at least implement the
	// ResolveOptions interface.
	AbbreviatedRef(any) (string, error)

	// FullNameRef returns the reference, but ensures abbreviations have been
	// expanded to the full book name. The WithAbbreviations option may be passed
	// to control the abbreviations to be used. This method does not ensure that
	// the verse can be resolved and represents a valid reference. The argument
	// is the options object to use, which should at least implement the
	// ResolveOptions interface.
	FullNameRef(any) (string, error)

	// IsSingleRange returns true if the reference is a single range of verses.
	IsSingleRange() bool

	// InCanon return a slice of CanonicalRef objects for thsi reference mapped
	// into that canon. This does not have to be fully validated, but the book
	// name, at least, will have to match a book in that canon. The final
	// argument is the options object to use, which should at least implement
	// the ResolveOptions interface.
	InCanon(*Canon, *Resolve) ([]CanonicalRef, error)
}

// CanonicalRef is a normalized reference to a single range of verses in a single
// book, which may have a length of one. Both Verse references are inclusive and
// must match the verse type of the book. (I.e., if the book has chapters, then
// both First and Last must be ref.CV references.)
type CanonicalRef struct {
	Book  *Book
	First V
	Last  V
}

func (r *CanonicalRef) Ref() string {
	if r.First.Equal(r.Last) {
		return fmt.Sprintf("%s %s", r.Book.Name, r.First.Ref())
	}
	return fmt.Sprintf("%s %s-%s", r.Book.Name, r.First.Ref(), r.Last.Ref())
}

func (r *CanonicalRef) Validate() error {
	if r.Book == nil {
		return fmt.Errorf("book is required")
	}
	if r.First.V == 0 {
		return fmt.Errorf("first reference is required")
	}
	if r.Last.V == 0 {
		return fmt.Errorf("last reference is required")
	}

	if err := r.First.Validate(r.Book.JustVerse); err != nil {
		return invalid("first reference is invalid: %w", unravelInvalid(err))
	}
	if err := r.Last.Validate(r.Book.JustVerse); err != nil {
		return invalid("last reference is invalid: %w", unravelInvalid(err))
	}

	if r.Book.JustVerse && r.First.C != 0 {
		return invalid("book has no chapters, but first reference is a chapter and verse reference")
	}
	if r.Book.JustVerse && r.Last.C != 0 {
		return invalid("book has no chapters, but last reference is a chapter and verse reference")
	}
	if !r.Book.JustVerse && r.First.C == 0 {
		return invalid("book has chapters, but first reference is not a chapter and verse reference")
	}
	if !r.Book.JustVerse && r.Last.C == 0 {
		return invalid("book has chapters, but last reference is not a chapter and verse reference")
	}

	if r.Last.RelativeTo(r.First).Before(r.First) {
		return fmt.Errorf("first reference must be before or equal to last reference")
	}

	return nil
}

func (r *CanonicalRef) Names() []string {
	return []string{r.Book.Name}
}

func (r *CanonicalRef) IsSingleRange() bool {
	return true
}

func (r *CanonicalRef) Verses() []V {
	verses := make([]V, 0, len(r.Book.Verses))
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

func (r *CanonicalRef) compactRef(name string) (string, error) {
	if r.First.Equal(r.Last) {
		return fmt.Sprintf("%s %s", name, r.First.Ref()), nil
	}

	if r.First.Equal(r.Book.Verses[0]) && r.Last.Equal(r.Book.Verses[len(r.Book.Verses)-1]) {
		return name, nil
	}

	isFCV := r.First.C != 0
	isLCV := r.Last.C != 0
	if isFCV && isLCV {
		if r.First.C == r.Last.C {
			lvInC, err := r.Book.LastVerseInChapter(r.First.C)
			if err != nil {
				return "", err
			}

			if r.First.V == 1 && r.Last.V == lvInC {
				return fmt.Sprintf("%s %d", name, r.First.C), nil
			}

			return fmt.Sprintf("%s %d:%d-%d", name, r.First.C, r.First.V, r.Last.V), nil
		} else {
			lvInC, err := r.Book.LastVerseInChapter(r.Last.C)
			if err != nil {
				return "", err
			}

			if r.First.V == 1 && r.Last.V == lvInC {
				return fmt.Sprintf("%s %d-%d", name, r.First.C, r.Last.C), nil
			}
		}
	}

	return fmt.Sprintf("%s %s-%s", name, r.First.Ref(), r.Last.Ref()), nil
}

// CompactRef returns a compact representation of the resolved reference. If the
// reference is a single verse, only one verse is returned (e.g., Genesis 12:4).
// If the reference is limited to a single chapter, the chapter is mentioned only
// once (e.g., Genesis 12:4-6). If the reference is for an entire chapter, the
// verse part is omitted (e.g., Genesis 12). If the reference is for an entire
// book, only the book name is returned (e.g., Genesis).
//
// This will ignore the AsSingleChapter option as a resoled reference inherently
// knows if it is a single chapter via IsSingleChapter.
func (r *CanonicalRef) CompactRef(opt ...ResolveOption) (string, error) {
	o := makeResolveOpts(opt)
	name := r.Book.Name
	if o.Abbreviations != nil && r.IsSingleChapter() {
		var err error
		name, err = o.Abbreviations.SingularName(name)
		if err != nil {
			return "", err
		}
	}
	return r.compactRef(name)
}

// FullNameRef is a synonym for CompactRef.
//
// This will ignore the AsSingleChapter option as a resoled reference inherently
// knows if it is a single chapter via IsSingleChapter.
func (r *CanonicalRef) FullNameRef(opt ...ResolveOption) (string, error) {
	return r.CompactRef(opt...)
}

// IsSingleChapter returns true if the reference is for a single chapter. This is
// decided by considering the kind of book this is and comparing First to Last.
func (r *CanonicalRef) IsSingleChapter() bool {
	if r.Book.JustVerse {
		return true
	}

	isFCV := r.First.C != 0
	isLCV := r.Last.C != 0
	if isFCV && isLCV {
		return r.First.C == r.Last.C
	}

	return false
}

// AbbreviatedRef returns a compact and abbreviated representation of the
// resolved reference. This works the same as CompactRef, but with the book name
// abbreviated using the Standard abbreviation for the book. You may use
// ref.WithAbbreviations to select an alternate set of abbreviations. If this
// option is not given, ref.Abbreviations will be used.
func (r *CanonicalRef) AbbreviatedRef(opt ...ResolveOption) (string, error) {
	o := makeResolveOpts(opt)
	abbrs := o.Abbreviations

	if abbrs == nil {
		return r.compactRef(r.Book.Name)
	}

	abbrName, err := abbrs.PreferredAbbreviation(r.Book.Name)
	if err != nil {
		return "", err
	}

	return r.compactRef(abbrName)
}
