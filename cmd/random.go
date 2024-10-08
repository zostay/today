package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"

	"github.com/zostay/today/pkg/ost"
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
	excludeIndex                 string
	exclude                      []string
)

func init() {
	randomCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	randomCmd.Flags().StringVarP(&fromCategory, "category", "c", "", "Pick a random verse from a category")
	randomCmd.Flags().StringVarP(&fromBook, "book", "b", "", "Pick a random verse from a book")
	randomCmd.Flags().UintVarP(&minimumVerses, "minimum-verses", "m", 1, "Minimum number of verses to include in the random selection")
	randomCmd.Flags().UintVarP(&maximumVerses, "maximum-verses", "M", 1, "Maximum number of verses to include in the random selection")
	randomCmd.Flags().StringVarP(&excludeIndex, "exclude-index", "X", "", "Exclude all passages references in the specified index file")
	randomCmd.Flags().StringSliceVarP(&exclude, "exclude", "x", []string{}, "Exclude the specified passage references")
}

func loadIndex(path string) (*ost.Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var idx ost.Index
	if filepath.Ext(path) == ".yaml" {
		err = ost.LoadIndexYaml(f, &idx)
	} else {
		err = ost.LoadIndexJson(f, &idx)
	}

	return &idx, err
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

	excludeRefs := make([]string, 0, len(exclude))
	if len(exclude) > 0 {
		excludeRefs = append(excludeRefs, exclude...)
	}

	if excludeIndex != "" {
		idx, err := loadIndex(excludeIndex)
		if err != nil {
			panic(err)
		}

		refs := make([]string, 0, len(idx.Verses))
		for _, v := range idx.Verses {
			refs = append(refs, v.Reference)
		}

		excludeRefs = append(excludeRefs, refs...)
	}

	if len(excludeRefs) > 0 {
		opts = append(opts, ref.ExcludeReferences(excludeRefs...))
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
		v  string
		vr *ref.Resolved
	)

	if asHtml {
		var vh template.HTML
		vr, vh, err = svc.RandomVerseHTML(cmd.Context(), opts...)
		v = string(vh)
	} else {
		vr, v, err = svc.RandomVerseText(cmd.Context(), opts...)
	}
	if err != nil {
		var ucerr *ref.UnknownCategoryError
		if errors.As(err, &ucerr) {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", ucerr)
			return nil
		}
		panic(err)
	}

	if asHtml {
		sref, err := vr.CompactRef()
		if err != nil {
			panic(err)
		}
		v = "<h1>" + sref + "</h1>\n" + v
	} else {
		sref, err := vr.CompactRef()
		if err != nil {
			panic(err)
		}
		v = sref + "\n\n" + v
	}

	fmt.Println(wrap.Wrap(v, 70))

	return nil
}
