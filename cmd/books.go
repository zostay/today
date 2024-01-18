package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
)

var listBooksCmd = &cobra.Command{
	Use:   "books",
	Short: "List the available books",
	Args:  cobra.NoArgs,
	Run:   RunListBooks,
}

func RunListBooks(cmd *cobra.Command, args []string) {
	for _, b := range ref.Canonical.Books {
		fmt.Println(b.Name)
	}
}
