package banks

import (
	"bd-backend/internal/models"
	"fmt"
	"strconv"
	"strings"
)

// scanBank formats a Bank model and converts it to a pointer of UserBankInfo.
func scanBank(bank *models.Bank) (*models.UserBankInfo, error) {
	if bank == nil {
		return nil, fmt.Errorf("bank is nil")
	}

	// Format fields
	bank.Routing = FormatRoutingNumber(bank.Routing)
	bank.Name = CapitalizeEachWord(bank.Name)
	bank.Address = CapitalizeEachWord(bank.Address)
	bank.Address = TrimWhitespace(bank.Address)
	bank.City = CapitalizeEachWord(bank.City)
	bank.City = TrimWhitespace(bank.City)
	bank.State = TrimWhitespace(bank.State)
	bank.Zip = TrimWhitespace(bank.Zip)

	userBankInfo := &models.UserBankInfo{
		BankID:  bank.BankID,
		Cert:    bank.Cert,
		Routing: bank.Routing,
		Name:    bank.Name,
		Address: bank.Address,
		City:    bank.City,
		State:   bank.State,
		Zip:     bank.Zip,
	}

	return userBankInfo, nil
}

// removeLeadingZeros removes leading zeros from a numeric string.
func removeLeadingZeros(s string) string {
	trimmed := strings.TrimLeft(s, "0")
	if len(trimmed) == 0 {
		return "000000000"
	}
	return trimmed
}

// FormatRoutingNumber pads the routing number to ensure it is exactly 9 digits long.
func FormatRoutingNumber(s *string) *string {
	if s == nil {
		return nil
	}
	formatted := padRoutingNumber(*s)
	return &formatted
}

// padRoutingNumber pads the routing number to ensure it is exactly 9 digits long.
func padRoutingNumber(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) >= 9 {
		return trimmed[:9]
	}
	return strings.Repeat("0", 9-len(trimmed)) + trimmed
}

// CapitalizeEachWord capitalizes the first letter of each word in a string.
func CapitalizeEachWord(s *string) *string {
	if s == nil {
		return nil
	}
	str := *s
	words := strings.Fields(str) // Split the string into words
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	capitalized := strings.Join(words, " ") // Join the words back into a single string
	return &capitalized
}

// TrimWhitespace trims leading and trailing whitespace from a string.
func TrimWhitespace(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

// convertToInt converts an interface{} to an int
func convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("invalid type for conversion to int: %T", v)
	}
}
