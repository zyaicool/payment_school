package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
)

type BillingHistoryRepositoryInterface interface {
    GetAllBillingHistory(
        page int, limit int, search string, studentID int, roleID int, schoolYearId int,
        paymentTypeId int, schoolID int, paymentStatusCode string, sortBy string, 
        sortOrder string, userID int, userLoginID int,
    ) ([]response.DataListBillingHistory, int, int64, error)
    GetDetailBillingHistoryIDRepositories(transactionId int) (response.BillingStudentByStudentIDBillingID, error)
    GetInstallmentHistoryDetails(id int) (response.BillingStudentForHistory, error)
    GetDataForInvoice(invoiceNumber string, schoolID int) ([]response.RespDataInvoice, error) // Menambahkan metode ini
    TotalAmountBillingStudent(transactionId int) (int64, error)
}

type BillingHistoryRepository struct{}

// NewBillingHistoryRepository mengembalikan pointer ke BillingHistoryRepository
func NewBillingHistoryRepository() *BillingHistoryRepository {
	return &BillingHistoryRepository{}
}

func (billingHistoryRepository *BillingHistoryRepository) GetAllBillingHistory(page int, limit int, search string, studentId int, roleID int, schoolYearId int, paymentTypeId int, schoolID int, paymentStatusCode string, sortBy string, sortOrder string, userID int, userLoginID int) ([]response.DataListBillingHistory, int, int64, error) {
	var billingStudent []response.DataListBillingHistory
	var total int64
	totalPages := 0
	offset := (page - 1) * limit

	fmt.Printf("schoolID: %d\n", schoolID) 

	// Base query for getting data
	query := `
		SELECT 
			tb.id,
			tb.invoice_number,
			CONCAT(s.nis, ' - ', s.full_name) AS student_name,
			tb.created_at AS payment_date,
    		CASE 
       		 	WHEN u.role_id = 2 THEN '-'
        		ELSE COALESCE((SELECT u1.username FROM users u1 WHERE u1.id = tb.created_by), 'system')
    		END AS username,
			tb.total_amount,
			CASE 
				WHEN tb.transaction_status = 'PS01' THEN 'menunggu'
				WHEN tb.transaction_status = 'PS02' THEN 'lunas'
				WHEN tb.transaction_status = 'PS03' THEN 'gagal'
				ELSE tb.transaction_status  -- Default case, if no match
			end as transaction_status,
			CASE 
				WHEN tb.transaction_type = 'PT01' THEN 'Kasir'
				ELSE (
					SELECT concat(mp.payment_method, ' - ', mp.bank_name) 
					FROM master_payment_method mp
					WHERE mp.id = tbd.master_payment_method_id
					LIMIT 1
				)
			END AS payment_method,
			tb.order_id
		FROM 
			transaction_billings tb 
		JOIN 
			transaction_billing_details tbd ON tbd.transaction_billing_id = tb.id
		JOIN 
			students s ON s.id = tb.student_id
		JOIN 
			user_students ust ON ust.student_id = s.id
		JOIN 
			user_schools usc ON	usc.user_id = ust.user_id
		join 
			schools sch on sch.id = usc.school_id
		JOIN 
			users u ON u.id = usc.user_id
	`

	if roleID == 2 && schoolYearId != 0 {
		query += fmt.Sprintf(" JOIN billings b ON array_position(string_to_array(tb.billing_id, ','), CAST(b.id AS TEXT)) > 0 ")
	}

	query += fmt.Sprintf("Where tb.deleted_at is null ")

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(s.full_name) LIKE '%s' OR s.nis LIKE '%s') ",
			"%"+search+"%", "%"+search+"%")
	}

	if userID != 0 {
		query += fmt.Sprintf(" AND tb.created_by = %d", userID)
	} else {
		if roleID == 2 {
			query += fmt.Sprintf(" AND ust.user_id = %d", userLoginID)
		}
	}

	if studentId != 0 {
		query += fmt.Sprintf(" AND s.id = %d ", studentId)
	}

	if paymentStatusCode != "" {
		query += fmt.Sprintf(" AND tb.transaction_status = '%s' ", paymentStatusCode)
	}

	if roleID == 2 {
		// filter khusus user ortu
		if schoolYearId != 0 {
			query += fmt.Sprintf(" AND b.school_year_id = %d ", schoolYearId)
		}
		query += fmt.Sprintf(" AND tb.transaction_status != 'PS03' ")
	} else if roleID == 1 || roleID == 3 || roleID == 4 || roleID == 5 {

		if paymentTypeId != 0 {
			if paymentTypeId == 99 {
				query += fmt.Sprintf(" AND tb.transaction_type = 'PT01' ")
			} else {
				query += fmt.Sprintf(" AND tbd.master_payment_method_id = %d ", paymentTypeId)
			}
		}

		fmt.Println(schoolID)
		if schoolID != 0 {
			query += fmt.Sprintf(" AND usc.school_id = %d", schoolID)
		}
	}

	if limit != 0 {
		// Get the total count of records without pagination (Offset/Limit)
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM (%s) AS total_query
		`, query)

		if err := database.DB.Raw(countQuery).Scan(&total).Error; err != nil {
			return nil, 0, 0, err
		}

		// Calculate total pages
		totalPages = int((total + int64(limit) - 1) / int64(limit)) // Ensure to round up for total pages

	} else {
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM (%s) AS total_query
		`, query)

		if err := database.DB.Raw(countQuery).Scan(&total).Error; err != nil {
			return nil, 0, 0, err
		}

		// Set totalPages to 1 since all data will be included in one response
		totalPages = 1
	}

	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	} else {
		query += " ORDER BY CASE WHEN tb.updated_at IS NOT NULL THEN 0 ELSE 1 END, tb.updated_at DESC, tb.created_at DESC"
	}

	// Apply pagination only if limit is not 0
	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	// query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)

	// Execute the query to get the billing records
	result := database.DB.Raw(query).Debug().Scan(&billingStudent)

	// Return billingStudent, totalPages, total, and error
	return billingStudent, totalPages, total, result.Error
}

