package mcpbridge

import (
	"fmt"

	"github.com/gigol/irecall/mcp/irecallapi"
	"github.com/gigol/irecall/mcp/tools"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func NewServer(cfg Config, version string) (*mcpserver.MCPServer, error) {
	client, err := irecallapi.NewClient(cfg.APIConfig())
	if err != nil {
		return nil, err
	}

	srv := mcpserver.NewMCPServer(
		"iRecall MCP",
		version,
		mcpserver.WithInstructions(instructions(cfg)),
		mcpserver.WithToolCapabilities(false),
	)

	registerTools(srv, client)
	return srv, nil
}

func instructions(cfg Config) string {
	return fmt.Sprintf(
		"Use these tools to query or write data through the local iRecall web API. Base URL: %s. Prefer the iRecall tools only when note recall or persistence is relevant.",
		cfg.BaseURL,
	)
}

func registerTools(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	tools.RegisterHealthTool(srv, client)
	tools.RegisterRecallTool(srv, client)
	tools.RegisterQuoteTools(srv, client)
	tools.RegisterHistoryTools(srv, client)
}
