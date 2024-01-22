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

var (
	randomCmd = &cobra.Command{
		Use:   "random",
		Short: "Pick a scripture to read at random",
		Args:  cobra.ExactArgs(0),
		RunE:  RunTodayRandom,
	}

	minimumVerses, maximumVerses uint
)

func init() {
	randomCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	randomCmd.Flags().StringVarP(&fromCategory, "category", "c", "", "Pick a random verse from a category")
	randomCmd.Flags().StringVarP(&fromBook, "book", "b", "", "Pick a random verse from a book")
	randomCmd.Flags().UintVarP(&minimumVerses, "minimum-verses", "m", 1, "Minimum number of verses to include in the random selection")
	randomCmd.Flags().UintVarP(&maximumVerses, "maximum-verses", "M", 1, "Maximum number of verses to include in the random selection")
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
	if minimumVerses != 0 {
		opts = append(opts, ref.WithAtLeast(minimumVerses))
	}
	if maximumVerses != 0 {
		opts = append(opts, ref.WithAtMost(maximumVerses))
	}

	if minimumVerses > maximumVerses {
		return errors.New("--minimum-verses cannot be greater than --maximum-verses")
	}

	ec, err := esv.NewFromEnvironment()
	if err != nil {
		panic(err)
	}
	svc := text.NewService(ec)

	var (
		v string
	)
	if asHtml {
		var vh template.HTML
		_, vh, err = svc.RandomVerseHTML(opts...)
		v = string(vh)
	} else {
		_, v, err = svc.RandomVerse(opts...)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(wrap.Wrap(v, 70))

	return nil
}
