package tools

import (
	"github.com/gigol/irecall/mcp/irecallapi"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func RegisterHistoryTools(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	// Placeholder for the next MCP pass. History endpoints already exist in the
	// REST API, but the first shell only wires the highest-value bridge tools.
}
