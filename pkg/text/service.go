package text

import (
	"context"
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

func (c *Service) VersionInformation(ctx context.Context) (*Version, error) {
	return c.Resolver.VersionInformation(ctx)
}

func (c *Service) Verse(ctx context.Context, vr string) (*Verse, error) {
	res, err := parseToResolved(vr)
	if err != nil {
		return nil, err
	}

	return c.Resolver.Verse(ctx, res)
}

func (c *Service) VerseText(ctx context.Context, vr string) (string, error) {
	res, err := parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return c.Resolver.VerseText(ctx, res)
}

func (c *Service) VerseHTML(ctx context.Context, vr string) (template.HTML, error) {
	res, err := parseToResolved(vr)
	if err != nil {
		return "", err
	}

	return c.Resolver.VerseHTML(ctx, res)
}

func (c *Service) RandomVerse(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, *Verse, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Resolver.Verse(ctx, res)
	return res, v, err
}

func (c *Service) RandomVerseText(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, string, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := c.Resolver.VerseText(ctx, res)
	return res, txt, err
}

func (c *Service) RandomVerseHTML(ctx context.Context, opt ...ref.RandomReferenceOption) (*ref.Resolved, template.HTML, error) {
	res, err := ref.Random(opt...)
	if err != nil {
		return nil, "", err
	}

	txt, err := c.Resolver.VerseHTML(ctx, res)
	return res, txt, err
}
