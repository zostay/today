package ref_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/ref"
)

func TestAbbrTree(t *testing.T) {
	t.Parallel()

	abbrs := ref.BookAbbreviations{
		Abbreviations: []ref.BookAbbreviation{
			{
				Name:    "Abc",
				Accepts: []string{"Abc", "Ac"},
			},
			{
				Name:    "Def",
				Accepts: []string{"Def", "Df"},
			},
			{
				Name:    "Acd",
				Accepts: []string{"Acd", "Ad"},
			},
		},
	}

	tree := ref.NewAbbrTree(&abbrs)

	res := tree.Get("A")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Abc": &abbrs.Abbreviations[0],
		"Acd": &abbrs.Abbreviations[2],
	}, res)

	res = tree.Get("Ab")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Abc": &abbrs.Abbreviations[0],
	}, res)

	res = tree.Get("Abc")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Abc": &abbrs.Abbreviations[0],
	}, res)

	res = tree.Get("Ac")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Abc": &abbrs.Abbreviations[0],
	}, res)

	res = tree.Get("acd")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Acd": &abbrs.Abbreviations[2],
	}, res)

	res = tree.Get(".ad")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Acd": &abbrs.Abbreviations[2],
	}, res)

	res = tree.Get("B")
	assert.Nil(t, res)

	res = tree.Get("D")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Def": &abbrs.Abbreviations[1],
	}, res)

	res = tree.Get("DE")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Def": &abbrs.Abbreviations[1],
	}, res)

	res = tree.Get("D E F.")
	assert.Equal(t, map[string]*ref.BookAbbreviation{
		"Def": &abbrs.Abbreviations[1],
	}, res)

	res = tree.Get("D E F. G")
	assert.Nil(t, res)

	res = tree.Get("")
	assert.Nil(t, res)
}
