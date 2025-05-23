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
		"create_beneficiary_tool",
		mcp.WithDescription("Create a beneficiary on Tazapay"),
		mcp.WithString("name", mcp.Required()),
		mcp.WithString("email"),
		mcp.WithString("type", mcp.Required(),mcp.Enum("individual", "business")),
		mcp.WithString("national_identification_number"),
		mcp.WithString("tax_id"),

		// Destination details object with correct required array at object level
		mcp.WithObject("destination_details", mcp.Required(),
			mcp.Properties(map[string]any{
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
			}),
			// Add required array for destination_details object
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
				"country":     map[string]any{"type": "string"},
			}),
		),

		// Documents object
		mcp.WithObject("documents",
			mcp.Properties(map[string]any{
				"type": map[string]any{"type": "string"},
				"url":  map[string]any{"type": "string"},
			}),
		),
	)
}

// Handle processes tool requests
func (t *CreateBeneficiaryTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	t.logger.InfoContext(ctx, "Handling CreateBeneficiaryTool request", "args", args)
	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	var payload types.DestinationRequest
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
