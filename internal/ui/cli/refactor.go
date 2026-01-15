package cli

import (
	"archon/internal/adapters/gemini"
	"archon/internal/adapters/vectordb"
	"archon/internal/config"
	"archon/internal/core"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
	Use:   "refactor [file]",
	Short: "Analyze and suggest refactorings for a file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		goal, _ := cmd.Flags().GetString("goal")
		apply, _ := cmd.Flags().GetBool("apply")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var contextText string
		if err == nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			fmt.Println("Gathering context...")
			contextText, _ = orchestrator.SearchContext(ctx, "refactor "+filePath+" with goal "+goal)
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		content, _ := os.ReadFile(filePath)
		
		prompt := fmt.Sprintf("%s\n\nTask: Perform refactoring on the following file: %s\nGoal: %s\n\nCode:\n```\n%s\n```", 
			contextText, filePath, goal, string(content))

		if apply {
			prompt += "\n\nRETURN ONLY THE REFACTORED CODE in a single Markdown code block (```). Do not provide any explanation outside that code block because your output will be written directly to the file."
		}

		fmt.Println("Analyzing and refactoring...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if apply {
			// Ekstrak kode dari blok ```
			newCode := extractCode(resp.Text)
			if newCode == "" {
				fmt.Println("Failed to extract code from AI response. Displaying text suggestions only:")
				fmt.Println(resp.Text)
				return
			}

			err = os.WriteFile(filePath, []byte(newCode), 0644)
			if err != nil {
				fmt.Printf("Failed to write to file: %v\n", err)
				return
			}
			fmt.Printf("✅ Successfully applied refactoring to %s\n", filePath)
		} else {
			fmt.Printf("\nRefactoring Suggestions:\n%s\n", resp.Text)
			fmt.Printf("\n(Tokens used: %d)\n", resp.TotalTokens)
		}
	},
}

func extractCode(resp string) string {
	lines := strings.Split(resp, "\n")
	var codeLines []string
	inBlock := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inBlock {
				break // End of block
			}
			inBlock = true
			continue
		}
		if inBlock {
			codeLines = append(codeLines, line)
		}
	}
	if len(codeLines) == 0 {
		// Jika tidak ada blok ```, coba ambil semua teks (siapa tahu Gemini lupa bloknya)
		// Tapi lebih aman return kosong jika format tidak sesuai
		return ""
	}
	return strings.Join(codeLines, "\n")
}

func init() {
	refactorCmd.Flags().String("goal", "improve code quality and performance", "Specific goal for refactoring")
	refactorCmd.Flags().Bool("apply", false, "Apply refactoring directly to the file")
	rootCmd.AddCommand(refactorCmd)
}
