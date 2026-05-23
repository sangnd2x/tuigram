package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIID   int    `toml:"api_id"`
	APIHash string `toml:"api_hash"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	if cfg.APIID == 0 || cfg.APIHash == "" {
		return nil, fmt.Errorf("config: api_id and api_hash are required in %s", path)
	}
	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

func LoadOrCreate(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("No config file found. Creating one at:", path)
		fmt.Println("Get your API ID and API Hash from https://my.telegram.org")
		var cfg Config
		fmt.Print("Enter API ID: ")
		if _, err := fmt.Scan(&cfg.APIID); err != nil {
			return nil, fmt.Errorf("reading api_id: %w", err)
		}
		fmt.Print("Enter API Hash: ")
		if _, err := fmt.Scan(&cfg.APIHash); err != nil {
			return nil, fmt.Errorf("reading api_hash: %w", err)
		}
		if err := Save(path, &cfg); err != nil {
			return nil, fmt.Errorf("saving config: %w", err)
		}
		fmt.Println("Config saved to", path)
		return &cfg, nil
	}
	return Load(path)
}
