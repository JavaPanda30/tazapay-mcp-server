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

// FundPayoutTool funds a payout in requires_funding state

type FundPayoutTool struct {
	logger *slog.Logger
}

func NewFundPayoutTool(logger *slog.Logger) *FundPayoutTool {
	logger.Info("Initializing FundPayoutTool")
	return &FundPayoutTool{logger: logger}
}

func (t *FundPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"fund_payout_tool",
		mcp.WithDescription("Fund a payout in requires_funding state by ID on Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the payout in requires_funding state")),
	)
}

func (t *FundPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling FundPayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pot_", id) != nil {
		err := errors.New("missing or invalid payout id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/payout/%s/fund", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fund payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in fund payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	status, _ := data["status"].(string)
	resultText := "Payout funded. Status: " + status

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled FundPayoutTool request", "result", result)

	return result, nil
}
