package main

import (
	"bytes"
	"flag"
	"fmt"
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
	baseURLFlag := flag.String("base-url", "", "override the iRecall web API base URL (default: IRECALL_BASE_URL or http://127.0.0.1:9527)")
	timeoutFlag := flag.Duration("timeout", 15*time.Second, "HTTP timeout for calls to the iRecall web API")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usageText(flag.CommandLine, os.Args[0]))
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println("iRecall MCP", binaryVersion())
		return
	}

	cfg, err := irecallmcp.LoadConfig(*baseURLFlag, *timeoutFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-mcp: %v\n", err)
		os.Exit(1)
	}

	srv, err := irecallmcp.NewServer(cfg, binaryVersion())
	if err != nil {
		fmt.Fprintf(os.Stderr, "irecall-mcp: %v\n", err)
		os.Exit(1)
	}

	if err := mcpserver.ServeStdio(srv); err != nil {
		fmt.Fprintf(os.Stderr, "irecall-mcp: %v\n", err)
		os.Exit(1)
	}
}

func usageText(fs *flag.FlagSet, program string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Usage: %s [flags]\n\n", program)
	buf.WriteString("irecall-mcp exposes the local iRecall web API as MCP tools over stdio.\n\n")
	buf.WriteString("Environment:\n")
	buf.WriteString("  IRECALL_BASE_URL   Base URL for the iRecall web server (default: http://127.0.0.1:9527)\n")
	buf.WriteString("  IRECALL_API_TOKEN  Bearer token used for authenticated REST requests\n\n")
	buf.WriteString("Flags:\n")
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&buf, "  -%s\n    \t%s\n", f.Name, f.Usage)
	})
	buf.WriteString("\nExamples:\n")
	fmt.Fprintf(&buf, "  IRECALL_API_TOKEN=... %s\n", program)
	fmt.Fprintf(&buf, "  IRECALL_BASE_URL=http://127.0.0.1:9527 IRECALL_API_TOKEN=... %s\n", program)
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
