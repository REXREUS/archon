package cli

import (
	"archon/internal/config"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage API credentials",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("Please provide an API key with --key")
			return
		}

		// Load existing config if any
		_, _ = config.LoadConfig()

		viper.Set("gemini_key", key)
		
		// If no config file found by LoadConfig, default to .archon.yaml in current dir
		if viper.ConfigFileUsed() == "" {
			viper.SetConfigFile(".archon.yaml")
		}

		err := viper.WriteConfig()
		if err != nil {
			// If WriteConfig fails (e.g. file doesn't exist), try SafeWriteConfig
			err = viper.SafeWriteConfig()
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}
		}
		fmt.Printf("API key saved successfully to %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	authCmd.Flags().String("key", "", "Google Gemini API Key")
	rootCmd.AddCommand(authCmd)
}
