package repositories

import (
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"time"

	database "schoolPayment/configs"
)

type ScheduleRepository struct{}

func NewScheduleRepository() ScheduleRepository {
	return ScheduleRepository{}
}

func GetStatusPayment() ([]models.TransactionBilling, error) {
	var transactions []models.TransactionBilling

	// Mendapatkan waktu 24 jam yang lalu dari saat ini
	cutoffTime := time.Now().Add(-24 * time.Hour)

	// Query untuk mendapatkan transaksi dengan kondisi di atas
	query := database.DB.Where("transaction_billings.deleted_at IS NULL AND transaction_billings.transaction_status = ? AND transaction_billings.created_at >= ?", "PS01", cutoffTime)

	result := query.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func GetCheckStatusPayment() ([]response.CheckPaymentStatusFailed, error) {
	var status []response.CheckPaymentStatusFailed

	// Mendapatkan waktu 24 jam yang lalu dari saat ini
	cutoffTime := time.Now().Add(-24 * time.Hour)

	query := `
		select tb.id, transaction_type, tb.virtual_account_number, total_amount, reference_number, description, 
		student_id, order_id, transaction_status, invoice_number, bank_name, tbd.master_payment_method_id, s.full_name as student_name, tb.created_by, s.nis
		from transaction_billings tb  
		join transaction_billing_details tbd on tb.id = tbd.transaction_billing_id
		join students s on  tb.student_id = s.id 
		WHERE tb.deleted_at IS NULL AND tb.transaction_status = ? AND tb.created_at <= ?
	`
	result := database.DB.Raw(query, "PS01", cutoffTime).Scan(&status)
	if result.Error != nil {
		return status, result.Error
	}

	return status, nil
}

func UpdateTransactionStatus(id uint) error {
	var transactionBilling models.TransactionBilling

	if err := database.DB.Where("id = ?", id).Find(&transactionBilling).Error; err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"transaction_status": "PS03",
		"expiry_time":        nil,
	}

	if transactionBilling.ExpiryTime != "" {
		updateData["expiry_time"] = transactionBilling.ExpiryTime
	}

	if err := database.DB.Model(&models.TransactionBilling{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func GetSchoolLogoSendEmailFailed(id uint) (response.SchoolLogoSendEmailFailedResponse, error) {
	var school response.SchoolLogoSendEmailFailedResponse
	query := `
			select school_name, school_logo from transaction_billings tb 
			join user_schools us on tb.created_by = us.user_id 
			join schools s on us.school_id = s.id
			where tb.id = ?
	`
	result := database.DB.Raw(query, id).Scan(&school)
	if result.Error != nil {
		return school, result.Error
	}

	return school, nil
}

func GetBillingStudentReminder() ([]response.BillingStudentReminderList, error) {
	var billingStudent []response.BillingStudentReminderList

	query := `
			select bs.id, bs.due_date, bs.amount, bs.student_id, s.full_name, s2.school_name, s2.school_logo, us.user_id from billing_students bs 
			inner join students s on s.id = bs.student_id 
			inner join user_students us on us.id = s.id 
			inner join user_schools us2 on us2.id = us.user_id 
			inner join schools s2 on s2.id = us2.school_id 
			where due_date::date in (
				CURRENT_DATE::date + INTERVAL '7 days',
				CURRENT_DATE::date + INTERVAL '3 days',
				CURRENT_DATE::date,
				CURRENT_DATE::date + INTERVAL '-1 days'
			)
	`

	result := database.DB.Raw(query).Scan(&billingStudent)

	if result.Error != nil {
		return billingStudent, result.Error
	}

	return billingStudent, nil
}
