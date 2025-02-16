package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	database "schoolPayment/configs"
	"schoolPayment/constants"
	"schoolPayment/dtos/request"
	models "schoolPayment/models"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type TransactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return TransactionRepository{}
}

func (transactionRepository *TransactionRepository) GetAllTransaction(page int, limit int, search string, schoolID int) ([]models.TransactionBilling, error) {
	var transactions []models.TransactionBilling
	query := database.DB.Where("tranasctions.deleted_at IS NULL")

	if search != "" {
		// query = query.Where("users.email LIKE ? or users.username Like ?", "%"+search+"%", "%"+search+"%")
	}

	if schoolID != 0 {
		query = query.Joins("JOIN billing_students ON billing_students.id = transactions.billing_student_id").
			Joins("JOIN user_students ON user_students.student_id = billing_students.student_id").
			Joins("JOIN user_schools ON user_schools.user_id = user_students.user_id").
			Joins("JOIN schools ON schools.id = user_schools.school_id AND schools.deleted_at is null").
			Where("schools.id = ?", schoolID)
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&transactions)
	return transactions, result.Error
}

func SaveLogPaymentMidtrans(orderId string, requestBody *snap.Request, responseString string) error {
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	responseCode := strconv.Itoa(fiber.StatusOK)

	rq := models.MidtransPaymentLog{
		OrderID:          orderId,
		StatusCode:       responseCode,
		Token:            "",
		RedirectUrl:      responseString,
		RequestBodyJson:  string(requestBodyJSON),
		ResponseBodyJson: responseString,
	}

	result := database.DB.Create(&rq)
	return result.Error
}

func SaveLogCheckPaymentMidtrans(orderID string, res *coreapi.TransactionStatusResponse) error {
	responseBodyJson, err := json.Marshal(res)
	if err != nil {
		return err
	}

	responseCode := strconv.Itoa(fiber.StatusOK)

	rq := models.MidtransCheckPaymentLog{
		OrderID:          orderID,
		StatusCode:       responseCode,
		ResponseBodyJson: string(responseBodyJson),
	}

	result := database.DB.Create(&rq)
	return result.Error
}