func (billingHistoryRepository *BillingHistoryRepository) GetDetailBillingHistoryIDRepositories(transactionId int) (response.BillingStudentByStudentIDBillingID, error) {
	var billingStudent response.BillingStudentByStudentIDBillingID
	query := `
		select 
		tb.id,
		CONCAT(s.nis, ' - ', s.full_name) AS student_name,
		CONCAT(sg.school_grade_name , ' - ', sc.school_class_name) AS school_class,
		tb.invoice_number,tbd.transaction_time as payment_date,
		CASE 
			WHEN tb.transaction_status = 'PS01' THEN 'menunggu'
			WHEN tb.transaction_status = 'PS02' THEN 'lunas'
			WHEN tb.transaction_status = 'PS03' THEN 'gagal'
			ELSE tb.transaction_status  -- Default case, if no match
		end as transaction_status,
		tbd.change_amount,
		tbd.discount_type,
		tbd.discount as discount_amount,
		tb.total_amount,
		tb.billing_student_ids,
		tbd.master_payment_method_id as payment_method_id,
		tb.transaction_type

		from transaction_billings tb 
		JOIN 
			transaction_billing_details tbd ON tbd.transaction_billing_id = tb.id
		JOIN 
			students s ON s.id = tb.student_id
		join 
			school_classes sc on sc.id = s.school_class_id 
		join
			school_grades sg on sg.id = s.school_grade_id 

	`

	query += fmt.Sprintf(" where tb.id = %d", transactionId)
	result := database.DB.Raw(query).Scan(&billingStudent)
	if result.Error != nil {
		return billingStudent, result.Error
	}

	return billingStudent, nil
}
func (billingHistoryRepository *BillingHistoryRepository) GetInstallmentHistoryDetails(id int) (response.BillingStudentForHistory, error) {
	var billingStudentCustom response.BillingStudentForHistory

	query := `
		SELECT b.amount, b.due_date, a.billing_type, case when a.is_donation = true then a.billing_name else b.detail_billing_name end
		FROM billings a
		JOIN billing_students b ON a.id = b.billing_id
		WHERE b.id = ?
	`
	err := database.DB.Raw(query, id).Scan(&billingStudentCustom).Error
	if err != nil {
		return billingStudentCustom, err
	}

	return billingStudentCustom, nil
}

