package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command

	asHtml bool

	fromCategory string
	fromBook     string
)

func init() {
	cmd = &cobra.Command{
		Use:   "today",
		Short: "Read some scripture today",
	}

	cmd.AddCommand(
		listBooksCmd,
		listCategoriesCmd,
		ostCmd,
		randomCmd,
		showCmd,
		versionCmd,
	)
}

func Execute() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
