package ref

import "fmt"

// Pericope represents a resolved extract from a canon.
type Pericope struct {
	Ref   *Resolved
	Canon *Canon
	Title string
}

// Lookup returns a Pericope from the given canon and reference.
func Lookup(c *Canon, ref, title string) (*Pericope, error) {
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

// Verses returns a channel which emits each verse in the Pericope.
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
