package bible

type Resolve struct {
	Abbreviations *BookAbbreviations
	Singular      bool
}

func (r *Resolve) apply(o *Resolve) {
	*r = *o
}

type ResolveOption interface {
	apply(*Resolve)
}

type ResolveOpt func(*Resolve)

func (r ResolveOpt) apply(o *Resolve) {
	r(o)
}

// WithAbbrevations will allow *Ref methods to use the given BookAbbreviations to
// resolve book names.
func WithAbbreviations(abbrs *BookAbbreviations) ResolveOpt {
	return func(o *Resolve) {
		o.Abbreviations = abbrs
	}
}

// WithoutAbbrevations will allow *Ref methods to use no abbreviations object
// during name/abbreviation resolution. (Othrwise, these methods will default to
// ref.Abbreviations.)
func WithoutAbbreviations() ResolveOpt {
	return func(o *Resolve) {
		o.Abbreviations = nil
	}
}

// AsSingleChapter will prefer the singular form of the book name when resolving
// references. (This is special casing for Psalms.)
func AsSingleChapter() ResolveOpt {
	return func(o *Resolve) {
		o.Singular = true
	}
}
