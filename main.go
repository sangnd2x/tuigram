package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	appauth "github.com/sangnguyen/tui-telegram/auth"
	"github.com/sangnguyen/tui-telegram/config"
	"github.com/sangnguyen/tui-telegram/internal/xdg"
	"github.com/sangnguyen/tui-telegram/telegram"
	"github.com/sangnguyen/tui-telegram/ui"
)

func main() {
	cfgPath := filepath.Join(xdg.ConfigDir(), "config.toml")
	cfg, err := config.LoadOrCreate(cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		os.Exit(1)
	}

	phone, err := appauth.PromptPhone()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading phone:", err)
		os.Exit(1)
	}

	sessionPath := filepath.Join(xdg.DataDir(), "session.json")
	if err := os.MkdirAll(filepath.Dir(sessionPath), 0700); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating data dir:", err)
		os.Exit(1)
	}

	bridge := telegram.NewUpdateBridge()
	client, err := telegram.New(cfg, sessionPath, bridge)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating client:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.Run(ctx, appauth.New(phone), func(ctx context.Context) error {
		p := tea.NewProgram(ui.New(client, ctx), tea.WithAltScreen())
		bridge.SetProgram(p)
		_, err := p.Run()
		cancel()
		return err
	}); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
