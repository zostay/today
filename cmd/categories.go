package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
)

var listCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List the available categories",
	Args:  cobra.NoArgs,
	Run:   RunListCategories,
}

func RunListCategories(cmd *cobra.Command, args []string) {
	for c := range ref.Canonical.Categories {
		fmt.Println(c)
	}
}
