package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ArchonCLI v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
