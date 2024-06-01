package ref

import (
	"fmt"

	"github.com/zostay/today/pkg/canon"
)

// Pericope represents a resolved extract from a canon.
type Pericope struct {
	Ref   *canon.Resolved
	Canon *canon.Canon
	Title string
}

// Category returns a list of Pericopes associated with that Category or nil if
// no such category is defined. Returns nil and error if there's a problem with
// the category definition.
func PericopeFromCanonCategory(c *canon.Canon, name string) ([]*Pericope, error) {
	if refs, hasCategory := c.Categories[name]; hasCategory {
		var ps []*Pericope
		for i := range refs {
			p, err := Lookup(c, refs[i], "")
			if err != nil {
				return nil, err
			}
			ps = append(ps, p)
		}
		return ps, nil
	}
	return nil, nil
}

func Lookup(c *canon.Canon, ref, title string) (*Pericope, error) {
	p, err := ParseProper(ref)
	if err != nil {
		return nil, err
	}

	if !p.IsSingleRange() {
		return nil, fmt.Errorf("pericope must be constructed from a single range: %s", ref)
	}

	res, err := c.Resolve(p)
	if err != nil {
		return nil, err
	}

	return &Pericope{
		Ref:   &res[0],
		Canon: c,
		Title: title,
	}, nil
}

func (p *Pericope) Verses() (<-chan Verse, error) {
	ch := make(chan Verse)

	go func() {
		defer close(ch)

		for _, v := range p.Ref.Verses() {
			ch <- v
		}
	}()

	return ch, nil
}
