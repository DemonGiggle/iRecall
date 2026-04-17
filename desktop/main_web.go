//go:build !wails

package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	charmterm "github.com/charmbracelet/x/term"
	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/desktop/backend"
	"github.com/gigol/irecall/desktop/web"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	dataPathFlag := flag.String("data-path", "", "store database, config, and logs under this root path")
	hostFlag := flag.String("host", "0.0.0.0", "host/interface to bind the web server to")
	portFlag := flag.Int("port", 0, "port to listen on (overrides saved web port)")
	flag.Parse()

	if *dataPathFlag != "" {
		config.SetRootPath(*dataPathFlag)
	}
	if err := config.EnsureDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: cannot create data directories: %v\n", err)
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	if *debugFlag {
		logLevel = slog.LevelDebug
	}
	logFile, err := os.OpenFile(config.LogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: cannot open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: logLevel})))

	app, err := backend.NewApp(config.RootPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
	defer app.Shutdown(nil)
	if err := ensureWebPasswordConfigured(app); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}

	port := *portFlag
	if port == 0 && app.GetSettings() != nil {
		port = app.GetSettings().Web.Port
	}
	if port < 1 || port > 65535 {
		port = 9527
	}

	server, err := web.NewServer(app, assets, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}

	addr := net.JoinHostPort(strings.TrimSpace(*hostFlag), fmt.Sprintf("%d", port))
	fmt.Printf("iRecall web UI listening on http://%s\n", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
}

func ensureWebPasswordConfigured(app *backend.App) error {
	status, err := app.AuthStatus()
	if err != nil {
		return err
	}
	if status.PasswordConfigured {
		return nil
	}
	if !charmterm.IsTerminal(os.Stdin.Fd()) {
		return errors.New("web password is not configured; launch from an interactive terminal to complete first-time setup")
	}

	fmt.Println("Create the iRecall web password before the server starts listening.")
	fmt.Println("Password policy: at least 12 characters and at least 3 of uppercase, lowercase, digit, symbol.")
	for {
		password, err := readPasswordPrompt("New password: ")
		if err != nil {
			return err
		}
		confirm, err := readPasswordPrompt("Confirm password: ")
		if err != nil {
			return err
		}
		if err := app.SetupPassword(password, confirm); err != nil {
			fmt.Fprintf(os.Stderr, "Password setup failed: %v\n", err)
			continue
		}
		fmt.Println("Web password configured.")
		return nil
	}
}

func readPasswordPrompt(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := charmterm.ReadPassword(os.Stdin.Fd())
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}
