package payout

import (
	"context"
	"encoding/json"
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
	return &GetPayoutTool{logger: logger}
}

func (*GetPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.GetPayoutToolName,
		mcp.WithDescription(constants.GetPayoutToolDesc),
		mcp.WithString(constants.GetPayoutIDField, mcp.Required(), mcp.Description(constants.GetPayoutIDDesc)),
	)
}

func (t *GetPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "Invalid arguments type for GetPayoutTool")
		return nil, fmt.Errorf("%w", constants.ErrInvalidArgumentsType)
	}

	t.logger.InfoContext(ctx, "Handling GetPayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pot_", id) != nil {
		t.logger.ErrorContext(ctx, constants.ErrMissingOrInvalidPayoutID.Error())
		return nil, constants.ErrMissingOrInvalidPayoutID
	}

	url := fmt.Sprintf("%s/payout/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal payout data", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: string(jsonBytes)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetPayoutTool request", "result", result)

	return result, nil
}
