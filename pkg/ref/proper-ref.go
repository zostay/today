package ref

import (
	"fmt"
	"strings"

	"github.com/zostay/go-std/slices"
)

type FullRef []Range

func (f FullRef) String() string {
	return strings.Join(
		slices.Map(f, func(r Range) string {
			return r.String()
		}),
		"; ")
}

type ProperRef struct {
	*Book
	FullRef
}

func (p *ProperRef) String() string {
	return fmt.Sprintf("%s %s", p.Book.Name, p.FullRef.String())
}

type MasterRef struct {
	refs map[string]*ProperRef
}

func NewMasterRef() *MasterRef {
	m := &MasterRef{
		refs: map[string]*ProperRef{},
	}

	for i := range Canonical {
		m.refs[Canonical[i].Name] = &ProperRef{
			Book:    &Canonical[i],
			FullRef: FullRef{},
		}
	}

	return m
}

func (m *MasterRef) Add(book string, rref Range) {
	m.refs[book].FullRef = append(m.refs[book].FullRef, rref)
}
