package cmd

import (
	"fmt"
	"html/template"

	"github.com/bbrks/wrap"
	"github.com/markusmobius/go-dateparser"
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
	ostCmd.AddCommand(ostIndexCmd, ostTodayCmd, ostOnCmd, ostPhotoCmd)

	ostCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")

	ostTodayCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostTodayCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostTodayCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")

	ostOnCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostOnCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostOnCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
}

func RunOst(cmd *cobra.Command, args []string) {
	opts := []ost.DayOption{}
	if len(args) == 1 {
		on := args[0]
		onTime, err := dateparser.Parse(nil, on)
		if err != nil {
			panic(err)
		}

		opts = append(opts, ost.On(onTime.Time))
	}

	client, err := ost.New(cmd.Context())
	if err != nil {
		panic(err)
	}

	var v string
	switch {
	case asYaml:
		vv, err := client.TodayVerse(cmd.Context(), opts...)
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
		vv, err := client.TodayVerse(cmd.Context(), opts...)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Reference: %s\n", vv.Reference)
		fmt.Printf("Version:   %s\n", vv.Version.Name)
		fmt.Printf("Link:      %s\n", vv.Version.Link)
		return
	case asHtml:
		var vh template.HTML
		vh, err = client.TodayHTML(cmd.Context(), opts...)
		v = string(vh)
	default:
		v, err = client.Today(cmd.Context(), opts...)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(wrap.Wrap(v, 70))
}
