package utilities

import (
	"fmt"
	"strings"
)

// FormatNumber menambahkan titik untuk angka dengan ribuan separator
func FormatNumber(number int) string {
	return fmt.Sprintf("%d", number)
}

// FormatWithThousandsSeparator menambahkan titik untuk angka dengan format ribuan
func FormatWithThousandsSeparator(number int) string {
	str := fmt.Sprintf("%d", number)
	n := len(str)

	// Jika angka kurang dari 4 digit, langsung kembalikan
	if n <= 3 {
		return str
	}

	// Tambahkan titik setiap 3 digit dari belakang
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}
	return result.String()
}
