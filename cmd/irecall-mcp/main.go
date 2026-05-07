package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	irecallmcp "github.com/gigol/irecall/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// Injected at link time: go build -ldflags "-X main.version=v0.1.0"
var version = "dev"

func main() {
	os.Exit(run(os.Args[0], os.Args[1:], os.Stdout, os.Stderr))
}

func run(program string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet(program, flag.ContinueOnError)
	fs.SetOutput(stderr)
	baseURLFlag := fs.String("base-url", "", "override the iRecall web API base URL (default: IRECALL_BASE_URL or http://127.0.0.1:9527)")
	tokenFileFlag := fs.String("token-file", "", "read the API token from this file (preferred over IRECALL_API_TOKEN)")
	timeoutFlag := fs.Duration("timeout", 15*time.Second, "HTTP timeout for calls to the iRecall web API")
	versionFlag := fs.Bool("version", false, "print version and exit")
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), usageText(fs, program))
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *versionFlag {
		fmt.Fprintln(stdout, "iRecall MCP", binaryVersion())
		return 0
	}

	cfg, err := irecallmcp.LoadConfig(*baseURLFlag, *tokenFileFlag, *timeoutFlag)
	if err != nil {
		fmt.Fprintf(stderr, "irecall-mcp: %v\n", err)
		return 1
	}

	srv, err := irecallmcp.NewServer(cfg, binaryVersion())
	if err != nil {
		fmt.Fprintf(stderr, "irecall-mcp: %v\n", err)
		return 1
	}

	if err := mcpserver.ServeStdio(srv); err != nil {
		fmt.Fprintf(stderr, "irecall-mcp: %v\n", err)
		return 1
	}
	return 0
}

func usageText(fs *flag.FlagSet, program string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Usage: %s [flags]\n\n", program)
	buf.WriteString("irecall-mcp exposes the local iRecall web API as MCP tools over stdio.\n\n")
	buf.WriteString("Environment:\n")
	buf.WriteString("  IRECALL_BASE_URL   Base URL for the iRecall web server (default: http://127.0.0.1:9527)\n")
	buf.WriteString("  IRECALL_API_TOKEN_FILE  Path to a protected file containing the bearer token\n")
	buf.WriteString("  IRECALL_API_TOKEN       Bearer token used for authenticated REST requests\n\n")
	buf.WriteString("Flags:\n")
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&buf, "  -%s\n    \t%s\n", f.Name, f.Usage)
	})
	buf.WriteString("\nExamples:\n")
	fmt.Fprintf(&buf, "  %s --token-file ~/.config/irecall/mcp-api-token\n", program)
	fmt.Fprintf(&buf, "  IRECALL_API_TOKEN=... %s\n", program)
	fmt.Fprintf(&buf, "  IRECALL_BASE_URL=http://127.0.0.1:9527 IRECALL_API_TOKEN_FILE=~/.config/irecall/mcp-api-token %s\n", program)
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
