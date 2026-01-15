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
	"github.com/spf13/viper"
)

var askCmd = &cobra.Command{
	Use:   "ask [query]",
	Short: "Ask a question about the codebase",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		
		ctx := context.Background()
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if cfg.GeminiKey == "" {
			fmt.Println("Error: Gemini API key not found. Use 'archon auth' or set ARCHON_GEMINI_KEY environment variable.")
			os.Exit(1)
		}

		// Initialize Vector DB and Orchestrator for RAG
		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		if err != nil {
			fmt.Printf("Warning: Vector DB not initialized. Asking without context. (%v)\n", err)
		} else {
			defer store.Close()
		}

		var prompt string
		if store != nil {
			orchestrator := core.NewOrchestrator(store)
			fmt.Printf("Searching context...\n")
			contextText, err := orchestrator.SearchContext(ctx, query)
			if err != nil {
				fmt.Printf("Error searching context: %v\n", err)
				prompt = query
			} else {
				prompt = fmt.Sprintf("%s\n\nUser Question: %s", contextText, query)
			}
		} else {
			prompt = query
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error creating Gemini client: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		// Sync and Use Context Cache
		hash, err := gemini.CalculateProjectHash(".")
		if err == nil {
			if cfg.CacheName != "" && cfg.ProjectHash == hash {
				client.SetCachedContent(cfg.CacheName)
			} else {
				// Try to create new cache if possible
				orchestrator := core.NewOrchestrator(store)
				files, err := orchestrator.GetFilesForIndexing(".")
				if err == nil && len(files) > 0 {
					cm := gemini.NewCacheManager(client.Client())
					cacheName, err := cm.CreateContextCache(ctx, cfg.ModelID, files)
					if err == nil {
						client.SetCachedContent(cacheName)
						// Save to config
						viper.Set("project_hash", hash)
						viper.Set("cache_name", cacheName)
						viper.WriteConfig()
					} else {
						// Update hash anyway to prevent constant retries if it's too small
						viper.Set("project_hash", hash)
						viper.Set("cache_name", "")
						viper.WriteConfig()
					}
				}
			}
		}

		fmt.Printf("Thinking...\n")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error asking Gemini: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nResponse:\n%s\n", resp.Text)
		fmt.Printf("\n(Tokens used: %d)\n", resp.TotalTokens)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
