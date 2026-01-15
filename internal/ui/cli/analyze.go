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

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform deep architectural analysis",
	Run: func(cmd *cobra.Command, args []string) {
		depth, _ := cmd.Flags().GetString("depth")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		store, _ := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var contextText string
		if store != nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			contextText, _ = orchestrator.SearchContext(ctx, "architectural overview and anomalies")
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		prompt := fmt.Sprintf("%s\n\nTask: Perform a deep architectural analysis on this project. Detect anomalies, code smells, or design pattern violations. Depth: %s", 
			contextText, depth)

		fmt.Println("Analyzing architecture...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nArchitectural Analysis:\n%s\n", resp.Text)
		fmt.Printf("\n(Tokens used: %d)\n", resp.TotalTokens)
	},
}

func init() {
	analyzeCmd.Flags().String("depth", "full", "Depth of analysis (basic, full)")
	rootCmd.AddCommand(analyzeCmd)
}
