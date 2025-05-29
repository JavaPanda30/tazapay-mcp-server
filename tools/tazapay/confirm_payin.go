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

// ConfirmPayinTool represents the confirm payin tool
// Similar structure to CreatePayinTool

type ConfirmPayinTool struct {
	logger *slog.Logger
}

func NewConfirmPayinTool(logger *slog.Logger) *ConfirmPayinTool {
	logger.Info("Initializing ConfirmPayinTool")
	return &ConfirmPayinTool{logger: logger}
}

func (t *ConfirmPayinTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"confirm_payin_tool",
		mcp.WithDescription("Confirm a payin and create a payment attempt on Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the already created payin")),
		mcp.WithObject("customer_details", mcp.Required(),
			mcp.Properties(map[string]any{
				"customer": map[string]any{"type": "string"},
			}),
		),
		mcp.WithString("success_url", mcp.Required()),
		mcp.WithString("cancel_url", mcp.Required()),
		mcp.WithObject("shipping_details"),
		mcp.WithObject("billing_details"),
		mcp.WithObject("transaction_documents"),
		mcp.WithObject("metadata"),
		mcp.WithString("reference_id"),
		mcp.WithString("statement_descriptor"),
		mcp.WithObject("payment_method_details", mcp.Required()),
		mcp.WithString("session_id", mcp.Required()),
	)
}

func (t *ConfirmPayinTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling ConfirmPayinTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := errors.New("Missing or invalid payin id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	// Remove id from args for payload
	delete(args, "id")
	payload := args

	url := fmt.Sprintf("%s/payin/%s/confirm", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to confirm payin", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in confirm payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	status, _ := data["status"].(string)
	statusDesc, _ := data["status_description"].(string)
	resultText := fmt.Sprintf("Payin confirmed. Status: %s. %s", status, statusDesc)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled ConfirmPayinTool request", "result", result)

	return result, nil
}
