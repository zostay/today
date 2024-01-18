package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
)

func init() {
	listCategories := &cobra.Command{
		Use:   "categories",
		Short: "List the available categories",
		Args:  cobra.NoArgs,
		Run:   RunListCategories,
	}

	cmd.AddCommand(listCategories)
}

func RunListCategories(cmd *cobra.Command, args []string) {
	for c := range ref.Canonical.Categories {
		fmt.Println(c)
	}
}
