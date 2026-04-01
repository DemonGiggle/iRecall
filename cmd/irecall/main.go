package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
	"github.com/gigol/irecall/tui"
)

// Injected at link time: go build -ldflags "-X main.version=v0.1.0"
var version = "dev"

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println("iRecall", version)
		os.Exit(0)
	}

	if err := config.EnsureDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall: cannot create data directories: %v\n", err)
		os.Exit(1)
	}

	// Set up file-based structured logging (never log to stdout/stderr — TUI owns them).
	logLevel := slog.LevelInfo
	if *debugFlag {
		logLevel = slog.LevelDebug
	}
	logFile, err := os.OpenFile(config.LogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall: cannot open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: logLevel})))

	// Open database.
	store, err := db.Open(config.DBPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall: cannot open database: %v\n", err)
		os.Exit(1)
	}

	// Bootstrap engine with default settings; load persisted settings from DB.
	defaults := core.DefaultSettings()
	engine := core.New(store, defaults)
	settings, err := engine.LoadSettings(nil) //nolint: staticcheck
	if err != nil || settings == nil {
		settings = defaults
	}
	engine.UpdateSettings(settings)

	// Start TUI.
	app := tui.NewApp(engine, settings, 0, 0)
	p := tea.NewProgram(app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall: %v\n", err)
		engine.Close()
		os.Exit(1)
	}
	engine.Close()
}
