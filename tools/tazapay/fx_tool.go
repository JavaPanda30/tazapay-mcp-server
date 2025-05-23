package tazapay

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// FXTool defines the tool structure
type FXTool struct {
	logger *slog.Logger
}

// NewFXTool returns a new instance of the FXTool
func NewFXTool(logger *slog.Logger) *FXTool {
	logger.Info("Initializing FXTool")

	return &FXTool{
		logger: logger,
	}
}

// Definition registers this tool with the MCP platform
func (*FXTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.FXToolName,
		mcp.WithDescription(constants.FXToolDescription),
		mcp.WithString(constants.FXFromField, mcp.Required(), mcp.Description(constants.FXFromDescription)),
		mcp.WithString(constants.FXToField, mcp.Required(), mcp.Description(constants.FXToDescription)),
		mcp.WithNumber(constants.FXAmountField, mcp.Required(), mcp.Description(constants.FXAmountDescription)),
	)
}

// Handle processes the tool request and returns a result
func (t *FXTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.InfoContext(ctx, "Handling FXTool request", slog.Any("params", req.Params.Arguments))

	args := req.Params.Arguments

	// validate and extract arguments
	params, err := validateAndExtractFXArgs(t, args)
	if err != nil {
		t.logger.Error("Argument validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	// construct URL for API call
	url := fmt.Sprintf("%s?initial_currency=%s&final_currency=%s&amount=%d",
		constants.PaymentFxBaseURLProd, params.From, params.To, int(params.Amount))

	t.logger.InfoContext(ctx, "Calling FX API", slog.String("url", url))

	// call FX API
	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.Error("FX API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("HandleGETHttpRequest failed: %w", err)
	}

	// extract required fields from response
	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.Error("No 'data' in FX API response")
		return nil, constants.ErrNoDataInResponse
	}

	exRate, ok1 := data["exchange_rate"].(float64)
	if !ok1 {
		t.logger.Error("Invalid type for exchange_rate")
		return nil, utils.WrapFieldTypeError(t.logger, "exchange_rate")
	}

	converted, ok2 := data["converted_amount"].(float64)
	if !ok2 {
		t.logger.Error("Invalid type for converted_amount")
		return nil, utils.WrapFieldTypeError(t.logger, "converted_amount")
	}

	result := fmt.Sprintf("Rate: %.2f, Converted Amount: %.2f", exRate, converted)
	t.logger.InfoContext(ctx, "FXTool result ready", slog.String("result", result))

	// return result
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

// validateAndExtractFXArgs validates request arguments and returns structured parameters
func validateAndExtractFXArgs(t *FXTool, args map[string]any) (types.FXParams, error) {
	var p types.FXParams
	var ok bool

	if p.Amount, ok = args[constants.FXAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.FXAmountField)
	}

	if p.From, ok = args[constants.FXFromField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.FXFromField)
	}

	if p.To, ok = args[constants.FXToField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.FXToField)
	}

	return p, nil
}
