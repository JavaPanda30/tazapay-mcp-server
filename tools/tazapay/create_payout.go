package tazapay

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
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
		mcp.WithString("purpose", mcp.Required(), mcp.Description("Reason for payout. Should be in form of PYR0XX where XX is a number from 01 TO 28. also default is PYR001.")),
		mcp.WithString("transaction_description", mcp.Required(), mcp.Description("Additional Details for the payout.")),
		mcp.WithString("reference_id", mcp.Description("Reference ID of the payout on your system.")),
		mcp.WithString("statement_descriptor", mcp.Description("Statement Descriptor for the payout.")),
		mcp.WithString("charge_type", mcp.Description("For wire transfers only."), mcp.Enum("shared", "ours")),
		mcp.WithString("type", mcp.Enum("local", "swift", "wallet")),
		mcp.WithString("holding_currency", mcp.Description("ISO 4217 standard, in uppercase. This is one of your balance currencies whose balance will fund the payout.")),
		mcp.WithString("on_behalf_of", mcp.Description("ID of the entity the payout is created on behalf of.")),
		mcp.WithString("metadata", mcp.Description("Set of key-value pairs to attach to the payout object.")),
		mcp.WithString("beneficiary", mcp.Description("ID of an existing payout beneficiary.")),
		mcp.WithObject("beneficiary_details",
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
									"type":        "object",
									"description": "Various banking codes required for international transfers. All code types such as swift_code, bic_code, ifsc_code, aba_code, sort_code, branch_code, bsb_code, bank_code, cnaps are grouped under 'bank_codes'. For example: {\"swift_code\":\"HSBCHCMCHKH\",\"ifsc_code\":\"HDFC0001234\"}",
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
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling CreatePayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	// Validate all arguments
	if err := t.validatePayoutArgs(ctx, args); err != nil {
		t.logger.ErrorContext(ctx, "Validation failed", "error", err)
		return nil, err
	}

	var payload types.PayoutRequest
	if err := utils.MapToStruct(args, &payload); err != nil {
		t.logger.ErrorContext(ctx, "Failed to map arguments to struct", "error", err)
		return nil, err
	}

	// clean payload to send only one of beneficiary or beneficiary_details
	if payload.Beneficiary != "" && utils.ValidatePrefixID("bnf_", payload.Beneficiary) == nil {
		// beneficiary is present and valid, remove beneficiary_details from payload
		temp := payload // create a copy

		type payoutRequestMap map[string]any
		var payloadMap payoutRequestMap
		b, _ := json.Marshal(temp)
		_ = json.Unmarshal(b, &payloadMap)
		delete(payloadMap, "beneficiary_details")

		resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreatePayoutAPIURL, payloadMap, constants.PostHTTPMethod)
		if err != nil {
			t.logger.ErrorContext(ctx, "Failed to create payout", "error", err)
			return nil, err
		}

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.logger.ErrorContext(ctx, "No data in create payout API response", "resp", resp)
			return nil, constants.ErrNoDataInResponse
		}

		payoutID, ok := data["id"].(string)
		if !ok || payoutID == "" {
			t.logger.ErrorContext(ctx, "No payout ID in response", "data", data)
			return nil, constants.ErrNoBeneficiaryID
		}

		resultText := "Payout created with ID: " + payoutID

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: resultText},
			},
		}
		t.logger.InfoContext(ctx, "Successfully handled CreatePayoutTool request", "result", result)

		return result, nil
	} else {
		// beneficiary_details is present and valid, clear beneficiary
		payload.Beneficiary = ""

		resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreatePayoutAPIURL, payload, constants.PostHTTPMethod)
		if err != nil {
			t.logger.ErrorContext(ctx, "Failed to create payout", "error", err)
			return nil, err
		}

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.logger.ErrorContext(ctx, "No data in create payout API response", "resp", resp)
			return nil, constants.ErrNoDataInResponse
		}

		payoutID, ok := data["id"].(string)
		if !ok || payoutID == "" {
			t.logger.ErrorContext(ctx, "No payout ID in response", "data", data)
			return nil, constants.ErrNoBeneficiaryID
		}

		resultText := "Payout created with ID: " + payoutID

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{Type: "text", Text: resultText},
			},
		}
		t.logger.InfoContext(ctx, "Successfully handled CreatePayoutTool request", "result", result)

		return result, nil
	}
}

// validatePayoutArgs validates the arguments for creating a payout
func (t *CreatePayoutTool) validatePayoutArgs(ctx context.Context, args map[string]any) error {
	// Enforce: either beneficiary or beneficiary_details, but not both or neither
	_, hasBeneficiary := args["beneficiary"]
	_, hasBeneficiaryDetails := args["beneficiary_details"]

	if (hasBeneficiary && hasBeneficiaryDetails) || (!hasBeneficiary && !hasBeneficiaryDetails) {
		return errors.New("either 'beneficiary' or 'beneficiary_details' must be provided, but not both or neither")
	}

	// In Handle, before mapping to struct, move any top-level bank codes into bank_codes
	if bd, ok := args["beneficiary_details"].(map[string]any); ok {
		if dd, ok := bd["destination_details"].(map[string]any); ok {
			if bank, ok := dd["bank"].(map[string]any); ok {
				bankCodes := map[string]any{}

				fields := []string{"swift_code", "bic_code", "ifsc_code", "aba_code", "sort_code", "branch_code", "bsb_code", "bank_code", "cnaps"}
				for _, field := range fields {
					if val, exists := bank[field]; exists {
						bankCodes[field] = val

						delete(bank, field)
					}
				}

				if len(bankCodes) > 0 {
					bank["bank_codes"] = bankCodes
				}
			}
			// Validate currency and country in destination_details.bank if present
			if bank, ok := dd["bank"].(map[string]any); ok {
				if currency, ok := bank["currency"].(string); ok && currency != "" {
					if err := utils.ValidateCurrency(currency); err != nil {
						t.logger.ErrorContext(ctx, err.Error())
						return err
					}
				}

				if country, ok := bank["country"].(string); ok && country != "" {
					if err := utils.ValidateCountry(country); err != nil {
						t.logger.ErrorContext(ctx, err.Error())
						return err
					}
				}
			}
			// Validate currency in destination_details.wallet if present
			if wallet, ok := dd["wallet"].(map[string]any); ok {
				if currency, ok := wallet["currency"].(string); ok && currency != "" {
					if err := utils.ValidateCurrency(currency); err != nil {
						t.logger.ErrorContext(ctx, err.Error())
						return err
					}
				}
			}
		}
		// Validate country in address if present
		if address, ok := bd["address"].(map[string]any); ok {
			if country, ok := address["country"].(string); ok && country != "" {
				if err := utils.ValidateCountry(country); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return err
				}
			}
		}
	}

	// Basic validation for required fields (amount, currency, purpose, transaction_description)
	missingFields := []string{}
	if args["amount"] == nil {
		missingFields = append(missingFields, "amount")
	}

	if args["currency"] == nil || args["currency"] == "" {
		missingFields = append(missingFields, "currency")
	}

	if args["purpose"] == nil || args["purpose"] == "" {
		missingFields = append(missingFields, "purpose")
	}

	if args["transaction_description"] == nil || args["transaction_description"] == "" {
		missingFields = append(missingFields, "transaction_description")
	}

	if len(missingFields) > 0 {
		err := utils.WrapMissingFieldsError(missingFields)
		t.logger.ErrorContext(ctx, err.Error())

		return err
	}

	return nil
}
