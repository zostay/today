package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Args:  cobra.NoArgs,
	Run:   RunVersion,
}

func init() {
}

//go:embed version.txt
var Version string

func RunVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("today v%s", Version)
}
