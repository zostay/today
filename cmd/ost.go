package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/bbrks/wrap"
	"github.com/markusmobius/go-dateparser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zostay/today/cmd/flag"
	"github.com/zostay/today/cmd/output"
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

	on       flag.Date
	download string
)

func init() {
	ostCmd.AddCommand(ostTodayCmd, ostOnCmd, ostPhotoCmd)

	ostCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
	ostCmd.Flags().VarP(&outputFormat, "output", "o", fmt.Sprintf("Output format (%s)", flag.ListOutputFormats()))

	ostTodayCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostTodayCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostTodayCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
	ostTodayCmd.Flags().VarP(&outputFormat, "output", "o", fmt.Sprintf("Output format (%s)", flag.ListOutputFormats()))

	ostOnCmd.Flags().BoolVarP(&asHtml, "html", "H", false, "Output as HTML")
	ostOnCmd.Flags().BoolVarP(&asMeta, "meta", "m", false, "Output information about Scripture")
	ostOnCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
	ostOnCmd.Flags().VarP(&outputFormat, "output", "o", fmt.Sprintf("Output format (%s)", flag.ListOutputFormats()))

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

	ofmt := collectOutputFormat().Name
	var v string
	switch ofmt {
	case output.YAMLFormat:
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
	case output.JSONFormat:
		vv, err := client.TodayVerse(cmd.Context(), opts...)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(cmd.OutOrStdout())
		err = enc.Encode(vv)
		if err != nil {
			panic(err)
		}
		return
	case output.MetaFormat:
		vv, err := client.TodayVerse(cmd.Context(), opts...)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Reference: %s\n", vv.Reference)
		fmt.Printf("Version:   %s\n", vv.Version.Name)
		fmt.Printf("Link:      %s\n", vv.Version.Link)
		return
	case output.HTMLFormat:
		var vh template.HTML
		vh, err = client.TodayHTML(cmd.Context(), opts...)
		v = string(vh)
	case output.TextFormat:
		v, err = client.Today(cmd.Context(), opts...)
	case output.JPEGFormat:
		panic("JPEG output not implemented")
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

	ofmt := collectOutputFormat().Name
	switch ofmt {
	case output.YAMLFormat:
		enc := yaml.NewEncoder(cmd.OutOrStdout())
		err = enc.Encode(pi)
		if err != nil {
			panic(err)
		}
	case output.JSONFormat:
		enc := json.NewEncoder(cmd.OutOrStdout())
		err = enc.Encode(pi)
		if err != nil {
			panic(err)
		}
	case output.TextFormat, output.MetaFormat:
		fmt.Printf("Link: %s\n", pi.Link)
		fmt.Printf("Author: %s (%s)\n", pi.Creator.Name, pi.Creator.Link)
	case output.HTMLFormat:
		panic("HTML output not supported for photos")
	case output.JPEGFormat:
		if download == "" {
			download = "download.jpg"
		}
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
