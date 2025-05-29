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

// GetPayoutTool fetches a payout by ID

type GetPayoutTool struct {
	logger *slog.Logger
}

func NewGetPayoutTool(logger *slog.Logger) *GetPayoutTool {
	logger.Info("Initializing GetPayoutTool")
	return &GetPayoutTool{logger: logger}
}

func (t *GetPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.GetPayoutToolName,
		mcp.WithDescription(constants.GetPayoutToolDesc),
		mcp.WithString(constants.GetPayoutIDField, mcp.Required(), mcp.Description(constants.GetPayoutIDDesc)),
	)
}

func (t *GetPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling GetPayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := errors.New("Missing or invalid payout id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/payout/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: fmt.Sprintf("Payout data: %+v", data)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetPayoutTool request", "result", result)

	return result, nil
}
