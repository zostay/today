package main

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
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

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show a specified scripture",
		Args:  cobra.MinimumNArgs(1),
		Run:   RunTodayShow,
	}

	listBooks := &cobra.Command{
		Use:   "books",
		Short: "List the available books",
		Args:  cobra.NoArgs,
		Run:   RunListBooks,
	}

	listCategories := &cobra.Command{
		Use:   "categories",
		Short: "List the available categories",
		Args:  cobra.NoArgs,
		Run:   RunListCategories,
	}

	randomCmd := &cobra.Command{
		Use:   "random",
		Short: "Pick a scripture to read at random",
		Args:  cobra.ExactArgs(0),
		RunE:  RunTodayRandom,
	}

	cmd.AddCommand(randomCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(listCategories)
	cmd.AddCommand(listBooks)

	showCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	randomCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	randomCmd.Flags().StringVarP(&fromCategory, "category", "c", "", "Pick a random verse from a category")
	randomCmd.Flags().StringVarP(&fromBook, "book", "b", "", "Pick a random verse from a book")
}

func RunListCategories(cmd *cobra.Command, args []string) {
	for c := range ref.Categories {
		fmt.Println(c)
	}
}

func RunListBooks(cmd *cobra.Command, args []string) {
	for _, b := range ref.Canonical {
		fmt.Println(b.Name)
	}
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
		vhtml, err = svc.VerseHTML(ref)
		v = string(vhtml)
	} else {
		v, err = svc.Verse(ref)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(wrap.Wrap(v, 70))
}

func main() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
