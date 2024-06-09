package ref

import (
	"fmt"
	"strings"
)

// resolveAll takes a list of references and resolves them.
func (c *Canon) resolveAll(exclude []string) ([]Resolved, error) {
	rs := make([]Resolved, 0, len(exclude))
	for _, e := range exclude {
		p, err := ParseMultiple(e)
		if err != nil {
			return nil, err
		}

		prs, err := c.Resolve(p, WithAbbreviations(Abbreviations))
		if err != nil {
			return nil, err
		}

		rs = append(rs, prs...)
	}

	return rs, nil
}

// vCmp compares two verses. Returns -1 if v1 < v2, 1 if v1 > v2, and 0 if
// v1 == v2.
func vCmp(v1, v2 Verse) int {
	switch {
	case v1.Before(v2):
		return -1
	case v2.Before(v1):
		return 1
	default:
		return 0
	}
}

// mergeReferences takes a listed of Resolved references and looks for any
// overlaps. It merges overlaps together and returns a new list of Resolved.
func mergeReferences(rs []Resolved) []Resolved {
	merged := make([]Resolved, 0, len(rs))
	for _, r := range rs {
		performedMerge := false
	INNER:
		for _, m := range merged {
			if m.Book.Name != r.Book.Name {
				continue
			}

			switch {
			case vCmp(m.First, r.First) < 0 && vCmp(m.Last, r.First) > 0:
				if vCmp(m.Last, r.Last) > 0 {
					m.Last = r.Last
				}
				performedMerge = true
				break INNER
			case vCmp(m.First, r.Last) < 0 && vCmp(m.Last, r.Last) > 0:
				if vCmp(m.First, r.First) < 0 {
					m.First = r.First
				}
				performedMerge = true
				break INNER
			}
		}

		if !performedMerge {
			merged = append(merged, r)
		}
	}
	return merged
}

// Filtered returns a new canon with the excluded references removed.
func (c *Canon) Filtered(exclude ...string) (*Canon, error) {
	copyCanon := c.Clone()
	copyCanon.Name += " (excluding " + strings.Join(exclude, ", ") + ")"

	// convert excluded verse ref strings to Resolved
	rs, err := c.resolveAll(exclude)
	if err != nil {
		return nil, err
	}

	// merge any overlapping ranges
	rs = mergeReferences(rs)

	err = copyCanon.filterOutCategories(rs)
	if err != nil {
		return nil, err
	}

	err = copyCanon.filterOutVerses(rs)
	if err != nil {
		return nil, err
	}

	copyCanon.filterOutBooks()

	return copyCanon, nil
}

// filterOutVerses removes verses from the canon that are listed in the
// exclusion list.
func (c *Canon) filterOutVerses(rs []Resolved) error {
	for _, r := range rs {
		bi := -1
		for i := range c.Books {
			if c.Books[i].Name == r.Book.Name {
				bi = i
			}
		}

		if bi < 0 {
			// should never happen, RIGHT?!?
			return fmt.Errorf("%w: unable to find book during filtering", ErrNotFound)
		}

		// Find the edges of each range inside the canon
		first, last := -1, -1
		for i := range c.Books[bi].Verses {
			if c.Books[bi].Verses[i].Equal(r.First) {
				first = i
			}
			if c.Books[bi].Verses[i].Equal(r.Last) {
				last = i
				break
			}
		}

		// clip out the excluded range
		orig := c.Books[bi].Verses
		c.Books[bi].Verses = orig[:first]
		c.Books[bi].Verses = append(c.Books[bi].Verses, orig[last+1:]...)
	}

	return nil
}

// filterOutBooks removes any books with no verses.
func (c *Canon) filterOutBooks() {
	for i := len(c.Books) - 1; i >= 0; i-- {
		if len(c.Books[i].Verses) == 0 {
			c.Books = append(c.Books[:i], c.Books[i+1:]...)
		}
	}
}

// filterOutCategories rewrites the categories to only include books that are
// found in the canon and any more specific ranges are pruned or removed based
// upon which passages remain in a canon after filtering.
func (c *Canon) filterOutCategories(rs []Resolved) error {
	for k, v := range c.Categories {
		newV := make([]string, 0, len(v))
		for _, sr := range v {
			pr, err := ParseProper(sr)
			if err != nil {
				return err
			}

			thisR, err := c.resolveProper(pr, &resolveOpts{})
			if err != nil {
				return err
			}

			in := make([]Resolved, len(thisR), len(rs))
			copy(in, thisR)
			newIn := make([]Resolved, 0, len(rs))
			for _, exR := range rs {
				for _, inR := range in {
					newIn = append(newIn, inR.Subtract(&exR)...)
				}
				in, newIn = newIn, in[:0]

				if len(in) == 0 {
					break
				}
			}

			for _, inR := range in {
				s, err := inR.CompactRef()
				if err != nil {
					return err
				}
				newV = append(newV, s)
			}
		}

		c.Categories[k] = newV
	}

	return nil
}
