package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zostay/today/pkg/ost"
	"gopkg.in/yaml.v3"
)

var (
	ostIndexCmd = &cobra.Command{
		Use:   "index",
		Short: "Fetch an index of scriptures posted to openscripture.today",
		Args:  cobra.NoArgs,
		Run:   RunOstIndex,
	}

	indexForMonth, indexForYear string
	asList                      bool
)

func init() {
	ostIndexCmd.Flags().StringVarP(&indexForMonth, "month", "m", "", "Fetch the index for a specific month (YYYY/MM)")
	ostIndexCmd.Flags().StringVarP(&indexForYear, "year", "y", "", "Fetch the index for a specific year (YYYY)")
	ostIndexCmd.Flags().BoolVarP(&asList, "list", "l", false, "Output the index as a list")
}

func RunOstIndex(cmd *cobra.Command, args []string) {
	client, err := ost.New(cmd.Context())
	if err != nil {
		panic(err)
	}

	if indexForMonth != "" && indexForYear != "" {
		panic("Cannot specify both --month and --year")
	}

	opts := []ost.IndexOption{}
	switch {
	case indexForMonth != "":
		parts := strings.SplitN(indexForMonth, "/", 2)
		year, month := parts[0], parts[1]
		if len(year) != 4 || len(month) != 2 {
			panic("Invalid month format")
		}
		opts = append(opts, ost.ForMonth(year, month))
	case indexForYear != "":
		if len(indexForYear) != 4 {
			panic("Invalid year format")
		}
		opts = append(opts, ost.ForYear(indexForYear))
	default:
		opts = append(opts, ost.ForAllTime())
	}

	idx, err := client.VerseIndex(cmd.Context(), opts...)
	if err != nil {
		panic(err)
	}

	if asList {
		for _, v := range idx.Verses {
			fmt.Printf("%s\n", v.Reference)
		}
		return
	}

	enc := yaml.NewEncoder(cmd.OutOrStdout())
	err = enc.Encode(idx)
	if err != nil {
		panic(err)
	}
}
