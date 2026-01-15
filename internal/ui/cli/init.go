package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Archon configuration",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".archon.yaml")
		
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Configuration file already exists.")
			return
		}

		viper.Set("model_id", "gemini-3-pro-preview")
		err := viper.SafeWriteConfigAs(configPath)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			return
		}
		fmt.Printf("Initialized configuration at %s\n", configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
