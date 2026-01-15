package cli

import (
	"archon/internal/config"
	"archon/internal/adapters/vectordb"
	"archon/internal/core"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	watch bool
	force bool
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Scan and index the codebase",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Indexing codebase...")
		ctx := context.Background()

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.GeminiKey == "" {
			fmt.Println("Error: Gemini API key not found. Use 'archon auth' to set it.")
			return
		}

		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		if err != nil {
			fmt.Printf("Error creating store: %v\n", err)
			return
		}
		defer store.Close()

		if force {
			fmt.Println("Force flag set, clearing existing index...")
			err = store.Clear(ctx)
			if err != nil {
				fmt.Printf("Error clearing store: %v\n", err)
				return
			}
		}

		orchestrator := core.NewOrchestrator(store)
		err = orchestrator.IndexDirectory(ctx, ".", func(current, total int, file string) {
			fmt.Printf("[%d/%d] Indexing %s...\n", current, total, file)
		})
		if err != nil {
			fmt.Printf("Error indexing: %v\n", err)
			return
		}
		fmt.Println("Indexing complete.")

		if watch {
			err = orchestrator.WatchDirectory(ctx, ".")
			if err != nil {
				fmt.Printf("Error watching directory: %v\n", err)
				return
			}
		}
	},
}

func init() {
	indexCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for file changes")
	indexCmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-indexing of all files")
	rootCmd.AddCommand(indexCmd)
}