func (billingHistoryRepository *BillingHistoryRepository) GetDataForInvoice(invoiceNumber string, schoolID int) ([]response.RespDataInvoice, error) {
	var dataInvoices []response.DataInvoice
	var res []response.RespDataInvoice

	// Main query to fetch invoice data
	query := `
	SELECT 
		tb.invoice_number, 
		tbd.transaction_time AS payment_date,
		NOW() AS print_date, 
		CASE 
			WHEN tb.transaction_type = 'PT01' THEN 'kasir'
			ELSE (
				SELECT concat(mp.payment_method, ' - ', mp.bank_name) 
				FROM master_payment_method mp
				WHERE mp.id = tbd.master_payment_method_id
				LIMIT 1
			)
		END AS transaction_type,
		s.nis, 
		s.full_name AS student_name,
		sc.school_class_name, 
		tb.billing_student_ids, 
		sc.school_id, 
		tb.total_amount, 
		0 AS sub_total,
		tbd.discount
	FROM transaction_billings tb 
	JOIN transaction_billing_details tbd 
		ON tbd.transaction_billing_id = tb.id
	JOIN students s 
		ON s.id = tb.student_id
	JOIN school_classes sc 
		ON sc.id = s.school_class_id
	WHERE tb.invoice_number = ? AND sc.school_id = ?`

	result := database.DB.Raw(query, invoiceNumber, schoolID).Scan(&dataInvoices)
	if result.Error != nil {
		return nil, result.Error
	}

	// Aggregate all billing_student_ids to fetch billing students in one query
	var billingStudentIDs []string
	for _, invoice := range dataInvoices {
		billingStudentIDs = append(billingStudentIDs, strings.Split(invoice.BillingStudentIds, ",")...)
	}

	// Remove duplicates and fetch billing students
	billingStudentIDs = uniqueStrings(billingStudentIDs)
	var billingStudents []models.BillingStudent
	if len(billingStudentIDs) > 0 {
		result = database.DB.Where("id IN ?", billingStudentIDs).Debug().Find(&billingStudents)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	fmt.Println("billing student", billingStudentIDs)

	// Group billing students by ID for easier mapping
	billingStudentMap := map[int]models.BillingStudent{}
	for _, student := range billingStudents {
		billingStudentMap[int(student.ID)] = student
	}

	// Map invoice data and calculate SubTotal
	for _, invoice := range dataInvoices {
		var studentDetails []models.BillingStudent
		for _, idStr := range strings.Split(invoice.BillingStudentIds, ",") {
			id, _ := strconv.Atoi(idStr) // Convert to int
			if student, exists := billingStudentMap[id]; exists {
				student.DetailBillingName = DeleteNISFromBillingName(student.DetailBillingName, invoice.Nis)
				studentDetails = append(studentDetails, student)
			}
		}

		subTotal := int64(0)
		for _, student := range studentDetails {
			subTotal += int64(student.Amount)
		}

		respDataInvoice := response.RespDataInvoice{
			InvoiceNumber:   invoice.InvoiceNumber,
			PaymentDate:     invoice.PaymentDate,
			PrintDate:       invoice.PrintDate,
			TransactionType: invoice.TransactionType,
			Nis:             invoice.Nis,
			StudentName:     invoice.StudentName,
			SchoolClassName: invoice.SchoolClassName,
			SchoolID:        invoice.SchoolID,
			SubTotal:        subTotal,
			Discount:        invoice.Discount,
			TotalAmount:     invoice.TotalAmount,
			BillingStudents: studentDetails,
		}
		res = append(res, respDataInvoice)
	}

	return res, nil
}

func DeleteNISFromBillingName(billingName string, nis string) string {
	parts := strings.Split(billingName, " - ")
	if len(parts) > 0 {
		// Extract the text before " - "
		string1 := parts[0]
		string2 := parts[1]

		if string2 == nis {
			return string1
		} else {
			return billingName
		}
	}
	return ""
}

// Helper function to remove duplicate strings
func uniqueStrings(input []string) []string {
	keys := make(map[string]bool)
	unique := []string{}
	for _, entry := range input {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			unique = append(unique, entry)
		}
	}
	return unique
}

func GetRedirectUrlsByOrderIds(orderIds []string) (map[string]response.RedirectUrlData, error) {
	redirectUrls := make(map[string]response.RedirectUrlData)

	// Struct to hold each row result with the raw JSON
	var results []struct {
		OrderId     string `json:"order_id"`
		RedirectUrl string `json:"redirect_url"`
	}

	// Query to fetch order_id and redirect_url JSON for the specified orderIds
	if err := database.DB.Table("midtrans_payment_logs").
		Select("order_id, redirect_url").
		Where("order_id IN ?", orderIds).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// Parse the JSON string in redirect_url field
	for _, result := range results {
		var redirectData struct {
			Token       string `json:"token"`
			RedirectUrl string `json:"redirect_url"`
		}
		// Parse the JSON from redirect_url
		if err := json.Unmarshal([]byte(result.RedirectUrl), &redirectData); err != nil {
			return nil, err
		}

		// Add to the map with the parsed data
		redirectUrls[result.OrderId] = response.RedirectUrlData{
			OrderId:     result.OrderId,
			Token:       redirectData.Token,
			RedirectUrl: redirectData.RedirectUrl,
		}
	}

	return redirectUrls, nil
}

func (billingHistoryRepository *BillingHistoryRepository) TotalAmountBillingStudent(transactionId int) (int64, error) {
	var totalAmount int64
	var billingStudentIDs string
	query := `Select tb.billing_student_ids from transaction_billings tb`
	query += fmt.Sprintf(" where tb.id = %d ", transactionId)
	result := database.DB.Raw(query).Scan(&billingStudentIDs)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to fetch total amount: %w", result.Error)
	}

	idStrings := strings.Split(billingStudentIDs, ",")
	for _, idStr := range idStrings {
		id, _ := strconv.Atoi(idStr)
		var amount int64
		query := `Select amount from billing_students bs`
		query += fmt.Sprintf(" where bs.id = %d ", id)
		result := database.DB.Raw(query).Scan(&amount)
		if result.Error != nil {
			return 0, fmt.Errorf("failed to fetch total amount: %w", result.Error)
		}

		totalAmount += amount
	}

	return totalAmount, nil
}
