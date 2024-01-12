package text

import (
	"errors"
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

var (
	ErrMultiVerse = errors.New("multiple verses not supported")
)

type Service struct {
	Resolver
}

func NewService(r Resolver) *Service {
	return &Service{r}
}

func parseToResolved(vr string) (*ref.Resolved, error) {
	pr, err := ref.ParseProper(vr)
	if err != nil {
		return nil, err
	}

	if !pr.IsSingleRange() {
		return nil, ErrMultiVerse
	}

	ref, err := ref.Canonical.Resolve(pr)
	if err != nil {
		return nil, err
	}

	return &ref[0], nil
}

func (c *Service) Verse(vr string) (string, error) {
	res, err := parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return c.Resolver.Verse(res)
}

func (c *Service) VerseHTML(vr string) (template.HTML, error) {
	res, err := parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return c.Resolver.VerseHTML(res)
}

func (c *Service) RandomVerse(opt ...ref.RandomReferenceOption) (*ref.Resolved, string, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := c.Resolver.Verse(res)
	return res, txt, err
}

func (c *Service) RandomVerseHTML(opt ...ref.RandomReferenceOption) (*ref.Resolved, template.HTML, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := c.Resolver.VerseHTML(res)
	return res, txt, err
}
