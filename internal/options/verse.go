package options

import "github.com/zostay/go-std/set"

type Verse struct {
	IncludeFormats set.Set[string]
}

type VerseOption interface {
	applyV(*Verse)
}

type VerseOpt func(*Verse)

func (o VerseOpt) applyV(v *Verse) {
	o(v)
}

func WithFormats(f ...string) VerseOpt {
	return VerseOpt(func(o *Verse) {
		fset := set.New(f...)
		if o.IncludeFormats == nil {
			o.IncludeFormats = set.Union(o.IncludeFormats, fset)
		}
	})
}

func WithOnlyFormats(f ...string) VerseOpt {
	return VerseOpt(func(o *Verse) {
		o.IncludeFormats = set.New(f...)
	})
}

func WithHTML() VerseOpt {
	return WithFormats("html")
}

func WithOnlyHTML() VerseOpt {
	return WithOnlyFormats("html")
}

func WithText() VerseOpt {
	return WithFormats("text")
}

func WithOnlyText() VerseOpt {
	return WithOnlyFormats("text")
}

func defaultVerse() *Verse {
	return &Verse{
		IncludeFormats: set.New("text", "html"),
	}
}

func MakeVerseOptions(opt []VerseOption) *Verse {
	v := defaultVerse()
	for _, f := range opt {
		f.applyV(v)
	}
	return v
}
