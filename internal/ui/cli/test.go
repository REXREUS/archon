package cli

import (
	"archon/internal/adapters/gemini"
	"archon/internal/adapters/vectordb"
	"archon/internal/config"
	"archon/internal/core"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [file]",
	Short: "Generate unit tests for a file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var contextText string
		if err == nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			fmt.Println("Gathering context...")
			contextText, _ = orchestrator.SearchContext(ctx, "unit test for "+filePath)
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		content, _ := os.ReadFile(filePath)
		prompt := fmt.Sprintf("%s\n\nTask: Create comprehensive unit tests for the following file: %s\n\nCode:\n```\n%s\n```", 
			contextText, filePath, string(content))

		fmt.Println("Generating tests...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nGenerated Unit Tests for %s:\n%s\n", filePath, resp.Text)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
