package main

import (
	"bytes"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	appbackend "github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/tui"
)

// Injected at link time: go build -ldflags "-X main.version=v0.1.0"
var version = "dev"

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	dataPathFlag := flag.String("data-path", "", "store database, config, and logs under this root path")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usageText(flag.CommandLine, os.Args[0]))
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println("iRecall", binaryVersion())
		os.Exit(0)
	}

	if *dataPathFlag != "" {
		config.SetRootPath(*dataPathFlag)
	} else if _, err := config.ApplyPreferredRootPath(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall: cannot load preferred data root: %v\n", err)
		os.Exit(1)
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

	runtimeState, err := appbackend.OpenRuntime(config.RootPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall: cannot open runtime: %v\n", err)
		os.Exit(1)
	}

	// Start TUI.
	app := tui.NewApp(runtimeState.Engine, runtimeState.Settings, runtimeState.Profile, runtimeState.Paths, 0, 0)
	p := tea.NewProgram(app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall: %v\n", err)
		runtimeState.Engine.Close()
		os.Exit(1)
	}
	runtimeState.Engine.Close()
}

func usageText(fs *flag.FlagSet, program string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Usage: %s [flags]\n\n", program)
	buf.WriteString("iRecall is a terminal-first personal knowledge recall app.\n")
	buf.WriteString("It stores quotes locally, retrieves reference quotes for recall, and supports manual quote sharing via exported JSON.\n\n")
	buf.WriteString("Flags:\n")
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&buf, "  -%s\n    \t%s\n", f.Name, f.Usage)
	})
	buf.WriteString("\nNotes:\n")
	buf.WriteString("  - On first launch, iRecall asks for your display name so shared quotes can show their source.\n")
	buf.WriteString("  - Use -data-path to run isolated local instances on the same machine.\n")
	buf.WriteString("\nExamples:\n")
	fmt.Fprintf(&buf, "  %s\n", program)
	fmt.Fprintf(&buf, "  %s -debug\n", program)
	fmt.Fprintf(&buf, "  %s --version\n", program)
	fmt.Fprintf(&buf, "  %s -data-path /tmp/irecall-alice\n", program)
	fmt.Fprintf(&buf, "  %s -data-path /tmp/irecall-bob\n", program)
	return buf.String()
}

func binaryVersion() string {
	if v := strings.TrimSpace(version); v != "" && v != "dev" {
		return v
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	var tag string
	var revision string
	var modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.tag":
			tag = strings.TrimSpace(setting.Value)
		case "vcs.revision":
			revision = strings.TrimSpace(setting.Value)
		case "vcs.modified":
			modified = strings.TrimSpace(setting.Value)
		}
	}

	if tag != "" {
		if modified == "true" {
			return tag + "-dirty"
		}
		return tag
	}

	if revision != "" {
		if len(revision) > 12 {
			revision = revision[:12]
		}
		if modified == "true" {
			return revision + "-dirty"
		}
		return revision
	}

	return "dev"
}
