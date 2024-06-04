package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/zostay/today/cmd/flag"
	"github.com/zostay/today/pkg/ost"
	"github.com/zostay/today/pkg/photo"
	"gopkg.in/yaml.v3"
)

var (
	ostPhotoCmd = &cobra.Command{
		Use:   "photo",
		Short: "Show information about or download the photo used with the scripture of the day from openscripture.today",
		Args:  cobra.NoArgs,
		Run:   RunOstPhoto,
	}

	download string
	on       flag.Date
)

func init() {
	ostPhotoCmd.Flags().StringVarP(&download, "download", "d", "openscripture.jpg", "Download the file photo to the named file")
	ostPhotoCmd.Flags().VarP(&on, "on", "o", "Specify the date to get the photo for")
	ostPhotoCmd.Flags().BoolVarP(&asYaml, "yaml", "y", false, "Output as YAML")
}

func RunOstPhoto(cmd *cobra.Command, args []string) {
	opts := []ost.DayOption{}
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
