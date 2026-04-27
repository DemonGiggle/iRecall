package tools

import (
	"context"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type getRecallHistoryArgs struct {
	ID int64 `json:"id"`
}

type deleteRecallHistoryArgs struct {
	IDs []int64 `json:"ids"`
}

func RegisterHistoryTools(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	listTool := mcpproto.NewTool(
		"irecall_list_history",
		mcpproto.WithDescription("List saved iRecall recall-history summaries."),
	)
	srv.AddTool(listTool, func(ctx context.Context, request mcpproto.CallToolRequest) (*mcpproto.CallToolResult, error) {
		history, err := client.ListRecallHistory(ctx)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to list iRecall history.", err), nil
		}
		return jsonResult(history)
	})

	getTool := mcpproto.NewTool(
		"irecall_get_history",
		mcpproto.WithDescription("Get one saved iRecall recall-history entry by ID, including referenced quotes."),
		mcpproto.WithNumber("id",
			mcpproto.Required(),
			mcpproto.Description("Recall-history entry ID."),
		),
	)
	srv.AddTool(getTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args getRecallHistoryArgs) (*mcpproto.CallToolResult, error) {
		if args.ID <= 0 {
			return mcpproto.NewToolResultError("id must be positive"), nil
		}
		entry, err := client.GetRecallHistory(ctx, args.ID)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to get iRecall history entry.", err), nil
		}
		return jsonResult(entry)
	}))

	deleteTool := mcpproto.NewTool(
		"irecall_delete_history",
		mcpproto.WithDescription("Delete one or more saved iRecall recall-history entries by ID."),
		mcpproto.WithArray("ids",
			mcpproto.Required(),
			mcpproto.Description("Recall-history entry IDs to delete."),
			mcpproto.Items(map[string]any{"type": "number"}),
		),
	)
	srv.AddTool(deleteTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args deleteRecallHistoryArgs) (*mcpproto.CallToolResult, error) {
		if len(args.IDs) == 0 {
			return mcpproto.NewToolResultError("ids is required"), nil
		}
		for _, id := range args.IDs {
			if id <= 0 {
				return mcpproto.NewToolResultError("ids must contain only positive IDs"), nil
			}
		}
		result, err := client.DeleteRecallHistory(ctx, args.IDs)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to delete iRecall history entries.", err), nil
		}
		return jsonResult(result)
	}))
}
