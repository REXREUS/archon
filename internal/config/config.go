package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GeminiKey   string `mapstructure:"gemini_key"`
	ModelID     string `mapstructure:"model_id"`
	ProjectHash string `mapstructure:"project_hash"`
	CacheName   string `mapstructure:"cache_name"`
}

func LoadConfig() (*Config, error) {
	viper.Reset()

	// Coba cari file .archon.yaml di folder saat ini atau home
	home, _ := os.UserHomeDir()
	searchPaths := []string{".", home}
	
	for _, p := range searchPaths {
		if p == "" {
			continue
		}
		path := filepath.Join(p, ".archon.yaml")
		if _, err := os.Stat(path); err == nil {
			viper.SetConfigFile(path)
			break
		}
	}

	// Jika tidak ditemukan file spesifik, atur default search
	if viper.ConfigFileUsed() == "" {
		viper.SetConfigName(".archon")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if home != "" {
			viper.AddConfigPath(home)
		}
	}

	viper.SetEnvPrefix("ARCHON")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.ModelID == "" {
		cfg.ModelID = "gemini-3-pro-preview"
	}

	return &cfg, nil
}
