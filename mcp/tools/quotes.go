package tools

import (
	"context"
	"strings"

	"github.com/gigol/irecall/mcp/irecallapi"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type addQuoteArgs struct {
	Content string `json:"content"`
}

type saveRecallAsQuoteArgs struct {
	Question string   `json:"question"`
	Response string   `json:"response"`
	Keywords []string `json:"keywords"`
}

func RegisterQuoteTools(srv *mcpserver.MCPServer, client *irecallapi.Client) {
	listTool := mcpproto.NewTool(
		"irecall_list_quotes",
		mcpproto.WithDescription("List stored quotes from iRecall."),
	)
	srv.AddTool(listTool, func(ctx context.Context, request mcpproto.CallToolRequest) (*mcpproto.CallToolResult, error) {
		quotes, err := client.ListQuotes(ctx)
		if err != nil {
			return mcpproto.NewToolResultErrorFromErr("Failed to list quotes from iRecall.", err), nil
		}
		return jsonResult(quotes)
	})

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