func CreateTransactionBilling(orderID string, studendID, billingAmount, paymentMethodId int, billingStudentIds []string, invoiceNumber string, listAccountNumber []string, listBillingId []int, bankName string, userId int) error {
	epochTime := time.Now().Unix()
	referenceNumber := fmt.Sprintf("000000%d", epochTime)
	billingStudentIdsStr := strings.Join(billingStudentIds, ",")
	listAccountNumberStr := strings.Join(listAccountNumber, ",")
	listBillingIdStr := IntsToString(listBillingId)
	expiryTime := time.Now().Format("2006-01-02T15:04:05")

	tx := database.DB.Begin()

	rq := models.TransactionBilling{
		BillingID:         listBillingIdStr,
		StudentID:         uint(studendID),
		TransactionType:   "PT02",
		TotalAmount:       billingAmount,
		ReferenceNumber:   referenceNumber,
		OrderID:           orderID,
		TransactionStatus: "PS01",
		Description:       "Payment for Billing",
		BillingStudentIds: billingStudentIdsStr,
		InvoiceNumber:     invoiceNumber,
		AccountNumber:     listAccountNumberStr,
		ExpiryTime:        expiryTime,
	}
	rq.CreatedBy = userId
	// Save the transaction to the database
	result := database.DB.Create(&rq)
	if result.Error != nil {
		return result.Error
	}

	// Prepare transaction detail data
	transactionDetail := models.TransactionBillingDetail{
		TransactionBillingID:  rq.ID,
		MasterPaymentMethodID: uint(paymentMethodId),
		BankName:              &bankName,
		TransactionTime:       func(t time.Time) *time.Time { return &t }(time.Now()),
	}
	transactionDetail.CreatedBy = userId
	_, errTransactionDetail := CreateTransactionBillingDetail(tx, &transactionDetail)
	if errTransactionDetail != nil {
		return errTransactionDetail
	}

	// Call the function to save the transaction history
	err := SaveTransactionBillingHistory(&rq)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTransactionBilling(payload request.WebhookPayload) error {
	var transaction models.TransactionBilling

	if err := database.DB.Where(constants.FilterOrderId, payload.OrderID).First(&transaction).Error; err != nil {
		return err
	}

	if len(payload.VANumbers) > 0 {
		transaction.VirtualAccountNumber = payload.VANumbers[0].VANumber
	} else if payload.PermataVaNumber != "" {
		transaction.VirtualAccountNumber = payload.PermataVaNumber
	} else if payload.PaymentType == "echannel" && payload.BillerCode == "70012" {
		transaction.VirtualAccountNumber = payload.BillerCode + payload.BillKey
	}

	switch payload.TransactionStatus {
	case "settlement", "capture":
		transaction.TransactionStatus = "PS02"
	case "expire", "failure", "cancel":
		transaction.TransactionStatus = "PS03"
	case "pending":
		transaction.TransactionStatus = "PS01"
	default:
		return fmt.Errorf("unknown status: %s", payload.TransactionStatus)
	}

	transaction.ExpiryTime = payload.ExpiryTime

	// Save the updated transaction
	if err := database.DB.Save(&transaction).Error; err != nil {
		return err
	}

	if payload.TransactionStatus == "settlement" {
		if err := UpdateBillingStudentPaymentStatus(transaction.BillingStudentIds, "2"); err != nil {
			return err
		}
	}

	// Save the transaction history after a successful update
	if err := SaveTransactionBillingHistory(&transaction); err != nil {
		return err
	}

	return nil
}

func SaveTransactionBillingHistory(transaction *models.TransactionBilling) error {
	rq := models.TransactionBillingHistory{
		TransactionBillingId: transaction.ID, // Assuming transaction.ID is set after saving
		OrderID:              transaction.OrderID,
		ReferenceNumber:      transaction.ReferenceNumber,
		TransactionStatus:    transaction.TransactionStatus,
		InvoiceNumber:        transaction.InvoiceNumber,
	}
	rq.CreatedBy = transaction.CreatedBy
	result := database.DB.Create(&rq)
	return result.Error
}

func (transactionRepository *TransactionRepository) CreateTransactionRepository(tx *gorm.DB, transaction *models.TransactionBilling) (*models.TransactionBilling, error) {
	// result := tx.Create(transaction)
	result := database.DB.Create(&transaction)
	return transaction, result.Error
}

func (transactionRepository *TransactionRepository) CreateTransactionHistoryRepository(tx *gorm.DB, transactionHistory *models.TransactionBillingHistory) (*models.TransactionBillingHistory, error) {
	result := tx.Create(transactionHistory)
	// result := database.DB.Create(&transactionHistory)
	return transactionHistory, result.Error
}

func UpdateStatusPayment(billingStudentId []int) error {
	if err := database.DB.Model(&models.BillingStudent{}).
	Where("id IN ?", billingStudentId).
	Update("payment_status", "2").Error; err != nil {
	return err
}

return nil
}

func GetEmailStudentID(stduentId int) (*models.Student, error) {
	var student models.Student

	result := database.DB.Where("id = ? AND deleted_at IS NULL AND is_block = FALSE", stduentId).First(&student)
	return &student, result.Error
}

func UpdateBillingStudentPaymentStatus(billingStudentIDs string, newStatus string) error {
	// Split the billingStudentIDs string by commas and trim spaces
	idStrings := strings.Split(billingStudentIDs, ",")
	for _, idStr := range idStrings {
		idStr = strings.TrimSpace(idStr)
		// Convert each ID to an integer
		billingStudentID, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		// Declare a new BillingStudent variable for each iteration
		var billingStudent models.BillingStudent

		// Retrieve the BillingStudent record by its ID only
		if err := database.DB.Where("id = ?", billingStudentID).First(&billingStudent).Error; err != nil {
			return err
		}

		// Update the PaymentStatus of the retrieved BillingStudent record
		billingStudent.PaymentStatus = newStatus
		if err := database.DB.Save(&billingStudent).Error; err != nil {
			return err
		}
	}

	return nil
}

func SaveLogCheckPaymentMidtransFromWebhook(orderID string, res *request.WebhookPayload) error {
	responseBodyJson, err := json.Marshal(res)
	if err != nil {
		return err
	}

	responseCode := strconv.Itoa(fiber.StatusOK)

	rq := models.MidtransCheckPaymentLog{
		OrderID:          orderID,
		StatusCode:       responseCode,
		ResponseBodyJson: string(responseBodyJson),
	}

	result := database.DB.Create(&rq)
	return result.Error
}

func CreateTransactionBillingDetail(tx *gorm.DB, billingTransactionDetail *models.TransactionBillingDetail) (*models.TransactionBillingDetail, error) {
	result := database.DB.Save(&billingTransactionDetail)
	return billingTransactionDetail, result.Error
}

func GetListAccountNumber(billingStudentIds []string) ([]string, error) {
	var listAccountNumber []string
	billingStudentIdin := formatBillingStudentIds(billingStudentIds)
	query := `select distinct account_number  from bank_accounts ba 
		join billings b on b.bank_account_id = ba.id 
		join billing_students bs on bs.billing_id = b.id`
	query += fmt.Sprintf(" where bs.id in %s", billingStudentIdin)
	result := database.DB.Raw(query).Scan(&listAccountNumber)
	if result.Error != nil {
		return []string{}, result.Error
	}
	return listAccountNumber, nil
}

func formatBillingStudentIds(ids []string) string {
	// Join the slice elements with commas and add parentheses around them
	return "(" + strings.Join(ids, ",") + ")"
}

func IntsToString(nums []int) string {
	// Convert each int to a string
	strNums := make([]string, len(nums))
	for i, num := range nums {
		strNums[i] = strconv.Itoa(num)
	}
	// Join the string representations with commas (or any separator you prefer)
	return strings.Join(strNums, ",")
}

func GetVirtualAccountNumberFromTransaction(billingId int) int {
	var bankAccount models.BankAccount
	query := `select ba.* from billings b 
			join bank_accounts ba on ba.id = b.bank_account_id`
	query += fmt.Sprintf(" where b.id = %d ", billingId)
	result := database.DB.Raw(query).First(&bankAccount)
	if result.Error != nil {
		return 0
	}

	accountNumber, err := strconv.Atoi(bankAccount.AccountNumber)
	if err != nil {
		return 0
	}
	return accountNumber
}

func GetBillingByOrderId(orderID string) (*models.TransactionBilling, error) {
	var transactionBiling *models.TransactionBilling

	result := database.DB.Where(constants.FilterOrderId, orderID).First(&transactionBiling)

	return transactionBiling, result.Error
}

func GetMidtransPaymentLogByOrderId(orderID string) (*models.MidtransPaymentLog, error) {
	var midtransPaymentLog *models.MidtransPaymentLog

	result := database.DB.Model(&midtransPaymentLog).Where(constants.FilterOrderId, orderID).First(&midtransPaymentLog)

	return midtransPaymentLog, result.Error
}
