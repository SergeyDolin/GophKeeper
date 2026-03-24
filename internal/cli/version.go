package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	build   = "dev"
	date    = "24.03.2026"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", version)
		fmt.Println("Build:", build)
		fmt.Println("Date:", date)
	},
}
