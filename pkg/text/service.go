package text

import (
	"errors"

	"github.com/zostay/today/pkg/bible"
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
		Abbreviations: bible.Abbreviations,
		Canon:         bible.Protestant,
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
