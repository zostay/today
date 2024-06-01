package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zostay/today/cmd/flag"
	"github.com/zostay/today/cmd/output"
)

var (
	cmd *cobra.Command

	asMeta, asYaml bool
	asHtml         bool
	outputFormat   flag.OutputFormat

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

func collectOutputFormat() output.Format {
	if asMeta && asYaml || asMeta && asHtml || asYaml && asHtml {
		panic("Only one of --meta, --yaml, or --html can be specified")
	}

	if (asMeta || asYaml || asHtml) && outputFormat.IsSet() {
		panic("You cannot use --output with --meta, --yaml, or --html")
	}

	if asMeta {
		outputFormat.Set("meta")
	} else if asYaml {
		outputFormat.Set("yaml")
	} else if asHtml {
		outputFormat.Set("html")
	}

	if outputFormat.IsSet() {
		return outputFormat.Value
	}

	return output.DefaultFormat()
}
