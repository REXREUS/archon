package cli

import (
	"archon/internal/adapters/gemini"
	"archon/internal/config"
	"archon/internal/utils"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a smart commit message",
	Long:  `Analyzes staged changes and generates a commit message following conventions (such as Conventional Commits).`,
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
			fmt.Println("No staged changes to generate a commit message for.")
			return
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		prompt := fmt.Sprintf(`Task: Create a descriptive commit message based on the following code changes (diff).
Use Conventional Commits format (type(scope): description).
Provide a brief explanation of WHAT changed and WHY (if it can be inferred).
Only return the commit message itself, no other additional text.

Diff:
%s`, string(diffOutput))

		fmt.Println("🤖 Generating commit message...")
		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		commitMsg := strings.TrimSpace(resp.Text)
		fmt.Printf("\nSuggested Commit Message:\n---\n%s\n---\n", commitMsg)
		
		fmt.Print("\nDo you want to commit now? (y/n): ")
		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "y" {
			commitExec := exec.Command("git", "commit", "-m", commitMsg)
			commitExec.Stdout = os.Stdout
			commitExec.Stderr = os.Stderr
			err := commitExec.Run()
			if err != nil {
				fmt.Printf("Failed to commit: %v\n", err)
			} else {
				fmt.Println("✅ Successfully committed!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
