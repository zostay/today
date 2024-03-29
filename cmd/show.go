package cmd

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specified scripture",
	Args:  cobra.MinimumNArgs(1),
	Run:   RunTodayShow,
}

func init() {
	showCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
}

func RunTodayShow(cmd *cobra.Command, args []string) {
	ec, err := esv.NewFromEnvironment()
	if err != nil {
		panic(err)
	}
	svc := text.NewService(ec)

	ref := strings.Join(args, " ")
	var v string
	if asHtml {
		var vhtml template.HTML
		vhtml, err = svc.VerseHTML(cmd.Context(), ref)
		v = string(vhtml)
	} else {
		v, err = svc.VerseText(cmd.Context(), ref)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(wrap.Wrap(v, 70))
}
