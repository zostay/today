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

// vSucceeds compares two verses. Returns true if v2 immediately follows v1.
func vSucceeds(b *Book, v1, v2 Verse) bool {
	for i, v := range b.Verses {
		if v.Equal(v1) && i+1 < len(b.Verses) {
			return b.Verses[i+1].Equal(v2)
		}
	}
	return false
}

// mergeReferences takes a listed of Resolved references and looks for any
// overlaps. It merges overlaps together and returns a new list of Resolved.
func mergeReferences(rs []Resolved) []Resolved {
	// TODO The recursion here to get it down to a minimum set is a little on
	//  the insane side. Should probably think of a more performant way.
	for {
		merged := make([]Resolved, 0, len(rs))

		for _, b := range rs {
			performedMerge := false
		INNER:
			for ai, a := range merged {
				if a.Book.Name != b.Book.Name {
					continue
				}

				// Overlaps come in four flavors, which can all be detected with three tests:
				// #1 A------B====A'-----B' - first range starts before second and second ends after first
				// #2 B------A====B'-----A' - second range starts before first and first ends after second
				// #3 A---B=====B'----A' - first fully contains second
				// #4 B---A=====A'----B' - second fully contains first
				//
				// A) #1 and #3 can be detected by testing if A <= B <= A'
				// B) #2 and #3 can be detected by testing if A <= B' <= A'
				// C) #4 can be detected by testing if B <= A <= B'
				// D) The verses do not overlap, but B is succeeded by A'.
				// E) The verses do not overlap, but B' immediately follows A.
				//
				// If we perform the tests in that order, then we can assume B is only detecting #2 and make
				// assumptions accordingly.

				switch {
				// Test (A): Detect #1 and #3 above
				case vCmp(a.First, b.First) <= 0 && vCmp(a.Last, b.First) >= 0:
					// If this is true, then this is #1, not #3
					if vCmp(a.Last, b.Last) < 0 {
						merged[ai].Last = b.Last
					}
					performedMerge = true
					break INNER
				// Test (B): Detect #2
				case vCmp(a.First, b.Last) <= 0 && vCmp(a.Last, b.Last) >= 0:
					merged[ai].First = b.First
					performedMerge = true
					break INNER
				// Test (C): Detect #4
				case vCmp(a.First, b.First) >= 0 && vCmp(a.First, b.Last) <= 0:
					merged[ai].First = b.First
					merged[ai].Last = b.Last
					performedMerge = true
				// Test (D): Detect if B is immediately follwed by A'
				case vSucceeds(a.Book, a.Last, b.First):
					merged[ai].First = a.First
					merged[ai].Last = b.Last
					performedMerge = true
				// Test (E): Detect if B' immediately followed by A
				case vSucceeds(a.Book, b.Last, a.First):
					merged[ai].First = b.First
					merged[ai].Last = a.Last
					performedMerge = true
				}
			}

			if !performedMerge {
				merged = append(merged, b)
			}
		}

		oldLen := len(rs)
		rs = merged
		if len(merged) < oldLen {
			continue
		}

		break
	}
	return rs
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
			cs, _ := r.CompactRef()
			return fmt.Errorf("%w: unable to find book while excluding %q", ErrNotFound, cs)
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

		if first < 0 || last < 0 {
			// should never happen, RIGHT?!?
			cs, _ := r.CompactRef()
			which := "both verses"
			if last >= 0 {
				which = "first verse"
			} else if first >= 0 {
				which = "last verse"
			}
			return fmt.Errorf("%w: unable to find %s while excluding %q", ErrNotFound, which, cs)
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
