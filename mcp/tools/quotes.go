package tools

import (
	"context"
	"strings"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type listQuotesArgs struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type addQuoteArgs struct {
	Content string `json:"content"`
}

type saveRecallAsQuoteArgs struct {
	Question string   `json:"question"`
	Response string   `json:"response"`
	Keywords []string `json:"keywords"`
}

type updateQuoteArgs struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

type deleteQuotesArgs struct {
	IDs []int64 `json:"ids"`
}

func RegisterQuoteTools(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	listTool := mcpproto.NewTool(
		"irecall_list_quotes",
		mcpproto.WithDescription("List stored quotes from iRecall using limit/offset pagination."),
		mcpproto.WithNumber("limit",
			mcpproto.Description("Maximum quotes to return. Defaults to 20; maximum 100."),
		),
		mcpproto.WithNumber("offset",
			mcpproto.Description("Number of newest quotes to skip before returning results. Defaults to 0."),
		),
	)
	srv.AddTool(listTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args listQuotesArgs) (*mcpproto.CallToolResult, error) {
		limit := args.Limit
		if limit < 0 {
			return mcpproto.NewToolResultError("limit must be non-negative"), nil
		}
		if limit == 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		offset := args.Offset
		if offset < 0 {
			return mcpproto.NewToolResultError("offset must be non-negative"), nil
		}
		quotes, err := client.ListQuotes(ctx, limit, offset)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to list quotes from iRecall.", err), nil
		}
		return jsonResult(struct {
			Limit  int                `json:"limit"`
			Offset int                `json:"offset"`
			Quotes []irecallapi.Quote `json:"quotes"`
		}{Limit: limit, Offset: offset, Quotes: quotes})
	}))

	addTool := mcpproto.NewTool(
		"irecall_add_quote",
		mcpproto.WithDescription("Add a free-form quote or note to iRecall."),
		mcpproto.WithString("content",
			mcpproto.Required(),
			mcpproto.Description("The note or quote content to store."),
		),
	)
	srv.AddTool(addTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args addQuoteArgs) (*mcpproto.CallToolResult, error) {
		content := strings.TrimSpace(args.Content)
		if content == "" {
			return mcpproto.NewToolResultError("content is required"), nil
		}
		quote, err := client.AddQuote(ctx, content)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to add a quote to iRecall.", err), nil
		}
		return jsonResult(quote)
	}))

	updateTool := mcpproto.NewTool(
		"irecall_update_quote",
		mcpproto.WithDescription("Update an existing iRecall quote by ID."),
		mcpproto.WithNumber("id",
			mcpproto.Required(),
			mcpproto.Description("Quote ID to update."),
		),
		mcpproto.WithString("content",
			mcpproto.Required(),
			mcpproto.Description("Replacement quote content."),
		),
	)
	srv.AddTool(updateTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args updateQuoteArgs) (*mcpproto.CallToolResult, error) {
		if args.ID <= 0 {
			return mcpproto.NewToolResultError("id must be positive"), nil
		}
		content := strings.TrimSpace(args.Content)
		if content == "" {
			return mcpproto.NewToolResultError("content is required"), nil
		}
		quote, err := client.UpdateQuote(ctx, args.ID, content)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to update quote in iRecall.", err), nil
		}
		return jsonResult(quote)
	}))

	deleteTool := mcpproto.NewTool(
		"irecall_delete_quotes",
		mcpproto.WithDescription("Delete one or more iRecall quotes by ID."),
		mcpproto.WithArray("ids",
			mcpproto.Required(),
			mcpproto.Description("Quote IDs to delete."),
			mcpproto.Items(map[string]any{"type": "number"}),
		),
	)
	srv.AddTool(deleteTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args deleteQuotesArgs) (*mcpproto.CallToolResult, error) {
		if len(args.IDs) == 0 {
			return mcpproto.NewToolResultError("ids is required"), nil
		}
		for _, id := range args.IDs {
			if id <= 0 {
				return mcpproto.NewToolResultError("ids must contain only positive IDs"), nil
			}
		}
		result, err := client.DeleteQuotes(ctx, args.IDs)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to delete quotes from iRecall.", err), nil
		}
		return jsonResult(result)
	}))

	saveTool := mcpproto.NewTool(
		"irecall_save_recall_as_quote",
		mcpproto.WithDescription("Persist a recall question/response pair as a normal iRecall quote."),
		mcpproto.WithString("question",
			mcpproto.Required(),
			mcpproto.Description("The original recall question."),
		),
		mcpproto.WithString("response",
			mcpproto.Required(),
			mcpproto.Description("The grounded recall response."),
		),
		mcpproto.WithArray("keywords",
			mcpproto.Description("Optional recall keywords to persist with the saved quote."),
			mcpproto.Items(map[string]any{"type": "string"}),
		),
	)
	srv.AddTool(saveTool, mcpproto.NewTypedToolHandler(func(ctx context.Context, request mcpproto.CallToolRequest, args saveRecallAsQuoteArgs) (*mcpproto.CallToolResult, error) {
		question := strings.TrimSpace(args.Question)
		response := strings.TrimSpace(args.Response)
		if question == "" {
			return mcpproto.NewToolResultError("question is required"), nil
		}
		if response == "" {
			return mcpproto.NewToolResultError("response is required"), nil
		}
		quote, err := client.SaveRecallAsQuote(ctx, question, response, args.Keywords)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to save the recall result as a quote in iRecall.", err), nil
		}
		return jsonResult(quote)
	}))
}
