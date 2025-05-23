package tazapay

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreatePayoutTool represents the create payout tool
type CreatePayoutTool struct {
	logger *slog.Logger
}

// NewCreatePayoutTool returns a new instance of the CreatePayoutTool
func NewCreatePayoutTool(logger *slog.Logger) *CreatePayoutTool {
	logger.Info("Initializing Create_Payout_Tool")
	return &CreatePayoutTool{logger: logger}
}

// Definition registers this tool with the MCP
func (t *CreatePayoutTool) Definition() mcp.Tool {
	t.logger.Info("Registering CreatePayoutTool with MCP")
	return mcp.NewTool(
		"create_payout_tool",
		mcp.WithDescription("Create a payout on Tazapay"),
		mcp.WithString("purpose", mcp.Required()),
		mcp.WithNumber("amount", mcp.Required()),
		mcp.WithString("currency", mcp.Required()),
		mcp.WithString("reference_id", mcp.Required()),
		mcp.WithString("beneficiary", mcp.Required()),
		mcp.WithString("transaction_description"),
		mcp.WithString("statement_descriptor"),
		mcp.WithString("charge_type"),
		mcp.WithString("type", mcp.Enum("individual", "company")),
		mcp.WithString("holding_currency"),
		mcp.WithString("on_behalf_of"),
		mcp.WithString("metadata"),

		// beneficiary_details object
		mcp.WithObject("beneficiary_details", mcp.Required(),
			mcp.Properties(map[string]any{
				"name":                           map[string]any{"type": "string"},
				"email":                          map[string]any{"type": "string"},
				"type":                           map[string]any{"type": "string"},
				"account_id":                     map[string]any{"type": "string"},
				"national_identification_number": map[string]any{"type": "string"},
				"tax_id":                         map[string]any{"type": "string"},
				"destination_details": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{
							"type": "string",
							"enum": []string{"bank", "wallet", "local_payment_network"},
						},
						"bank": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"bank_codes": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"swift_code":  map[string]any{"type": "string"},
										"bic_code":    map[string]any{"type": "string"},
										"ifsc_code":   map[string]any{"type": "string"},
										"aba_code":    map[string]any{"type": "string"},
										"sort_code":   map[string]any{"type": "string"},
										"branch_code": map[string]any{"type": "string"},
										"bsb_code":    map[string]any{"type": "string"},
										"bank_code":   map[string]any{"type": "string"},
										"cnaps":       map[string]any{"type": "string"},
									},
								},
								"account_number": map[string]any{"type": "string"},
								"iban":           map[string]any{"type": "string"},
								"bank_name":      map[string]any{"type": "string"},
								"branch_name":    map[string]any{"type": "string"},
								"country":        map[string]any{"type": "string"},
								"currency":       map[string]any{"type": "string"},
								"purpose_code":   map[string]any{"type": "string"},
								"firc_required":  map[string]any{"type": "boolean"},
								"account_type":   map[string]any{"type": "string"},
							},
						},
						"wallet": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"deposit_address": map[string]any{"type": "string"},
								"type":            map[string]any{"type": "string"},
								"currency":        map[string]any{"type": "string"},
							},
						},
						"local_payment_network": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"type":             map[string]any{"type": "string"},
								"deposit_key_type": map[string]any{"type": "string"},
								"deposit_key":      map[string]any{"type": "string"},
							},
						},
					},
				},
				"phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"number":       map[string]any{"type": "string"},
						"calling_code": map[string]any{"type": "string"},
					},
				},
				"address": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"line1":       map[string]any{"type": "string"},
						"line2":       map[string]any{"type": "string"},
						"city":        map[string]any{"type": "string"},
						"state":       map[string]any{"type": "string"},
						"postal_code": map[string]any{"type": "string"},
						"country":     map[string]any{"type": "string"},
					},
				},
				// documents object (not array)
				"documents": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string"},
						"url":  map[string]any{"type": "string"},
					},
				},
			}),
		),

		// logistics_tracking_details object
		mcp.WithObject("logistics_tracking_details",
			mcp.Properties(map[string]any{
				"logistics_provider_name": map[string]any{"type": "string"},
				"logistics_provider_code": map[string]any{"type": "string"},
				"tracking_number":         map[string]any{"type": "string"},
			}),
		),
	)
}

// Handle processes tool requests
func (t *CreatePayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.InfoContext(ctx, "Handling CreatePayoutTool request", "request", req)
	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: "Payout created with ID: ABC123"},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled CreatePayoutTool request", "result", result)
	return result, nil
}
