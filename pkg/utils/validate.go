package utils

import (
	"errors"
	"strings"
)

// ValidateCurrency checks if the currency is 3 uppercase letters (ISO 4217)
func ValidateCurrency(currency string) error {
	if len(currency) != 3 || currency != strings.ToUpper(currency) {
		return errors.New("invalid currency format: must be 3 uppercase letters (e.g., 'USD')")
	}
	return nil
}

// ValidateCountry checks if the country is 2 uppercase letters (ISO 3166 alpha-2)
func ValidateCountry(country string) error {
	if len(country) != 2 || country != strings.ToUpper(country) {
		return errors.New("invalid country format: must be 2 uppercase letters (e.g., 'US')")
	}
	return nil
}

// ValidatePrefixId checks if the id starts with the given prefix
func ValidatePrefixId(prefix, id string) error {
	if !strings.HasPrefix(id, prefix) {
		return errors.New("invalid id format: must start with '" + prefix + "'")
	}
	return nil
}
