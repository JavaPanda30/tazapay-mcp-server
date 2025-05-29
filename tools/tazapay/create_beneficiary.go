package tazapay

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// CreateBeneficiaryTool represents the create beneficiary tool
type CreateBeneficiaryTool struct {
	logger *slog.Logger
}

// NewCreateBeneficiaryTool returns a new instance of the CreateBeneficiaryTool
func NewCreateBeneficiaryTool(logger *slog.Logger) *CreateBeneficiaryTool {
	logger.Info("Initializing Create_Beneficiary_Tool")
	return &CreateBeneficiaryTool{logger: logger}
}

// Definition : registers this tool with the MCP
func (t *CreateBeneficiaryTool) Definition() mcp.Tool {
	t.logger.Info("Registering CreateBeneficiaryTool with MCP")

	return mcp.NewTool(
		constants.CreateBeneficiaryToolName,
		mcp.WithDescription(constants.CreateBeneficiaryToolDesc),
		mcp.WithString(constants.BeneficiaryNameField, mcp.Required()),
		mcp.WithString(constants.BeneficiaryEmailField),
		mcp.WithString(constants.BeneficiaryTypeField, mcp.Required(), mcp.Enum("individual", "business")),
		mcp.WithString(constants.BeneficiaryNationalIdField),
		mcp.WithString(constants.BeneficiaryTaxIdField),
		mcp.WithObject(constants.BeneficiaryDestinationDetailsField, mcp.Required(),
			mcp.Properties(map[string]any{
				"type": map[string]any{
					"type": "string",
					"enum": []string{"bank", "wallet", "local_payment_network"},
				},
				"bank": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"account_number": map[string]any{"type": "string"},
						"iban":           map[string]any{"type": "string"},
						"bank_name":      map[string]any{"type": "string"},
						"branch_name":    map[string]any{"type": "string"},
						"country":        map[string]any{"type": "string", "description": "ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc."},
						"currency":       map[string]any{"type": "string", "description": "Currency in which the beneficiary will receive funds (in uppercase, ISO-4217 standard, e.g., USD, EUR)"},
						"purpose_code":   map[string]any{"type": "string"},
						"firc_required":  map[string]any{"type": "boolean"},
						"account_type":   map[string]any{"type": "string"},
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
					},
				},
				"wallet": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"deposit_address": map[string]any{"type": "string"},
						"type":            map[string]any{"type": "string"},
						"currency":        map[string]any{"type": "string", "description": "Currency in which the beneficiary will receive funds (in uppercase, ISO-4217 standard, e.g., USD, EUR)"},
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
			}),
			mcp.Required(),
		),

		// Phone object
		mcp.WithObject("phone",
			mcp.Properties(map[string]any{
				"number":       map[string]any{"type": "string"},
				"calling_code": map[string]any{"type": "string"},
			}),
		),

		// Address object
		mcp.WithObject("address",
			mcp.Properties(map[string]any{
				"line1":       map[string]any{"type": "string"},
				"line2":       map[string]any{"type": "string"},
				"city":        map[string]any{"type": "string"},
				"state":       map[string]any{"type": "string"},
				"postal_code": map[string]any{"type": "string"},
				"country":     map[string]any{"type": "string", "description": "ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc."},
			}),
		),

		// Documents object
		mcp.WithObject("document",
			mcp.Properties(map[string]any{
				"type": map[string]any{"type": "string"},
				"url":  map[string]any{"type": "string"},
			}),
		),
	)
}

// Handle processes tool requests
func (t *CreateBeneficiaryTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling CreateBeneficiaryTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	var payload types.CreateBeneficiaryRequest
	if err := utils.MapToStruct(args, &payload); err != nil {
		t.logger.ErrorContext(ctx, "Failed to map arguments to struct", "error", err)
		return nil, err
	}

	t.logger.InfoContext(ctx, "Mapped arguments to struct", "payload", payload)
	// Basic validation for required fields
	if payload.Name == "" || payload.Type == "" || payload.DestinationDetails.Type == "" {
		err := utils.WrapMissingFieldsError([]string{"name", "type", "account_id", "destination_details.type"})
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	// Validate currency and country in destination_details if present
	if dest, ok := args[constants.BeneficiaryDestinationDetailsField].(map[string]any); ok {
		if bank, ok := dest["bank"].(map[string]any); ok {
			if currency, ok := bank["currency"].(string); ok && currency != "" {
				if err := utils.ValidateCurrency(currency); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}
			if country, ok := bank["country"].(string); ok && country != "" {
				if err := utils.ValidateCountry(country); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}
		}
		if wallet, ok := dest["wallet"].(map[string]any); ok {
			if currency, ok := wallet["currency"].(string); ok && currency != "" {
				if err := utils.ValidateCurrency(currency); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}
		}
	}
	// Validate country in address if present
	if address, ok := args["address"].(map[string]any); ok {
		if country, ok := address["country"].(string); ok && country != "" {
			if err := utils.ValidateCountry(country); err != nil {
				t.logger.ErrorContext(ctx, err.Error())
				return nil, err
			}
		}
	}

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreateBeneficiaryAPIURL, payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to create beneficiary", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in create beneficiary API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// The beneficiary ID is in data["id"]
	beneficiaryID, ok := data["id"].(string)
	if !ok || beneficiaryID == "" {
		t.logger.ErrorContext(ctx, "No beneficiary ID in response", "data", data)
		return nil, constants.ErrNoBeneficiaryId
	}

	// Optionally, you can include the destination as well
	destinationID, _ := data["destination"].(string)

	resultText := "Beneficiary created with ID: " + beneficiaryID
	if destinationID != "" {
		resultText += ", destinationID: " + destinationID
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled CreateBeneficiaryTool request", "result", result)

	return result, nil
}
