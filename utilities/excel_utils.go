package utilities

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type ExcelUtility struct {
	File *excelize.File
}

// NewExcelUtility creates a new Excel utility instance
func NewExcelUtility() *ExcelUtility {
	return &ExcelUtility{
		File: excelize.NewFile(),
	}
}

// WriteHeaders writes headers to the Excel file
func (e *ExcelUtility) WriteHeaders(sheetName string, headers []string, makeHeadersBold bool) error {
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		e.File.SetCellValue(sheetName, cell, header)

		if makeHeadersBold {
			style, err := e.File.NewStyle(&excelize.Style{
				Font: &excelize.Font{Bold: true},
			})
			if err != nil {
				return fmt.Errorf("failed to create header style: %v", err)
			}
			e.File.SetCellStyle(sheetName, cell, cell, style)
		}
	}
	return nil
}

// SetColumnWidths sets the width for specified columns
func (e *ExcelUtility) SetColumnWidths(sheetName string, columnWidths map[string]float64) error {
	for col, width := range columnWidths {
		err := e.File.SetColWidth(sheetName, col, col, width)
		if err != nil {
			return fmt.Errorf("failed to set column width: %v", err)
		}
	}
	return nil
}

// CreateRupiahStyle creates a style for Rupiah currency format
func (e *ExcelUtility) CreateRupiahStyle() (int, error) {
	formatString := "[$Rp-421]#,##0"
	style, err := e.File.NewStyle(&excelize.Style{
		CustomNumFmt: &formatString,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create Rupiah style: %v", err)
	}
	return style, nil
}
func (e *ExcelUtility) Write(buffer *bytes.Buffer) error {
	return e.File.Write(buffer)
}

// WriteToBuffer writes the Excel file to a buffer
func (e *ExcelUtility) WriteToBuffer() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	err := e.File.Write(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to write excel to buffer: %v", err)
	}
	return buffer, nil
}

// SetResponseHeaders sets the appropriate headers for Excel file download
func SetExcelResponseHeaders(c *fiber.Ctx, filename string) {
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
}

// Close closes the Excel file
func (e *ExcelUtility) Close() error {
	return e.File.Close()
}

// AddAutoFilter adds auto filter to specified range
func (e *ExcelUtility) AddAutoFilter(sheetName, startCell, endCell string) error {
	err := e.File.AutoFilter(sheetName, fmt.Sprintf("%s:%s", startCell, endCell), nil)
	if err != nil {
		return fmt.Errorf("failed to add auto filter: %v", err)
	}
	return nil
}

// SetCellValue sets the value for a specified cell
func (e *ExcelUtility) SetCellValue(sheet, cell string, value interface{}) error {
	return e.File.SetCellValue(sheet, cell, value)
}

// SetCellStyle sets the style for a specified cell or range
func (e *ExcelUtility) SetCellStyle(sheet, startCell, endCell string, styleID int) error {
	err := e.File.SetCellStyle(sheet, startCell, endCell, styleID)
	if err != nil {
		return fmt.Errorf("failed to set cell style: %v", err)
	}
	return nil
}

// CreateCenterStyle creates a center-aligned style for cells
func (e *ExcelUtility) CreateCenterStyle() (int, error) {
	style, err := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create center style: %v", err)
	}
	return style, nil
}

func (e *ExcelUtility) CreateBoldStyle() (int, error) {
	style, err := e.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	return style, err
}

// func AddDropdownList(file *excelize.File, sheetName, cellRange string, options []string) error {
// 	validation := excelize.NewDataValidation(true)
// 	validation.Sqref = cellRange
// 	validation.SetDropList(options)
// 	return file.AddDataValidation(sheetName, validation)
// }

func (e *ExcelUtility) AutoFitColumn(sheetName, column string) error {
	// Get all rows in the column
	rows, err := e.File.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows for sheet %s: %v", sheetName, err)
	}

	// Calculate maximum width
	maxWidth := 5.0 // Minimum width
	for _, row := range rows {
		if len(row) > 0 {
			for i, cell := range row {
				if string(rune('A'+i)) == column {
					cellWidth := float64(len(cell)) * 1.0 // Approximate width factor
					if cellWidth > maxWidth {
						maxWidth = cellWidth
					}
				}
			}
		}
	}

	// Set column width
	if err := e.File.SetColWidth(sheetName, column, column, maxWidth); err != nil {
		return fmt.Errorf("failed to set column width: %v", err)
	}
	return nil
}
