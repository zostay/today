package text

import "html/template"

type Service struct {
	Resolver
}

func NewService(r Resolver) *Service {
	return &Service{r}
}

func (c *Service) Verse(ref string) (string, error) {
	pr, err := ref.Parse(ref)
	if err != nil {
		return "", err
	}

	return c.Resolver.Verse(pr)
}

func (c *Service) VerseHTML(ref string) (template.HTML, error) {
	pr, err := ref.Parse(ref)
	if err != nil {
		return "", err
	}

	return c.Resolver.VerseHTML(pr)
}
