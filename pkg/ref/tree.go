package ref

import "unicode"

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

	return cur.Final
}
