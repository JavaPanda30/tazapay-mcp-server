package types

// DestinationRequest represents the payload for registering a destination
type DestinationRequest struct {
	Name                         string             `json:"name"`
	Email                        string             `json:"email,omitempty"`
	Type                         string             `json:"type"`
	AccountID                    string             `json:"account_id"`
	NationalIdentificationNumber string             `json:"national_identification_number,omitempty"`
	TaxID                        string             `json:"tax_id,omitempty"`
	DestinationDetails           DestinationDetails `json:"destination_details"`
	Phone                        *Phone             `json:"phone,omitempty"`
	Address                      *Address           `json:"address,omitempty"`
	Documents                    *Documents         `json:"documents,omitempty"`
}

// DestinationDetails represents the details of a destination
type DestinationDetails struct {
	Type                string               `json:"type"`
	Bank                *Bank                `json:"bank,omitempty"`
	Wallet              *Wallet              `json:"wallet,omitempty"`
	LocalPaymentNetwork *LocalPaymentNetwork `json:"local_payment_network,omitempty"`
}

// Bank represents bank details
type Bank struct {
	BankCodes     BankCodes `json:"bank_codes,omitempty"`
	AccountNumber string    `json:"account_number,omitempty"`
	IBAN          string    `json:"iban,omitempty"`
	BankName      string    `json:"bank_name,omitempty"`
	BranchName    string    `json:"branch_name,omitempty"`
	Country       string    `json:"country,omitempty"`
	Currency      string    `json:"currency,omitempty"`
	PurposeCode   string    `json:"purpose_code,omitempty"`
	FIRCRequired  bool      `json:"firc_required,omitempty"`
	AccountType   string    `json:"account_type,omitempty"`
}

// BankCodes holds various bank identification codes
type BankCodes struct {
	SwiftCode  string `json:"swift_code,omitempty"`
	BICCode    string `json:"bic_code,omitempty"`
	IFSCCode   string `json:"ifsc_code,omitempty"`
	ABACode    string `json:"aba_code,omitempty"`
	SortCode   string `json:"sort_code,omitempty"`
	BranchCode string `json:"branch_code,omitempty"`
	BSBCode    string `json:"bsb_code,omitempty"`
	BankCode   string `json:"bank_code,omitempty"`
	CNAPS      string `json:"cnaps,omitempty"`
}

// Wallet represents wallet details
type Wallet struct {
	DepositAddress string `json:"deposit_address,omitempty"`
	Type           string `json:"type,omitempty"`
	Currency       string `json:"currency,omitempty"`
}

// LocalPaymentNetwork represents local payment network details
type LocalPaymentNetwork struct {
	Type           string `json:"type,omitempty"`
	DepositKeyType string `json:"deposit_key_type,omitempty"`
	DepositKey     string `json:"deposit_key,omitempty"`
}

// Phone represents phone details
type Phone struct {
	Number      string `json:"number,omitempty"`
	CallingCode string `json:"calling_code,omitempty"`
}

// Address represents address details
type Address struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Documents represents documents details
type Documents struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
