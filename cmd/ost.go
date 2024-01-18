package cmd

import (
	"fmt"
	"html/template"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ost"
)

var ostCmd = &cobra.Command{
	Use:     "openscripture.today",
	Short:   "Work with current scripture of the day from openscripture.today",
	Args:    cobra.NoArgs,
	Run:     RunOst,
	Aliases: []string{"ost"},
}

func init() {
	ostCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
}

func RunOst(cmd *cobra.Command, args []string) {
	client, err := ost.New()
	if err != nil {
		panic(err)
	}

	var v string
	if asHtml {
		var vh template.HTML
		vh, err = client.TodayHTML()
		v = string(vh)
	} else {
		v, err = client.Today()
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(wrap.Wrap(v, 70))
}
