package types

// PayoutRequest represents the payload for creating a payout
// Reuses beneficiary types for nested fields
// Add/modify fields as needed for your payout API

type PayoutRequest struct {
	Purpose                  string                    `json:"purpose"`
	Amount                   float64                   `json:"amount"`
	Currency                 string                    `json:"currency"`
	ReferenceID              string                    `json:"reference_id"`
	Beneficiary              string                    `json:"beneficiary"`
	BeneficiaryDetails       DestinationDetails        `json:"beneficiary_details"`
	TransactionDesc          string                    `json:"transaction_desc,omitempty"`
	StatementDesc            string                    `json:"statement_desc,omitempty"`
	ChargeType               string                    `json:"charge_type,omitempty"`
	Type                     string                    `json:"type,omitempty"`
	HoldingCurrency          string                    `json:"holding_currency,omitempty"`
	OnBehalfOf               string                    `json:"on_behalf_of,omitempty"`
	Metadata                 string                    `json:"metadata,omitempty"`
	LogisticsTrackingDetails *LogisticsTrackingDetails `json:"logistics_tracking_details,omitempty"`
}

// LogisticsTrackingDetails represents logistics tracking info for a payout
// (Add/modify fields as needed)
type LogisticsTrackingDetails struct {
	LogisticsProviderName string `json:"logistics_provider_name,omitempty"`
	LogisticsProviderCode string `json:"logistics_provider_code,omitempty"`
	TrackingNumber        string `json:"tracking_number,omitempty"`
}
