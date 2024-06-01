package text

import (
	"context"
	"html/template"

	"github.com/zostay/today/internal/options"
	"github.com/zostay/today/pkg/ref"
)

type VerseRandomReferenceOption = options.VerseRandomReferenceOption
type RandomReferenceOption = options.RandomReferenceOption

var (
	// FromCanon limits the random selection to a specific canon of the Bible.
	FromCanon = options.FromCanon
	// FromBook limits the random selection to a specific book of the Bible.
	FromBook = options.FromBook
	// FromCategory limits the random selection to a specific category of books.
	FromCategory = options.FromCategory
	// WithAtLeast ensures the random selection has at least the given number of verses.
	WithAtLeast = options.WithAtLeast
	// WithAtMost ensures the random selection has at most the given number of verses.
	WithAtMost = options.WithAtMost
)

func xxx() {
	s := &Service{}
	s.RandomVerse(context.Background(), FromBook("John"), WithHTML())
}

func (s *Service) RandomVerse(ctx context.Context, opt ...VerseRandomReferenceOption) (*ref.Resolved, *Verse, error) {
	ro := options.MakeVerseRandomReferenceOptions(opt)
	res, err := ref.Random(ro.RandomReference)
	if err != nil {
		return nil, nil, err
	}

	v, err := s.verse(ctx, res, []VerseOption{ro.Verse})
	return res, v, err
}

func (s *Service) RandomVerseText(ctx context.Context, opt ...RandomReferenceOption) (*ref.Resolved, string, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := s.Resolver.VerseAs(ctx, res, "text")
	return res, txt, err
}

func (s *Service) RandomVerseHTML(ctx context.Context, opt ...RandomReferenceOption) (*ref.Resolved, template.HTML, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := s.Resolver.VerseAs(ctx, res, "html")
	return res, template.HTML(txt), err
}
