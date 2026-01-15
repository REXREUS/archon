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

var explainCmd = &cobra.Command{
	Use:   "explain [file/symbol]",
	Short: "Explain a file or symbol",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		
		cfg, _ := config.LoadConfig()
		ctx := context.Background()

		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var contextText string
		if err == nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			fmt.Println("Gathering context...")
			contextText, _ = orchestrator.SearchContext(ctx, "Explain "+target)
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		fmt.Printf("Analyzing %s...\n", target)
		
		var prompt string
		if contextText != "" {
			prompt = fmt.Sprintf("%s\n\nTask: Explain in detail about: %s", contextText, target)
		} else {
			// Fallback if no context found
			content, err := os.ReadFile(target)
			if err == nil {
				prompt = fmt.Sprintf("Explain the following code:\n```\n%s\n```", string(content))
			} else {
				prompt = fmt.Sprintf("Explain about: %s", target)
			}
		}

		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(resp.Text)
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
