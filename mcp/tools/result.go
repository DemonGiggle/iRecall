package tools

import (
	"encoding/json"

	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func jsonResult(value any) (*mcpproto.CallToolResult, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcpproto.NewToolResultText(string(data)), nil
}
