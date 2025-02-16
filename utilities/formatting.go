package utilities

import (
	"fmt"
	"math/big"
	models "schoolPayment/models"
	"strings"
	"time"
)

func FormatPhoneNumber(phoneNumber string) string {
	if strings.HasPrefix(phoneNumber, "62") {
		// Add '+' prefix if the number starts with '62'
		if strings.HasPrefix(phoneNumber, "+") {
			return phoneNumber[1:]
		}
		return "+" + phoneNumber
	} else {
		// Add '0' prefix if the number does not start with '62'
		return "0" + phoneNumber
	}
}

func ParseBirthDate(dateStr string) (time.Time, error) {
	// List of possible date formats
	// formats := []string{"01-02-06", "01/02/06", "01-02-2006", "01/02/2006"}

	// var birthDate time.Time
	// var err error

	// // Try parsing the date with each format
	// for _, format := range formats {
	//     birthDate, err = time.Parse(format, dateStr)
	//     if err == nil {
	//         // If the date has a two-digit year, adjust it to four digits
	//         if birthDate.Year() < 100 {
	//             birthDate = birthDate.AddDate(2000, 0, 0) // Assumes years in 2000s
	//         }
	//         return birthDate, nil
	//     }
	// }

	// return birthDate, err // Return the last error if all formats fail

	formats := []string{
		"02-01-06",   // dd-mm-yy
		"02/01/06",   // dd/mm/yy
		"02-01-2006", // dd-mm-yyyy
		"02/01/2006", // dd/mm/yyyy
	}

	var parsedDate time.Time
	var err error

	// Iterate through all formats and attempt to parse the date
	for _, format := range formats {
		parsedDate, err = time.Parse(format, dateStr)
		if err == nil {
			// If the parsing is successful, return the parsed date
			return parsedDate, nil
		}
	}

	// Return error if none of the formats matched
	return time.Time{}, fmt.Errorf("could not parse date: %s", dateStr)
}

func FormatBigInt(value string) *big.Int {
	num := new(big.Int)
	rsp, _ := num.SetString(value, 10)
	return rsp
}

// Add these helper functions after the imports
func GetMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MakeSchoolGradeMap(grades []models.SchoolGrade) map[string]uint {
	m := make(map[string]uint)
	for _, grade := range grades {
		m[grade.SchoolGradeName] = grade.ID
	}
	return m
}

func MakeSchoolClassMap(classes []models.SchoolClass) map[string]uint {
	m := make(map[string]uint)
	for _, class := range classes {
		m[class.SchoolClassName] = class.ID
	}
	return m
}

func MakeSchoolYearMap(years []models.SchoolYear) map[string]uint {
	m := make(map[string]uint)
	for _, year := range years {
		m[year.SchoolYearName] = year.ID
	}
	return m
}

func MakeExistingStudentMap(students []models.Student) map[string]models.Student {
	m := make(map[string]models.Student)
	for _, student := range students {
		m[student.Nis] = student
	}
	return m
}
