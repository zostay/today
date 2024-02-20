package ref

import (
	"strings"
	"unicode"
)

type AbbrTree struct {
	Children map[rune]*AbbrTree
	Final    map[string]*BookAbbreviation
}

func isSkipped(c rune) bool {
	return unicode.IsSpace(c) || unicode.IsPunct(c)
}

func NewAbbrTree(abbrs *BookAbbreviations) *AbbrTree {
	root := AbbrTree{}
	for i, abbr := range abbrs.Abbreviations {
		for _, acc := range abbr.Accepts {
			cur := &root
			for _, c := range acc {
				if isSkipped(c) {
					continue
				}
				if cur.Children == nil {
					cur.Children = map[rune]*AbbrTree{}
				}
				lc := unicode.ToLower(c)
				if cur.Children[lc] == nil {
					cur.Children[lc] = &AbbrTree{}
				}
				cur = cur.Children[lc]

				if cur.Final == nil {
					cur.Final = map[string]*BookAbbreviation{}
				}
				cur.Final[abbr.Name] = &abbrs.Abbreviations[i]
			}
		}
	}
	return &root
}

func cleanAbbreviation(abbr string) string {
	var cleaned string
	for _, c := range abbr {
		if !isSkipped(c) {
			cleaned += string(c)
		}
	}
	return strings.ToLower(cleaned)
}

func (t *AbbrTree) Get(abbr string) map[string]*BookAbbreviation {
	cur := t
	for _, c := range abbr {
		if isSkipped(c) {
			continue
		}

		cur = cur.Children[unicode.ToLower(c)]

		if cur == nil {
			return nil
		}
	}

	// If there are multiple answers, we will check to see if any of the answers
	// a complete name. If they are, then all the ones tha tare complete names
	// will be returned. If none are complete names, then all are returned. For
	// example, if the input is "Jn" and the accepts list for "John" includes
	// "Jn" and the accepts list for Jonah includes "Jnh", only "John" will be
	// returned even though the initial match includes both.
	if len(cur.Final) > 1 {
		cleanAbbr := cleanAbbreviation(abbr)
		completeNames := map[string]*BookAbbreviation{}
		allOrdinals, someOrdinals := true, false
		for name, finalAbbr := range cur.Final {
			for _, acc := range finalAbbr.Accepts {
				allOrdinals = allOrdinals && finalAbbr.Ordinal != 0
				someOrdinals = someOrdinals || finalAbbr.Ordinal != 0
				cleanAcc := cleanAbbreviation(acc)
				if cleanAcc == cleanAbbr {
					completeNames[name] = finalAbbr
					break
				}
			}
		}

		if len(completeNames) > 0 {
			return completeNames
		}

		// And then, there's this one weird trick we need to distinguish Isaiah
		// from I Samuel...
		if someOrdinals && !allOrdinals {
			for name, finalAbbr := range cur.Final {
				if finalAbbr.Ordinal == 0 {
					completeNames[name] = finalAbbr
				}
			}

			return completeNames
		}
	}

	return cur.Final
}
