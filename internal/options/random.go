package options

import (
	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/canon"
)

type RandomReference struct {
	Category string
	Book     string
	Canon    *canon.Canon
	Min, Max int
}

type RandomReferenceOption interface {
	applyRR(*RandomReference)
}

type RandomReferenceOpt func(*RandomReference)

func (o RandomReferenceOpt) applyRR(rr *RandomReference) {
	o(rr)
}

func FromCanon(canon *canon.Canon) RandomReferenceOpt {
	return func(o *RandomReference) {
		o.Canon = canon
	}
}

func FromBook(name string) RandomReferenceOpt {
	return func(o *RandomReference) {
		o.Book = name
	}
}

func FromCategory(name string) RandomReferenceOpt {
	return func(o *RandomReference) {
		o.Category = name
	}
}

func WithAtLeast(n uint) RandomReferenceOpt {
	return func(o *RandomReference) {
		o.Min = int(n)
	}
}

func WithAtMost(n uint) RandomReferenceOpt {
	return func(o *RandomReference) {
		o.Max = int(n)
	}
}

func defaultRandomReference() *RandomReference {
	return &RandomReference{
		Canon: bible.Protestant,
		Min:   1,
		Max:   30,
	}
}

func MakeRandomReferenceOpts(opts []RandomReferenceOption) *RandomReference {
	rr := defaultRandomReference()
	for _, f := range opts {
		f.applyRR(rr)
	}
	return rr
}
