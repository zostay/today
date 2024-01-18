package cmd

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Pick a scripture to read at random",
	Args:  cobra.ExactArgs(0),
	RunE:  RunTodayRandom,
}

func init() {
	randomCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	randomCmd.Flags().StringVarP(&fromCategory, "category", "c", "", "Pick a random verse from a category")
	randomCmd.Flags().StringVarP(&fromBook, "book", "b", "", "Pick a random verse from a book")
}

func RunTodayRandom(cmd *cobra.Command, args []string) error {
	if fromCategory != "" && fromBook != "" {
		return errors.New("cannot specify both --category and --book")
	}

	var opts []ref.RandomReferenceOption
	if fromCategory != "" {
		opts = append(opts, ref.FromCategory(fromCategory))
	}
	if fromBook != "" {
		opts = append(opts, ref.FromBook(fromBook))
	}

	ec, err := esv.NewFromEnvironment()
	if err != nil {
		panic(err)
	}
	esvClient := text.NewService(ec)

	var (
		r ref.Ref
		v string
	)
	if asHtml {
		var vh template.HTML
		r, vh, err = esvClient.RandomVerseHTML(opts...)
		v = string(vh)
	} else {
		r, v, err = esvClient.RandomVerse(opts...)
	}
	if err != nil {
		panic(err)
	}
	v += "\n\n" + r.Ref()
	fmt.Println(wrap.Wrap(v, 70))

	return nil
}
