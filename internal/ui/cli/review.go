package cli

import (
	"archon/internal/adapters/gemini"
	"archon/internal/adapters/vectordb"
	"archon/internal/config"
	"archon/internal/core"
	"archon/internal/utils"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review staged changes using AI",
	Long:  `Analyzes staged code changes (git add) and provides feedback regarding quality, potential bugs, and architectural compliance.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		// Ambil diff staged
		diffOutput, err := utils.GetStagedDiff()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(diffOutput) == "" {
			fmt.Println("No staged changes (git add) to review.")
			return
		}

		// Inisialisasi store untuk RAG (opsional tapi bagus untuk konteks)
		store, _ := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var contextText string
		if store != nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			// Cari konteks berdasarkan file yang berubah
			contextText, _ = orchestrator.SearchContext(ctx, "Review changes in these files")
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		// Gunakan cache jika tersedia
		hash, err := gemini.CalculateProjectHash(".")
		if err == nil && cfg.CacheName != "" && cfg.ProjectHash == hash {
			client.SetCachedContent(cfg.CacheName)
		}

		prompt := fmt.Sprintf(`%s

Task: Perform a deep code review on the following changes (diff). 
Focus on:
1. Potential bugs or missed edge cases.
2. Compliance with best practices (Clean Code, SOLID).
3. Code smells or unnecessary complexity.
4. Concrete improvement suggestions.

Diff:
%s`, contextText, string(diffOutput))

		fmt.Println("🚀 Analyzing your changes...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nAI Code Review:\n%s\n", resp.Text)
		fmt.Printf("\n(Tokens used: %d)\n", resp.TotalTokens)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
