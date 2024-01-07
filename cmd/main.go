package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

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
		fmt.Println(b.Name())
	}
}

func RunTodayRandom(cmd *cobra.Command, args []string) error {
	//keeper.RequiresSecretKeeper()

	if fromCategory != "" && fromBook != "" {
		return errors.New("cannot specify both --category and --book")
	}

	opts := []text.RandomReferenceOption{}
	if fromCategory != "" {
		opts = append(opts, text.FromCategory(fromCategory))
	}
	if fromBook != "" {
		opts = append(opts, text.FromBook(fromBook))
	}

	r, err := esv.NewFromAuthFile(esv.AuthFile)
	if err != nil {
		panic(err)
	}
	esvClient := text.NewService(r)

	rand.Seed(time.Now().UTC().UnixNano())
	var v string
	if asHtml {
		v, err = esvClient.RandomVerseHTML(opts...)
	} else {
		v, err = esvClient.RandomVerse(opts...)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(wrap.Wrap(v, 70))

	return nil
}

func RunTodayShow(cmd *cobra.Command, args []string) {
	keeper.RequiresSecretKeeper()

	ref := strings.Join(args, " ")
	var (
		err error
		v   string
	)
	if asHtml {
		v, err = esv.GetVerseHTML(ref)
	} else {
		v, err = esv.GetVerse(ref)
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
