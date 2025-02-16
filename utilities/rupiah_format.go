package utilities

import (
	"math/big"
	"strings"
)

func RupiahFormat(amount *big.Int) string {
	// Convert the big.Int to a string
	amountStr := amount.String()

	// Insert commas as thousand separators
	n := len(amountStr)
	if n <= 3 {
		return "Rp " + amountStr
	}

	var result []string
	for i, digit := range amountStr {
		result = append(result, string(digit))
		if (n-1-i)%3 == 0 && i != n-1 {
			result = append(result, ".")
		}
	}

	return "Rp " + strings.Join(result, "")
}