package utilities

import (
	"fmt"
	"math/big"
	"strings"
)

// FormatCurrency formats an integer amount into a currency string (e.g., "Rp10.000")
func FormatCurrency(amount *big.Int) string {
	if amount == nil {
		return "Rp0"
	}

	// Konversi nilai big.Int ke string
	amountStr := amount.String()
	n := len(amountStr)
	if n <= 3 {
		return fmt.Sprintf("Rp%s", amountStr)
	}

	// Tambahkan pemisah ribuan (.)
	var result strings.Builder
	result.WriteString("Rp")
	for i, digit := range amountStr {
		result.WriteRune(digit)
		if (n-i-1)%3 == 0 && i != n-1 {
			result.WriteRune('.')
		}
	}
	return result.String()
}