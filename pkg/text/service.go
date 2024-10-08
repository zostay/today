package text

import (
	"context"
	"errors"
	"fmt"
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

var (
	ErrMultiVerse = errors.New("multiple verses not supported")
)

type Service struct {
	Resolver
	Abbreviations *ref.BookAbbreviations
	Canon         *ref.Canon
}

type ServiceOption func(*Service)

func WithAbbreviations(abbr *ref.BookAbbreviations) ServiceOption {
	return func(s *Service) {
		s.Abbreviations = abbr
	}
}

func WithoutAbbreviations() ServiceOption {
	return func(s *Service) {
		s.Abbreviations = nil
	}
}

func WithCanon(c *ref.Canon) ServiceOption {
	return func(s *Service) {
		s.Canon = c
	}
}

func NewService(r Resolver, opt ...ServiceOption) *Service {
	s := &Service{
		Resolver:      r,
		Abbreviations: ref.Abbreviations,
		Canon:         ref.Canonical,
	}

	for _, o := range opt {
		o(s)
	}

	return s
}

func (s *Service) parseToResolved(vr string) (*ref.Resolved, error) {
	pr, err := ref.ParseProper(vr)
	if err != nil {
		return nil, err
	}

	if !pr.IsSingleRange() {
		return nil, ErrMultiVerse
	}

	var opts []ref.ResolveOption
	if s.Abbreviations != nil {
		opts = append(opts, ref.WithAbbreviations(s.Abbreviations))
	}

	ref, err := s.Canon.Resolve(pr, opts...)
	if err != nil {
		return nil, err
	}

	return &ref[0], nil
}

func (s *Service) VersionInformation(ctx context.Context) (*Version, error) {
	return s.Resolver.VersionInformation(ctx)
}

func (s *Service) Verse(ctx context.Context, vr string) (*Verse, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return nil, err
	}

	return s.Resolver.Verse(ctx, res)
}

func (s *Service) VerseText(ctx context.Context, vr string) (string, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return s.Resolver.VerseText(ctx, res)
}

func (s *Service) VerseHTML(ctx context.Context, vr string) (template.HTML, error) {
	res, err := s.parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return s.Resolver.VerseHTML(ctx, res)
}

func (s *Service) RandomVerse(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, *Verse, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, nil, err
	}

	v, err := s.Resolver.Verse(ctx, res)
	return res, v, err
}

func (s *Service) RandomVerseText(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, string, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", fmt.Errorf("unable to select random verse: %w", err)
	}

	txt, err := s.Resolver.VerseText(ctx, res)
	if err != nil {
		return res, txt, fmt.Errorf("unable to resolve text for verse %q: %w", res.Ref(), err)
	}

	return res, txt, nil
}

func (s *Service) RandomVerseHTML(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, template.HTML, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", fmt.Errorf("unable to select random verse: %w", err)
	}

	txt, err := s.Resolver.VerseHTML(ctx, res)
	if err != nil {
		return res, txt, fmt.Errorf("unable to resolve HTML for verse %q: %w", res.Ref(), err)
	}

	return res, txt, nil
}
