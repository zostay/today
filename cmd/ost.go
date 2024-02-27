package cmd

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/bbrks/wrap"
	"github.com/markusmobius/go-dateparser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zostay/today/cmd/flag"
	"github.com/zostay/today/pkg/ost"
	"github.com/zostay/today/pkg/photo"
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

	ostPhotoCmd = &cobra.Command{
		Use:   "photo",
		Short: "Show information about or download the photo used with the scripture of the day from openscripture.today",
		Args:  cobra.NoArgs,
		Run:   RunOstPhoto,
	}

	asMeta, asYaml bool
	on             flag.Date
	download       string
)

func init() {
	ostCmd.AddCommand(ostTodayCmd, ostOnCmd, ostPhotoCmd)

	ostCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")

	ostTodayCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostTodayCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostTodayCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")

	ostOnCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostOnCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostOnCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")

	ostPhotoCmd.Flags().StringVarP(&download, "download", "d", "openscripture.jpg", "Download the file photo to the named file")
	ostPhotoCmd.Flags().VarP(&on, "on", "o", "Specify the date to get the photo for")
	ostPhotoCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
}

func RunOst(cmd *cobra.Command, args []string) {
	opts := []ost.Option{}
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

func RunOstPhoto(cmd *cobra.Command, args []string) {
	opts := []ost.Option{}
	if !on.Value.IsZero() {
		opts = append(opts, ost.On(on.Value.Time))
	}

	client, err := ost.New(cmd.Context())
	if err != nil {
		panic(err)
	}

	pi, err := client.TodayPhoto(cmd.Context(), opts...)
	if err != nil {
		panic(err)
	}

	switch {
	case asYaml:
		enc := yaml.NewEncoder(cmd.OutOrStdout())
		err = enc.Encode(pi)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("Link: %s\n", pi.Link)
		fmt.Printf("Author: %s (%s)\n", pi.Creator.Name, pi.Creator.Link)
	}

	if download != "" {
		if !pi.HasImage(photo.Original) {
			panic("No image available to download")
		}

		item := pi.GetImage(photo.Original)
		if err != nil {
			panic(err)
		}

		r, err := item.Reader()
		if err != nil {
			panic(err)
		}

		f, err := os.Create(download)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		if err != nil {
			panic(err)
		}
	}
}
