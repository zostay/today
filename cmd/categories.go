package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/zostay/go-std/maps"

	"github.com/zostay/today/pkg/ref"
)

var (
	listCategoriesCmd = &cobra.Command{
		Use:   "categories",
		Short: "List the available categories",
		Args:  cobra.NoArgs,
		Run:   RunListCategories,
	}

	listPericopes bool
)

func init() {
	listCategoriesCmd.Flags().BoolVarP(&listPericopes, "pericopes", "p", false, "List the pericopes in each category")
}

func RunListCategories(cmd *cobra.Command, args []string) {
	cats := maps.Keys(ref.Canonical.Categories)
	sort.Strings(cats)

	for _, c := range cats {
		fmt.Println(c)
		if listPericopes {
			ps, err := ref.Canonical.Category(c)
			if err != nil {
				panic(err)
			}

			sort.Slice(ps, func(i, j int) bool {
				return ps[i].Ref.Ref() < ps[j].Ref.Ref()
			})

			for _, p := range ps {
				fmt.Printf("  %s\n", p.Ref.Ref())
			}
		}
	}
}
