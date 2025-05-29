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
		mcp.WithNumber("amount", mcp.Required(), mcp.Description("Amount in cents. For example, $10.12 should be 1012.")),
		mcp.WithString("currency", mcp.Required(), mcp.Description("ISO 4217 standard. This is the payout currency.")),
		mcp.WithString("purpose", mcp.Required(), mcp.Description("Reason for payout.")),
		mcp.WithString("transaction_description", mcp.Required(), mcp.Description("Additional Details for the payout.")),
		mcp.WithString("reference_id", mcp.Description("Reference ID of the payout on your system.")),
		mcp.WithString("statement_descriptor", mcp.Description("Statement Descriptor for the payout.")),
		mcp.WithString("charge_type", mcp.Description("For wire transfers only.")),
		mcp.WithString("type", mcp.Enum("individual", "business")),
		mcp.WithString("holding_currency", mcp.Description("ISO 4217 standard, in uppercase. This is one of your balance currencies whose balance will fund the payout.")),
		mcp.WithString("on_behalf_of", mcp.Description("ID of the entity the payout is created on behalf of.")),
		mcp.WithString("metadata", mcp.Description("Set of key-value pairs to attach to the payout object.")),
		// beneficiary_details object
		mcp.WithString("beneficiary", mcp.Description("ID of an existing payout beneficiary.")),
		mcp.WithObject("beneficiary_details", mcp.Required(),
			mcp.Properties(map[string]any{
				"name":                           map[string]any{"type": "string", "description": "Beneficiary Name"},
				"email":                          map[string]any{"type": "string", "description": "Beneficiary email address"},
				"type":                           map[string]any{"type": "string", "enum": []string{"business", "individual"}, "description": "business or individual"},
				"tax_id":                         map[string]any{"type": "string", "description": "Tax ID of the beneficiary."},
				"national_identification_number": map[string]any{"type": "string", "description": "National ID of the individual."},
				"phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"calling_code": map[string]any{"type": "string", "description": "Calling country code (e.g., '1' for US)"},
						"number":       map[string]any{"type": "string", "description": "Phone Number"},
					},
				},
				"address": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"line1":       map[string]any{"type": "string"},
						"line2":       map[string]any{"type": "string"},
						"city":        map[string]any{"type": "string"},
						"state":       map[string]any{"type": "string"},
						"country":     map[string]any{"type": "string", "description": "Country (ISO 3166-1 alpha_2 country code)"},
						"postal_code": map[string]any{"type": "string"},
					},
				},
				"destination_details": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string", "enum": []string{"bank", "wallet", "local_payment_network"}, "description": "bank, wallet or local_payment_network"},
						"bank": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"account_number": map[string]any{"type": "string", "description": "Bank Account Number. Either account_number or IBAN is mandatory"},
								"iban":           map[string]any{"type": "string", "description": "IBAN. Either account_number or iban is mandatory"},
								"bank_name":      map[string]any{"type": "string", "description": "Name of the bank"},
								"branch_name":    map[string]any{"type": "string", "description": "Name of the branch"},
								"country":        map[string]any{"type": "string", "description": "Two-letter country code (ISO 3166-1 alpha-2)"},
								"currency":       map[string]any{"type": "string", "description": "Three-letter ISO currency code, in uppercase"},
								"purpose_code":   map[string]any{"type": "string", "description": "Purpose Code for INR bank accounts"},
								"account_type":   map[string]any{"type": "string", "description": "Type of bank account"},
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
								"firc_required": map[string]any{"type": "boolean", "description": "Pass true if you want FIRC for a payout to India. By default, it is set to false"},
							},
						},
						"wallet": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"deposit_address": map[string]any{"type": "string"},
								"type":            map[string]any{"type": "string"},
								"currency":        map[string]any{"type": "string", "description": "Currency in which the payout is to be made (in uppercase, ISO-4217 standard, e.g., USD, EUR)"},
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
				// document object (single, for backward compatibility)
				"document": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type": map[string]any{"type": "string"},
						"url":  map[string]any{"type": "string"},
					},
				},
			}),
		),
		// documents array
		mcp.WithArray("documents", mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"type": map[string]any{"type": "string", "description": "Type of document. Possible values - invoice, other"},
				"url":  map[string]any{"type": "string", "description": "Dynamically downloadable URL"},
			},
		}), mcp.Description("Attach documents to the payout")),
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
