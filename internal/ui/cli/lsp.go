package cli

import (
	"archon/internal/ui/lsp"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Start Language Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Starting Archon LSP server...")
		server := lsp.NewServer()
		if err := server.Start(); err != nil {
			fmt.Printf("LSP server error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
