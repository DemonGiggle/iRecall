package tools

import (
	"context"
	"strings"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type recallArgs struct {
	Question string `json:"question"`
}

func RegisterRecallTool(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	tool := mcpproto.NewTool(
		"irecall_recall",
		mcpproto.WithDescription("Run the iRecall recall flow against the local note corpus."),
		mcpproto.WithString("question",
			mcpproto.Required(),
			mcpproto.Description("The recall question to ask iRecall."),
		),
	)
	srv.AddTool(tool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args recallArgs) (*mcpproto.CallToolResult, error) {
		question := strings.TrimSpace(args.Question)
		if question == "" {
			return mcpproto.NewToolResultError("question is required"), nil
		}
		result, err := client.RunRecall(ctx, question)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to run recall in iRecall.", err), nil
		}
		return jsonResult(result)
	}))
}
