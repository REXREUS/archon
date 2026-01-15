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

var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Generate architecture diagrams (Mermaid/PlantUML)",
	Run: func(cmd *cobra.Command, args []string) {
		diagType, _ := cmd.Flags().GetString("type")
		focus, _ := cmd.Flags().GetString("focus")

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
			contextText, _ = orchestrator.SearchContext(ctx, "structure and relationships for "+focus)
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		prompt := fmt.Sprintf("%s\n\nTask: Create %s diagram code based on the current code structure. Focus on: %s", 
			contextText, diagType, focus)

		fmt.Println("Generating diagram...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nGenerated Diagram Code (%s):\n%s\n", diagType, resp.Text)
	},
}

func init() {
	diagramCmd.Flags().String("type", "mermaid", "Type of diagram (mermaid, plantuml)")
	diagramCmd.Flags().String("focus", "all", "Focus area for the diagram")
	rootCmd.AddCommand(diagramCmd)
}
