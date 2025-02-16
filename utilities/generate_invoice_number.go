package utilities

import (
	"fmt"
	database "schoolPayment/configs"
	"schoolPayment/constants"
	models "schoolPayment/models"

	"schoolPayment/repositories"
	"strconv"
	"time"
)

func GenerateInvoiceNumber(schoolId uint) (string, error) {
	// Get invoice format for the school
	invoiceFormat, err := repositories.NewInvoiceFormatRepository(database.DB).GetBySchoolID(schoolId)
	if err != nil {
		return "", fmt.Errorf("failed to get invoice format: %v", err)
	}

	if invoiceFormat == nil {
		return "", fmt.Errorf("invoice format not found for school ID: %d", schoolId)
	}

	prefix := invoiceFormat.Prefix
	format := invoiceFormat.Format
	currentDate := time.Now()

	// Prepare the date part based on format
	var datePart string
	var whereClause string

	switch format {
	case "DDMMYY00001":
		datePart = currentDate.Format("020106")
		whereClause = fmt.Sprintf("DATE(transaction_billings.created_at) = DATE('%s')", currentDate.Format(constants.DateFormatYYYYMMDD))
	case "MMYY00001":
		datePart = currentDate.Format("0106")
		whereClause = fmt.Sprintf("DATE_TRUNC('month', transaction_billings.created_at) = DATE_TRUNC('month', '%s'::timestamp)",
			currentDate.Format(constants.DateFormatYYYYMMDD))
	case "YY00001":
		datePart = currentDate.Format("06")
		whereClause = fmt.Sprintf("DATE_TRUNC('year', transaction_billings.created_at) = DATE_TRUNC('year', '%s'::timestamp)",
			currentDate.Format(constants.DateFormatYYYYMMDD))
	case "00001", "1":
		datePart = ""
		whereClause = "1=1"
	default:
		return "", fmt.Errorf("unsupported invoice format: %s", format)
	}

	// Get last invoice number with the same format
	var lastInvoice models.TransactionBilling
	query := database.DB.
		Joins("JOIN students ON students.id = transaction_billings.student_id").
		Joins("JOIN user_students ON user_students.student_id = students.id").
		Joins("JOIN user_schools ON user_schools.user_id = user_students.user_id").
		Where("invoice_number LIKE ? AND user_schools.school_id = ?", prefix+datePart+"%", schoolId)

	if format != "00001" && format != "1" {
		query = query.Where(whereClause)
	}
	err = query.Order("transaction_billings.created_at DESC").First(&lastInvoice).Error

	newSequenceNumber := 1
	if err == nil && lastInvoice.InvoiceNumber != "" {
		// Extract sequence number from last invoice
		seqStr := lastInvoice.InvoiceNumber[len(prefix)+len(datePart):]
		lastSeq, err := strconv.Atoi(seqStr)
		if err != nil {
			return "", fmt.Errorf("failed to parse last sequence number: %v", err)
		}
		newSequenceNumber = lastSeq + 1
	}

	// Format sequence number based on format
	var sequenceFormat string
	switch format {
	case "1":
		sequenceFormat = "%d"
	default:
		sequenceFormat = "%05d"
	}

	// Generate new invoice number
	newInvoiceNumber := fmt.Sprintf("%s%s"+sequenceFormat, prefix, datePart, newSequenceNumber)
	return newInvoiceNumber, nil
}
