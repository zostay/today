package text

import (
	"context"
	"fmt"
	"html/template"

	"github.com/zostay/today/internal/options"
	"github.com/zostay/today/pkg/ref"
)

type VerseOption = options.VerseOption

var (
	// WithVerseFormats specifies the returned verse should return the verse in the given formats.
	WithFormats = options.WithFormats
	// WithOnlyVerseFormats specifies the returned verse should be limited to the given formats.
	WithOnlyFormats = options.WithOnlyFormats
	// WithHTMLVerse specifies the returned verse should include the verse in HTML format.
	WithHTML = options.WithHTML
	// WithOnlyHTMLVerse specifies the returned verse should be limited to the HTML format.
	WithOnlyHTML = options.WithOnlyHTML
	// WithTextVerse specifies the returned verse should include the verse in text format.
	WithText = options.WithText
	// WithOnlyTextVerse specifies the returned verse should be limited to the text format.
	WithOnlyText = options.WithOnlyText
)

// Verse fetches a complete verse with assocaited metadata for the given verse
// reference.
func (s *Service) Verse(ctx context.Context, vr string, vsopt ...VerseOption) (*Verse, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return nil, fmt.Errorf("error parsing the verse reference %q: %w", vr, err)
	}

	return s.verse(ctx, res, vsopt)
}

func (s *Service) verse(ctx context.Context, res *ref.Resolved, vsopt []VerseOption) (*Verse, error) {
	opts := options.MakeVerseOptions(vsopt)

	l, err := s.Resolver.VerseURI(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("error getting verse URI for %q: %w", res.Ref(), err)
	}

	vi, err := s.VersionInformation(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting verse information for %q: %w", res.Ref(), err)
	}

	c := make(Content, opts.IncludeFormats.Len())
	for _, f := range opts.IncludeFormats.Keys() {
		c[f], err = s.Resolver.VerseAs(ctx, res, f)
		if err != nil {
			return nil, fmt.Errorf("error getting verse as %q: %w", f, err)
		}
	}

	var resOpts []ref.ResolveOption
	if s.Abbreviations != nil {
		resOpts = append(resOpts, ref.WithAbbreviations(s.Abbreviations))
	}

	cref, err := res.CompactRef(resOpts...)
	if err != nil {
		return nil, fmt.Errorf("error compacting the verse reference %q: %w", res.Ref(), err)
	}

	return &Verse{
		Reference: cref,
		Content:   c,
		Link:      l,
		Version:   *vi,
	}, nil
}

func (s *Service) VerseText(ctx context.Context, vr string) (string, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return s.Resolver.VerseAs(ctx, res, "text")
}

func (s *Service) VerseHTML(ctx context.Context, vr string) (template.HTML, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return "", err
	}

	vs, err := s.Resolver.VerseAs(ctx, res, "html")
	return template.HTML(vs), err
}
