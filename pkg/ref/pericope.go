package ref

// Pericope represents a resolved extract from a canon.
type Pericope struct {
	Ref *Multiple

	Canon Canon
	Book  *Book

	Name   string
	Title  string
	Verses []Verse
}

//func NewPericope(c Canon, ref, title string) (*Pericope, error) {
//	m, err := ParseMultiple(ref)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := m.Validate(); err != nil {
//		return nil, err
//	}
//
//	b, err := c.Book(m[0].Book)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Pericope{
//		Ref: m,
//
//		Canon: c,
//		Book:  b,
//
//		Name:  b.Name,
//		Title: title,
//	}, nil
//}
