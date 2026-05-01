//go:build !wails

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	charmterm "github.com/charmbracelet/x/term"
	"github.com/gigol/irecall/app"
	"github.com/gigol/irecall/config"
	frontendassets "github.com/gigol/irecall/frontend"
)

// Injected at link time: go build -ldflags "-X main.version=v0.1.0"
var version = "dev"

func main() {
	maybeHandleAuthCommand(os.Args[1:])
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	dataPathFlag := flag.String("data-path", "", "store database, config, and logs under this root path")
	hostFlag := flag.String("host", "0.0.0.0", "host/interface to bind the web server to")
	portFlag := flag.Int("port", 0, "port to listen on (overrides saved web port)")
	apiOnlyFlag := flag.Bool("api-only", false, "run a headless API server without serving the frontend UI or requiring a web password at startup")
	providerOptions := bindProviderStartupFlags(flag.CommandLine)
	resetPasswordFlag := flag.Bool("reset-passwd", false, "clear the configured web password and prompt for a new one before startup")
	versionFlag := flag.Bool("version", false, "print version and exit")
	// WARNING: This flag disables the interactive web password check. For testing only.
	// Do NOT enable this in production environments.
	unsafeNoPasswordCheckFlag := flag.Bool("unsafe-no-password-check", false, "DISABLE web password check (testing only). Do NOT use in production.")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usageText(flag.CommandLine, os.Args[0]))
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println("iRecall web", binaryVersion())
		return
	}

	serverOptions := ServerOptions{
		APIOnly:               *apiOnlyFlag,
		UnsafeNoPasswordCheck: *unsafeNoPasswordCheckFlag,
	}
	if err := serverOptions.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
	if err := providerOptions.Validate(serverOptions); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}

	if *dataPathFlag != "" {
		config.SetRootPath(*dataPathFlag)
	} else if _, err := config.ApplyPreferredRootPath(); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: cannot load preferred data root: %v\n", err)
		os.Exit(1)
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

	runtimeApp, err := app.NewApp(config.RootPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
	defer runtimeApp.Shutdown(nil)
	if err := applyAPIOnlyProviderStartupConfig(runtimeApp, serverOptions, *providerOptions); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
	if *resetPasswordFlag {
		if err := runtimeApp.ResetPassword(); err != nil {
			fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Existing web password cleared.")
	}
	if serverOptions.requiresWebPasswordBootstrap() {
		if err := ensureWebPasswordConfigured(runtimeApp); err != nil {
			fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
			os.Exit(1)
		}
	} else if serverOptions.UnsafeNoPasswordCheck {
		fmt.Fprintln(os.Stderr, "WARNING: running with --unsafe-no-password-check; web password check is disabled (testing only). Do NOT use in production.")
	} else if serverOptions.APIOnly {
		fmt.Fprintln(os.Stderr, "Running in --api-only mode; frontend UI and browser-session auth are disabled.")
	}

	port := *portFlag
	if port == 0 && runtimeApp.GetSettings() != nil {
		port = runtimeApp.GetSettings().Web.Port
	}
	if port < 1 || port > 65535 {
		port = 9527
	}

	server, err := NewServer(runtimeApp, frontendassets.Assets, port, serverOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}

	addr := net.JoinHostPort(strings.TrimSpace(*hostFlag), fmt.Sprintf("%d", port))
	fmt.Printf("%s listening on http://%s\n", serverOptions.listenerLabel(), addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-web: %v\n", err)
		os.Exit(1)
	}
}

func usageText(fs *flag.FlagSet, program string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Usage: %s [flags]\n", program)
	fmt.Fprintf(&buf, "       %s auth <subcommand> [flags]\n\n", program)
	buf.WriteString("iRecall web serves the local HTTP UI and API.\n\n")
	buf.WriteString("Flags:\n")
	buf.WriteString(flagDefaultsText(fs))
	buf.WriteString("\nExamples:\n")
	fmt.Fprintf(&buf, "  %s --version\n", program)
	fmt.Fprintf(&buf, "  %s -host 127.0.0.1 -port 9527\n", program)
	fmt.Fprintf(&buf, "  %s --api-only --provider-host api.openai.example/v1 --provider-port 443 --provider-https --provider-model gpt-4.1 --provider-keyword-model gpt-4.1-mini\n", program)
	fmt.Fprintf(&buf, "  %s auth issue-token --write-token-file ~/.config/irecall/mcp-api-token\n", program)
	return buf.String()
}

func flagDefaultsText(fs *flag.FlagSet) string {
	var buf bytes.Buffer
	originalOutput := fs.Output()
	fs.SetOutput(&buf)
	defer fs.SetOutput(originalOutput)
	fs.PrintDefaults()
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

func ensureWebPasswordConfigured(app *app.App) error {
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
