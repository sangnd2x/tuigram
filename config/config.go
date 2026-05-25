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
	Phone   string `toml:"phone,omitempty"`
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

func ensurePhone(path string, cfg *Config) error {
	if cfg.Phone == "" {
		fmt.Print("Enter your phone number (with country code, e.g. +1234567890): ")
		if _, err := fmt.Scan(&cfg.Phone); err != nil {
			return fmt.Errorf("reading phone: %w", err)
		}
		if err := Save(path, cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
	}
	return nil
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
		fmt.Print("Enter your phone number (with country code, e.g. +1234567890): ")
		if _, err := fmt.Scan(&cfg.Phone); err != nil {
			return nil, fmt.Errorf("reading phone: %w", err)
		}
		if err := Save(path, &cfg); err != nil {
			return nil, fmt.Errorf("saving config: %w", err)
		}
		fmt.Println("Config saved to", path)
		return &cfg, nil
	}
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}
	if err := ensurePhone(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
