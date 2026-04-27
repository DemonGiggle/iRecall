package tools

import (
	"context"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func RegisterHealthTool(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	tool := mcpproto.NewTool(
		"irecall_health",
		mcpproto.WithDescription("Check whether the iRecall web API is reachable and the bearer token is accepted."),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcpproto.CallToolRequest) (*mcpproto.CallToolResult, error) {
		state, err := client.BootstrapState(ctx)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to reach the iRecall web API.", err), nil
		}
		return jsonResult(struct {
			OK             bool   `json:"ok"`
			ProductName    string `json:"productName,omitempty"`
			ProfilePresent bool   `json:"profilePresent"`
		}{
			OK:             true,
			ProductName:    state.ProductName,
			ProfilePresent: state.Profile != nil,
		})
	})
}
