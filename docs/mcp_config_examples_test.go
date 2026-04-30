package docs_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenClawMCPDocsAvoidShellTokenWrappers(t *testing.T) {
	paths := []string{
		filepath.FromSlash("MCP_OPENCLAW.md"),
		filepath.FromSlash("../packaging/systemd/README.md"),
	}

	for _, rel := range paths {
		path := filepath.Join(".", rel)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}
		text := string(data)
		for _, bad := range []string{
			`"command": "bash"`,
			`"command": "/bin/sh"`,
			"bash -lc IRECALL_API_TOKEN",
			"sh -c IRECALL_API_TOKEN",
			"IRECALL_API_TOKEN=$(cat",
			"IRECALL_API_TOKEN=`cat",
		} {
			if strings.Contains(text, bad) {
				t.Fatalf("%s contains unsafe shell-token wrapper %q", rel, bad)
			}
		}
	}
}

func TestPackagedOpenClawMCPConfigUsesDirectCommandAndTokenFile(t *testing.T) {
	path := filepath.Join("..", "packaging", "openclaw", "irecall-mcp-server.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	var cfg struct {
		MCP struct {
			Servers map[string]struct {
				Command string            `json:"command"`
				Args    []string          `json:"args"`
				Env     map[string]string `json:"env"`
			} `json:"servers"`
		} `json:"mcp"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal packaged config error = %v", err)
	}

	server, ok := cfg.MCP.Servers["irecall"]
	if !ok {
		t.Fatal("packaged config missing mcp.servers.irecall")
	}
	if strings.Contains(server.Command, "sh") || strings.Contains(server.Command, "bash") {
		t.Fatalf("Command = %q, want direct irecall-mcp binary", server.Command)
	}
	if got, want := server.Command, "/usr/local/bin/irecall-mcp"; got != want {
		t.Fatalf("Command = %q, want %q", got, want)
	}
	if len(server.Args) != 2 || server.Args[0] != "--token-file" || server.Args[1] != "/etc/irecall/api-token" {
		t.Fatalf("Args = %#v, want --token-file /etc/irecall/api-token", server.Args)
	}
	if got, want := server.Env["IRECALL_BASE_URL"], "http://127.0.0.1:9527"; got != want {
		t.Fatalf("IRECALL_BASE_URL = %q, want %q", got, want)
	}
}
