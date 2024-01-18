package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
)

func init() {
	listBooks := &cobra.Command{
		Use:   "books",
		Short: "List the available books",
		Args:  cobra.NoArgs,
		Run:   RunListBooks,
	}

	cmd.AddCommand(listBooks)
}

func RunListBooks(cmd *cobra.Command, args []string) {
	for _, b := range ref.Canonical.Books {
		fmt.Println(b.Name)
	}
}
