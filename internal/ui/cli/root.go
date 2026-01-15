package cli

import (
	"fmt"
	"os"

	"archon/internal/ui/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "archon",
	Short: "ArchonCLI - AI Architect Assistant for your codebase",
	Long: `ArchonCLI is a revolutionary CLI & TUI tool designed to interact with complex codebases 
using semantic syntax-aware indexing and Google Gemini 3.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments, start TUI mode
		if len(args) == 0 {
			tui.StartApp()
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
