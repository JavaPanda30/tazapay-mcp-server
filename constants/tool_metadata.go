package constants

// Payment Link Tool constants
const (
	PaymentLinkToolName = "generate_payment_link_tool"
	PaymentLinkToolDesc = "Generates a payment link with specified invoice details and customer information"

	InvoiceCurrencyField = "invoice_currency"
	InvoiceCurrencyDesc  = "Currency in which the invoice is to be raised (e.g., USD, EUR)"

	PaymentAmountField = "payment_amount"
	PaymentAmountDesc  = "Total invoice amount to be paid"

	CustomerNameField = "customer_name"
	CustomerNameDesc  = "Full name of the customer"

	CustomerEmailField = "customer_email"
	CustomerEmailDesc  = "Email address of the customer"

	CustomerCountryField = "customer_country"
	CustomerCountryDesc  = "Country of the customer"

	TransactionDescField = "transaction_description"
	TransactionDesc      = "Short description or purpose of the transaction"
)

// FX Tool constants
const (
	FXToolName        = "tazapay_fetch_fx_tool"
	FXToolDescription = "Get FX rate from one currency to another using Tazapay FX rate"

	FXFromField       = "from"
	FXFromDescription = "Currency to convert from. It should be in 3 letter currency code. Example : USD, INR"

	FXToField       = "to"
	FXToDescription = "Currency to convert to. It should be in 3 letter currency code. Example : USD, INR"

	FXAmountField       = "amount"
	FXAmountDescription = "Amount to convert. It should be a number and should not have any decimal places."
)

// Balance Fetch tool
const (
	BalanceToolName = "tazapay_fetch_balance_tool"
	BalanceToolDesc = "Get balance from Tazapay. Send currency code to fetch balance for that currency." +
		" For all the balances available in Tazapay send empty string."

	BalanceCurrencyField = "currency"
	BalanceCurrencyDesc  = "Currency to fetch balance for. It should be in 3 letter currency code. Example : USD, INR"
)

// create beneficiary tool
const (
	CreateBeneficiaryToolName = "create_beneficiary_tool"
	CreateBeneficiaryToolDesc = "Create a new beneficiary in Tazapay."
)
