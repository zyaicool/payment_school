package utilities

import (
	"bytes"
	"fmt"
	"schoolPayment/constants"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// Helper function to get pointer to string
func ptr(s string) *string {
	return &s
}

// Add this to utilities/excel.go
type DropdownOption struct {
	Label string
	Value string
}

func GenerateFileExcel(c *fiber.Ctx, headers []string, fileName string, data []string) (*bytes.Buffer, error) {
	// Create a new Excel file
	file := excelize.NewFile()

	// Create a new sheet and add headers to the first row
	sheetName := "Sheet1"
	file.SetSheetName(file.GetSheetName(0), sheetName)
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1" // For example: A1, B1, C1...
		file.SetCellValue(sheetName, cell, header)
	}

	for i, value := range data {
		cell := string(rune('A'+i)) + "2" // For example: A2, B2, C2...
		file.SetCellValue(sheetName, cell, value)
	}

	// Write the Excel file to a buffer
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	// Set the response headers for file download
	c.Set(fiber.HeaderContentType, constants.ContentTypeExcel)
	c.Set(fiber.HeaderContentDisposition, fileName)

	// Send the file buffer as the response
	return buffer, nil
}

func GenerateFileExcelUser(c *fiber.Ctx, headers []string, fileName string, data [][]string, roleOptions []string) (*bytes.Buffer, error) {
	// Create a new Excel file
	file := excelize.NewFile()

	// Create a new sheet and set the name
	sheetName := "Sheet1"
	file.SetSheetName(file.GetSheetName(0), sheetName)

	// Add headers to the first row
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1" // Example: A1, B1, C1...
		file.SetCellValue(sheetName, cell, header)
	}

	// Add example data to the second row
	for i, rowData := range data {
		for j, value := range rowData {
			cell := string(rune('A'+j)) + strconv.Itoa(i+2) // Example: A2, B2, C2...
			file.SetCellValue(sheetName, cell, value)
		}
	}

	// Add dropdown for the "Role" column
	if len(roleOptions) > 0 {
		dropdownFormula := "\"" + strings.Join(roleOptions, ",") + "\""

		// Define messages as variables to take their address
		promptTitle := "Valid test"
		promptMessage := "Please select a valid role from the dropdown."

		validation := &excelize.DataValidation{
			Type:             "list",
			Formula1:         dropdownFormula,
			AllowBlank:       false,
			ShowDropDown:     false,
			ShowInputMessage: true,
			PromptTitle:      &promptTitle,
			Prompt:           &promptMessage,
			Sqref:            "A2:A1048576", // Apply to column A, starting from row 2
		}

		// Add the data validation to the sheet
		if err := file.AddDataValidation(sheetName, validation); err != nil {
			return nil, err
		}
	}

	// Write the Excel file to a buffer
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	// Set the response headers for file download
	c.Set(fiber.HeaderContentType, constants.ContentTypeExcel)
	c.Set(fiber.HeaderContentDisposition, fileName)

	return buffer, nil
}

func GenerateFileExcelStudent(c *fiber.Ctx, headers []string, filename string, exampleData [][]interface{}, dropdowns map[string][]DropdownOption, formats map[string]string) (*bytes.Buffer, error) {
	// Create a new Excel file
	file := excelize.NewFile()

	// Create a new sheet and set the name
	sheetName := "Sheet1"
	file.SetSheetName(file.GetSheetName(0), sheetName)

	// Create text style for phone number column
	textStyle, err := file.NewStyle(&excelize.Style{
		CustomNumFmt: ptr("@"), // @ is the format code for text
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create text style: %v", err)
	}

	// Add headers to the first row
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1" // Example: A1, B1, C1...
		file.SetCellValue(sheetName, cell, header)
	}

	// Add data to the sheet starting from row 2
	for i, rowData := range exampleData {
		for j, value := range rowData {
			cell := string(rune('A'+j)) + strconv.Itoa(i+2) // Example: A2, B2, C2...
			file.SetCellValue(sheetName, cell, value)

			// Apply text format to phone number column (column J)
			if j == 9 { // Phone number is in column J (0-based index 9)
				if err := file.SetCellStyle(sheetName, cell, cell, textStyle); err != nil {
					return nil, err
				}
			}
		}
	}

	// Inside GenerateFileExcelStudent function
	for column, options := range dropdowns {
		// Create a hidden sheet to store the mapping
		mappingSheetName := "Mapping_" + column
		file.NewSheet(mappingSheetName)

		// Add labels and values to the mapping sheet
		for i, opt := range options {
			labelCell := fmt.Sprintf("A%d", i+1)
			valueCell := fmt.Sprintf("B%d", i+1)
			file.SetCellValue(mappingSheetName, labelCell, opt.Label)
			file.SetCellValue(mappingSheetName, valueCell, opt.Value)
		}

		// Create named range for the labels
		labelRange := fmt.Sprintf("%s!$A$1:$A$%d", mappingSheetName, len(options))
		dropdownFormula := "=" + labelRange

		promptTitle := "Valid input"
		promptMessage := "Please select a valid option from the dropdown."

		validation := &excelize.DataValidation{
			Type:             "list",
			Formula1:         dropdownFormula,
			AllowBlank:       false,
			ShowDropDown:     false,
			ShowInputMessage: true,
			PromptTitle:      &promptTitle,
			Prompt:           &promptMessage,
			Sqref:            column + "2:" + column + "1048576",
		}

		if err := file.AddDataValidation(sheetName, validation); err != nil {
			return nil, err
		}

		// Hide the mapping sheet
		file.SetSheetVisible(mappingSheetName, false)
	}

	// Apply custom number format for date column
	if dateFormat, ok := formats["D"]; ok {
		dateStyle, err := file.NewStyle(&excelize.Style{
			CustomNumFmt: &dateFormat,
		})
		if err != nil {
			return nil, err
		}

		// Apply the date style to the example data cells in column D
		for i := range exampleData {
			cell := fmt.Sprintf("D%d", i+2) // +2 because data starts after header
			file.SetCellStyle(sheetName, cell, cell, dateStyle)
		}
	}

	// Write the Excel file to a buffer
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	// Set the response headers for file download
	c.Set(fiber.HeaderContentType, constants.ContentTypeExcel)
	c.Set(fiber.HeaderContentDisposition, filename)

	// Return the buffer containing the Excel file
	return buffer, nil
}

func MapDisplayLabelToValueStudent(field string, displayLabel string) string {
	// Define the mappings based on your dropdown configurations
	mappings := map[string]map[string]string{
		"gender": {
			"Laki-laki": "laki-laki",
			"Perempuan": "perempuan",
		},
		"religion": {
			"Islam":     "islam",
			"Kristen":   "kristen",
			"Katolik":   "katolik",
			"Hindu":     "hindu",
			"Budha":     "budha",
			"Khonghucu": "khonghucu",
			"Lainnya":   "other",
		},
		"status": {
			"Aktif":          "aktif",
			"Tamat":          "tamat",
			"Pindah Sekolah": "pindah_sekolah",
			"Dropout":        "dropout",
		},
	}

	if valueMap, exists := mappings[field]; exists {
		if value, ok := valueMap[displayLabel]; ok {
			return value
		}
	}
	return displayLabel // Return original value if no mapping found
}

func GetColumnName(col int) string {
	var result string
	for col >= 0 {
		result = string('A'+col%26) + result
		col = col/26 - 1
	}
	return result
}
