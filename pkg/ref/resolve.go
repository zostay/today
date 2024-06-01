package ref

import (
	"fmt"

	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/canon"
)

func makeResolveOpts(opts []bible.ResolveOption) *bible.Resolve {
	r := &bible.Resolve{
		Abbreviations: canon.Abbreviations,
	}
	for _, o := range opts {
		o.apply(r)
	}
	return r
}

// Book will return the Book with the exact given name.
func Book(c *bible.Canon, in string, opt ...bible.ResolveOption) (*bible.Book, error) {
	name := in

	opts := makeResolveOpts(opt)
	if opts.Abbreviations != nil {
		var err error
		name, err = opts.Abbreviations.BookName(in)
		if err != nil {
			return nil, err
		}
	}

	for i := range c.Books {
		b := &c.Books[i]
		if b.Name == name {
			return b, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", bible.ErrNotFound, name)
}
