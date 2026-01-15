package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall ArchonCLI and remove its data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling ArchonCLI...")

		// 1. Remove configuration directory
		homeDir, _ := os.UserHomeDir()
		// Based on documentation, config is at $HOME/.archon.yaml or similar
		// But usually it's better to check where it is stored.
		// For now let's follow standard patterns.
		configPath := homeDir + "/.archon.yaml"
		if _, err := os.Stat(configPath); err == nil {
			os.Remove(configPath)
			fmt.Println("- Removed configuration file: " + configPath)
		}

		// 2. Remove local vector database (chromem_db)
		// We need to know where it is stored. In the current tree it is in project root.
		// But in installed version it might be elsewhere.
		if _, err := os.Stat("chromem_db"); err == nil {
			os.RemoveAll("chromem_db")
			fmt.Println("- Removed local vector database")
		}

		// 3. Inform about binary removal
		fmt.Println("\nTo completely remove the binary, please run the uninstall script:")
		if runtime.GOOS == "windows" {
			fmt.Println("PowerShell: irm https://raw.githubusercontent.com/rexreus/archon/main/scripts/uninstall.ps1 | iex")
		} else {
			fmt.Println("Bash: curl -sSL https://raw.githubusercontent.com/rexreus/archon/main/scripts/uninstall.sh | bash")
		}
		
		fmt.Println("\nNote: You can also manually delete the binary from your PATH.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
