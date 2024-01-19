package cmd

import (
	"fmt"
	"html/template"
	"time"

	"github.com/bbrks/wrap"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/ost"
)

var (
	ostCmd = &cobra.Command{
		Use:     "openscripture",
		Short:   "Work with current scripture of the day from openscripture.today",
		Args:    cobra.NoArgs,
		Run:     RunOst,
		Aliases: []string{"ost"},
	}

	ostTodayCmd = &cobra.Command{
		Use:   "today",
		Short: "Show the current scripture of the day from openscripture.today",
		Args:  cobra.NoArgs,
		Run:   RunOst,
	}

	ostOnCmd = &cobra.Command{
		Use:   "on",
		Short: "Show the scripture of the day for a specific date from openscripture.today",
		Args:  cobra.ExactArgs(1),
		Run:   RunOst,
	}

	asMeta, asYaml bool
)

func init() {
	ostCmd.AddCommand(ostTodayCmd, ostOnCmd)

	ostCmd.PersistentFlags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostCmd.PersistentFlags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostCmd.PersistentFlags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
}

func RunOst(cmd *cobra.Command, args []string) {
	opts := []ost.Option{}
	if len(args) == 1 {
		on := args[0]
		onTime, err := time.ParseInLocation("2006-01-02", on, time.Local)
		if err != nil {
			panic(err)
		}

		opts = append(opts, ost.On(onTime))
	}

	client, err := ost.New()
	if err != nil {
		panic(err)
	}

	var v string
	switch {
	case asYaml:
		vv, err := client.TodayVerse(opts...)
		if err != nil {
			panic(err)
		}

		enc := yaml.NewEncoder(cmd.OutOrStdout())
		err = enc.Encode(vv)
		if err != nil {
			panic(err)
		}
		return
	case asMeta:
		vv, err := client.TodayVerse(opts...)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Reference: %s\n", vv.Reference)
		fmt.Printf("Version:   %s\n", vv.Version.Name)
		fmt.Printf("Link:      %s\n", vv.Version.Link)
		return
	case asHtml:
		var vh template.HTML
		vh, err = client.TodayHTML(opts...)
		v = string(vh)
	default:
		v, err = client.Today(opts...)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(wrap.Wrap(v, 70))
}
