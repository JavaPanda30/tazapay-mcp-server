package tazapay

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// FetchCheckoutTool fetches the details of a checkout session by ID

type FetchCheckoutTool struct {
	logger *slog.Logger
}

func NewFetchCheckoutTool(logger *slog.Logger) *FetchCheckoutTool {
	logger.Info("Initializing FetchCheckoutTool")
	return &FetchCheckoutTool{logger: logger}
}

func (t *FetchCheckoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"fetch_checkout_tool",
		mcp.WithDescription("Fetch the details of a checkout session by ID from Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the checkout session")),
	)
}

func (t *FetchCheckoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling FetchCheckoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := errors.New("missing or invalid checkout session id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/checkout/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch checkout session", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in fetch checkout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: fmt.Sprintf("Checkout session data: %+v", data)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled FetchCheckoutTool request", "result", result)

	return result, nil
}
