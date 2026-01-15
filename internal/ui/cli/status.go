package cli

import (
	"fmt"
	"os"
	"archon/internal/config"
	"archon/internal/adapters/gemini"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Archon status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		fmt.Printf("Archon Status:\n")
		fmt.Printf("- Model: %s\n", cfg.ModelID)
		if cfg.GeminiKey != "" {
			fmt.Printf("- API Key: Configured\n")
		} else {
			fmt.Printf("- API Key: NOT Configured\n")
		}
		
		if _, err := os.Stat("./chromem_db"); err == nil {
			fmt.Printf("- Vector DB: Ready (chromem_db)\n")
		} else {
			fmt.Printf("- Vector DB: Not initialized (Use 'archon index')\n")
		}

		// Caching status
		hash, _ := gemini.CalculateProjectHash(".")
		if hash != "" {
			fmt.Printf("- Project Hash: %s\n", hash[:8]+"...")
		}
		
		cacheStatus := "Inactive"
		if cfg.CacheName != "" {
			if cfg.ProjectHash == hash {
				cacheStatus = fmt.Sprintf("Active (%s)", cfg.CacheName)
			} else {
				cacheStatus = "Inactive (Hash Mismatch)"
			}
		} else {
			cacheStatus = "Inactive (Not Created)"
		}
		fmt.Printf("- Context Cache: %s\n", cacheStatus)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
