package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	database "schoolPayment/configs"
	"schoolPayment/constants"
	"schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/coreapi"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type TransactionService struct {
	transactionRepository      repositories.TransactionRepository
	userRepositories           repositories.UserRepository
	schoolRepositories         repositories.SchoolRepository
	billingHistoryService      BillingHistoryServiceInterface
	billingStudentRepositories repositories.BillingStudentRepository
	billingRepositories        repositories.BillingRepositoryInterface
	studentRepository          repositories.StudentRepositoryInteface
	
}

func NewTransactionService(
	transactionRepository repositories.TransactionRepository,
	userRepositories repositories.UserRepository,
	schoolRepositories repositories.SchoolRepository,
	billingHistoryService BillingHistoryServiceInterface,
	billingStudentRepositories repositories.BillingStudentRepository,
	billingRepositories repositories.BillingRepositoryInterface,
	studentRepository repositories.StudentRepositoryInteface,
) TransactionService {
	return TransactionService{
		transactionRepository:      transactionRepository,
		userRepositories:           userRepositories,
		schoolRepositories:         schoolRepositories,
		billingHistoryService:      billingHistoryService,
		billingStudentRepositories: billingStudentRepositories,
		billingRepositories:        billingRepositories,
		studentRepository:          studentRepository,
	}
}

func (transactionService *TransactionService) GetAllTransaction(page int, limit int, search string, userID int) (response.TransactionListResponse, error) {
	var resp response.TransactionListResponse
	resp.Limit = limit
	resp.Page = page

	getUser, err := transactionService.userRepositories.GetUserByID(uint(userID))
	if err != nil {
		resp.Data = []models.TransactionBilling{}
		return resp, nil
	}

	schoolID := 0
	if getUser.UserSchool != nil {
		schoolID = int(getUser.UserSchool.School.ID)
	}

	dataTransactions, err := transactionService.transactionRepository.GetAllTransaction(page, limit, search, schoolID)
	if err != nil {
		resp.Data = []models.TransactionBilling{}

		return resp, nil
	}
	if len(dataTransactions) > 0 {
		resp.Data = dataTransactions
	} else {
		resp.Data = []models.TransactionBilling{}
	}

	return resp, nil
}

func (transactionService *TransactionService) MidtransPayment(billingService *BillingService, studendID, paymentMethodId, userId int, billingStudentIds []string) (*response.MidtransResponse, error) {
	var rspPayment response.MidtransResponse

	listBillingId, err := transactionService.billingStudentRepositories.GetListBillingId(billingStudentIds)
	if err != nil {
		return nil, err
	}

	totalAmount, err := transactionService.billingStudentRepositories.GetTotalAmountByBillingStudentIds(billingStudentIds)
	if err != nil {
		return nil, err
	}

	listAccountNumber, err := repositories.GetListAccountNumber(billingStudentIds)
	if err != nil {
		return nil, err
	}

	user, err := transactionService.userRepositories.GetUserByID(uint(userId))
	if err != nil {
		return nil, err
	}

	invoiceNumber, err := utilities.GenerateInvoiceNumber(user.UserSchool.SchoolID)
	if err != nil {
		fmt.Errorf(constants.MessageErrorGenerateInvoiceNumber, err)
	}

	paymentMethod, err := repositories.GetPaymentMethodByID(paymentMethodId)
	if err != nil {
		return nil, err
	}

	orderID := fmt.Sprintf("SCPAY-%s-%d", paymentMethod.PaymentMethod, time.Now().Unix())

	resp, err := utilities.SendRequestPayment(
		orderID,
		studendID,
		totalAmount,
		paymentMethodId,
		billingStudentIds,
		invoiceNumber,
		listAccountNumber,
		listBillingId,
		paymentMethod.BankName,
		userId,
	)
	if err != nil {
		return nil, err
	}
	rspPayment.OrderID = orderID
	rspPayment.Token = &resp.Token
	rspPayment.RedirectURL = &resp.RedirectURL

	return &rspPayment, nil
}

func (transactionService *TransactionService) PaymentDonation(studentID, billingID, amount, paymentMethodID, userID int) (*response.MidtransResponse, error) {
	var rspPayment response.MidtransResponse

	billing, err := transactionService.billingRepositories.GetBillingByID(billingID)
	user, err := transactionService.userRepositories.GetUserByID(uint(userID))
	student, err := transactionService.studentRepository.GetStudentByID(uint(studentID), user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	belumBayarCode, _ := BillingStatusCode("Belum bayar")
	strBelumBayarCode := strconv.Itoa(belumBayarCode)

	billingDetail := models.BillingDetail{
		BillingID:         billing.ID,
		DetailBillingName: fmt.Sprintf("%s - %s", billing.BillingName, student.Nis),
		Amount:            int64(amount),
		DueDate:           &firstOfMonth,
	}

	_, err = repositories.CreateBillingDetail(&billingDetail, student.Nis, false)
	if err != nil {
		return nil, err
	}

	billingStudent := models.BillingStudent{
		BillingID:         billing.ID,
		StudentID:         uint(studentID),
		PaymentStatus:     strBelumBayarCode,
		DetailBillingName: fmt.Sprintf("%s - %s", billing.BillingName, student.Nis),
		Amount:            int64(amount),
		DueDate:           &firstOfMonth,
		BillingDetailID:   billingDetail.ID,
	}

	rspBillingStudent, err := repositories.CreateBillingStudent(&billingStudent, student.Nis, false)
	if err != nil {
		return nil, err
	}

	invoiceNumber, err := utilities.GenerateInvoiceNumber(user.UserSchool.SchoolID)
	if err != nil {
		fmt.Errorf(constants.MessageErrorGenerateInvoiceNumber, err)
	}

	paymentMethod, err := repositories.GetPaymentMethodByID(paymentMethodID)
	if err != nil {
		return nil, err
	}

	listAccountNumber, err := repositories.GetListAccountNumber([]string{strconv.Itoa(studentID)})
	if err != nil {
		return nil, err
	}

	orderID := fmt.Sprintf("SCPAY-%s-%d", paymentMethod.PaymentMethod, time.Now().Unix())

	resp, err := utilities.SendRequestPayment(
		orderID,
		studentID,
		amount,
		paymentMethodID,
		[]string{strconv.Itoa(int(rspBillingStudent.ID))},
		invoiceNumber,
		listAccountNumber,
		[]int{int(billing.ID)},
		paymentMethod.BankName,
		userID,
	)
	if err != nil {
		return nil, err
	}
	rspPayment.OrderID = orderID
	rspPayment.Token = &resp.Token
	rspPayment.RedirectURL = &resp.RedirectURL

	return &rspPayment, nil
}

func MidtransCheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	rsp, err := utilities.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (transactionService *TransactionService) CreateTransactionService(request request.CreateTransactionRequest, userId int, c *fiber.Ctx) (models.TransactionBilling, error) {
	// Get user information to get schoolId
	user, err := transactionService.userRepositories.GetUserByID(uint(userId))
	if err != nil {
		return models.TransactionBilling{}, err
	}

	if user.UserSchool == nil {
		return models.TransactionBilling{}, fmt.Errorf("user not associated with any school")
	}

	invoiceNumber, err := utilities.GenerateInvoiceNumber(user.UserSchool.SchoolID)
	if err != nil {
		return models.TransactionBilling{}, fmt.Errorf(constants.MessageErrorGenerateInvoiceNumber, err)
	}

	epoch := time.Now().Unix()
	referenceNumber := fmt.Sprintf("%s%d", "000000", epoch)
	expiryTime := time.Now().Format("2006-01-02T15:04:05")

	tx := database.DB.Begin()

	listBillingId, err := transactionService.billingStudentRepositories.GetListBillingId(request.BillingStudentIds)
	if err != nil {
		return models.TransactionBilling{}, err
	}

	totalAmount, err := transactionService.billingStudentRepositories.GetTotalAmountByBillingStudentIds(request.BillingStudentIds)
	if err != nil {
		return models.TransactionBilling{}, err
	}

	listAccountNumber, err := repositories.GetListAccountNumber(request.BillingStudentIds)
	if err != nil {
		return models.TransactionBilling{}, err
	}

	// Calculate total billing after discount
	totalBillingAfterDiscon := totalAmount
	if request.Discount != 0 {
		if request.DiscountType == "%" {
			diskon := totalAmount * request.Discount / 100
			totalBillingAfterDiscon = totalAmount - diskon
		} else {
			totalBillingAfterDiscon = totalAmount - request.Discount
		}
	}

	if totalBillingAfterDiscon > request.AmountToPay {
		return models.TransactionBilling{}, fmt.Errorf("ERROR: total billing does not match amount to pay")
	}

	// accountNumber := repositories.GetVirtualAccountNumberFromTransaction(request.BillingId)

	// Prepare transaction data
	billingStudentIdsStr := strings.Join(request.BillingStudentIds, ",")
	listAccountNumberStr := strings.Join(listAccountNumber, ",")
	listBillingIdStr := utilities.IntsToString(listBillingId)
	transaction := models.TransactionBilling{
		StudentID:         uint(request.StudentId),
		BillingID:         listBillingIdStr,
		TransactionType:   "PT01",
		TotalAmount:       request.AmountToPay,
		ReferenceNumber:   referenceNumber,
		TransactionStatus: "PS02",
		BillingStudentIds: billingStudentIdsStr,
		InvoiceNumber:     invoiceNumber,
		AccountNumber:     listAccountNumberStr,
		ExpiryTime:        expiryTime,
	}
	transaction.CreatedBy = userId
	// Create transaction in the repository
	dataTransaction, err := transactionService.transactionRepository.CreateTransactionRepository(tx, &transaction)
	if err != nil {
		return models.TransactionBilling{}, err
	}

	// Prepare transaction detail data
	transactionDetail := models.TransactionBillingDetail{
		TransactionBillingID: transaction.ID,
		Discount:             request.Discount,
		DiscountType:         request.DiscountType,
		ChangeAmount:         request.ChangeAmount,
		BankName:             nil,
		VirtualAccountNumber: nil,
		TransactionTime:      func(t time.Time) *time.Time { return &t }(time.Now()),
	}
	transactionDetail.CreatedBy = userId
	// Create transaction billing detail
	if dataTransaction != nil {
		_, err := repositories.CreateTransactionBillingDetail(tx, &transactionDetail)
		if err != nil {
			return models.TransactionBilling{}, err
		}
	}

	if dataTransaction != nil {
		history := models.TransactionBillingHistory{
			TransactionBillingId: dataTransaction.ID,
			// TransactionDate:      time.Now(),
			// TransactionAmount:    request.TotalAmount,
			ReferenceNumber: referenceNumber,
			// TransactionType:      "kasir",
			// Description:          request.Description,
			OrderID:           "",
			InvoiceNumber:     invoiceNumber,
			TransactionStatus: "PS02",
		}
		history.CreatedBy = userId
		_, err = transactionService.transactionRepository.CreateTransactionHistoryRepository(tx, &history)
		if err != nil {
			// tx.Rollback()
			return models.TransactionBilling{}, err
		}
	}

	billingStudentIds := request.BillingStudentIds
	var billingStudentIdsInt []int
	for _, idStr := range billingStudentIds {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			// Tangani error jika konversi gagal
			fmt.Println("Error konversi:", err)
			return models.TransactionBilling{}, err
		}
		billingStudentIdsInt = append(billingStudentIdsInt, idInt)
	}
	if dataTransaction != nil {
		err := repositories.UpdateStatusPayment(billingStudentIdsInt)
		if err != nil {
			return models.TransactionBilling{}, err
		}
	}

	// Commit the transaction
	tx.Commit()

	user, errUser := transactionService.userRepositories.GetUserByID(uint(userId))
	student, errStudent := transactionService.studentRepository.GetStudentByID(uint(request.StudentId), user)
	school, errSchool := transactionService.schoolRepositories.GetSchoolByID(user.UserSchool.SchoolID)
	emailParent, errEmailParent := repositories.GetEmailParentById(request.StudentId)
	if errUser != nil || errStudent != nil || errSchool != nil || errEmailParent != nil {
		return models.TransactionBilling{}, fmt.Errorf("failed to get user, student, school, or email parent")
	}

	formattedAmount := formatToIDR(int64(transaction.TotalAmount))
	//parts := strings.Split(transactionDetail.TransactionTime.String(), " m=")
	//cleanedInput := parts[0]
	// t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 -07", cleanedInput)
	// if err != nil {
	// 	fmt.Println("Error parsing time:", err)
	// 	// return
	// }

	// Format the time to "12 November 2024 13:27"
	formattedDate := time.Now().Format(constants.DateFormatDDMMMYYYhhmm)
	schoolLogo := utilities.ConvertPath(school.SchoolLogo)
	bodyEmail := BodyEmailTransaction{
		StudentName:   student.FullName,
		InvoiceNumber: invoiceNumber,
		PaymentDate:   formattedDate,
		TotalPayment:  formattedAmount,
		SchoolName:    user.UserSchool.School.SchoolName,
		Year:          time.Now().Year(),
		SchoolLogo:    schoolLogo,
	}

	filename, err := transactionService.billingHistoryService.GenerateInvoice(c, invoiceNumber, false, userId)
	if err != nil {
		return models.TransactionBilling{}, err
	}
	// send email to parent when transaction is success
	emailTemplate := utilities.GenerateEmailBodyTransactionSukses()
	_, err = SendEmailPaymentSuccess(emailParent.Email, "Bukti Pembayaran Transaksi", emailTemplate, bodyEmail, filename)
	if err != nil {
		return models.TransactionBilling{}, fmt.Errorf("failed to send payment success email: %v", err)
	}
	return *dataTransaction, nil
}

// Validate totalBilling and amountToPay
// func validateTotalBilling(request request.CreateTransactionRequest) error {
// 	totalBillingAfterDiscon := request.TotalBilling
// 	if request.Discount != 0 {
// 		if request.DiscountType == "%" {
// 			totalBillingAfterDiscon = request.TotalBilling * request.Discount / 100
// 		} else {
// 			totalBillingAfterDiscon = request.TotalBilling - request.Discount
// 		}
// 	}
// 	if totalBillingAfterDiscon != request.AmountToPay {
// 		return fmt.Errorf("total billing does not match amount to pay")
// 	}
// 	return nil
// }

func (transactionService *TransactionService) UpdateFromWebHook(payload request.WebhookPayload, c *fiber.Ctx) error {
	err := utilities.ValidateSignature(payload.OrderID, payload.StatusCode, payload.GrossAmount, os.Getenv("SERVER_KEY"), payload.SignatureKey)
	if err != nil {
		return err
	}

	err = utilities.ChangeStatusTransactionFromWebhook(payload)
	if err != nil {
		return err
	}

	midtransPayment, err := repositories.GetMidtransPaymentLogByOrderId(payload.OrderID)

	if err != nil {
		return err
	}

	// Parse the JSON into a map
	var data map[string]interface{}
	err = json.Unmarshal([]byte(midtransPayment.RedirectUrl), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}

	redirectURL, _ := data["redirect_url"].(string)

	bodyEmail, user, transactionID, err := transactionService.prepareEmailData(payload)
	if err != nil {
		return err
	}

	switch payload.TransactionStatus {
	case "pending":
		emailTemplate := utilities.GenerateEmailBodyTransactionWaiting()

		cleanedTotalPayment := strings.ReplaceAll(bodyEmail.TotalPayment, "Rp", "")
		cleanedTotalPayment = strings.ReplaceAll(cleanedTotalPayment, ".00", "")

		totalPaymentBigInt, success := new(big.Int).SetString(cleanedTotalPayment, 10)
		if !success {
			return fmt.Errorf(constants.MessageErrorConvertTotalPaymenr, cleanedTotalPayment)
		}

		bodyEmail.TotalPayment = utilities.RupiahFormat(totalPaymentBigInt)
		fmt.Println("cek ", bodyEmail.TotalPayment)

		message := fmt.Sprintf("Hai %s, yuk selesaikan pembayaran transaksi kamu sebelum %s. Terima kasih", user.Username, bodyEmail.ExpireDate)

		err = SendEmailPaymentWaiting(user.Email, "Segera Melakukan Pembayaran", emailTemplate, bodyEmail)
		if err != nil {
			return err
		}

		scheduleSvc := &scheduleService{}

		err = scheduleSvc.SendPaymentNotif(int(user.ID), bodyEmail.StudentNis, bodyEmail.StudentName, transactionID, message, "payment_waiting", redirectURL)

		if err != nil {
			return err
		}
	case "settlement":
		filename, err := transactionService.billingHistoryService.GenerateInvoice(c, bodyEmail.InvoiceNumber, false, int(user.ID))
		if err != nil {
			return err
		}

		cleanedTotalPayment := strings.ReplaceAll(bodyEmail.TotalPayment, "Rp", "")
		cleanedTotalPayment = strings.ReplaceAll(cleanedTotalPayment, ".00", "")

		totalPaymentBigInt, success := new(big.Int).SetString(cleanedTotalPayment, 10)
		if !success {
			return fmt.Errorf(constants.MessageErrorConvertTotalPaymenr, cleanedTotalPayment)
		}

		bodyEmail.TotalPayment = utilities.RupiahFormat(totalPaymentBigInt)

		emailTemplate := utilities.GenerateEmailBodyTransactionSukses()
		_, err = SendEmailPaymentSuccess(user.Email, "Bukti Pembayaran Transaksi", emailTemplate, bodyEmail, filename)
		if err != nil {
			return fmt.Errorf("failed to send payment success email: %v", err)
		}

		message := fmt.Sprintf("Hai %s, Pembayaran Anda Berhasil. Terima kasih telah melakukan Pembayaran", user.Username)
		scheduleSvc := &scheduleService{}

		err = scheduleSvc.SendPaymentNotif(int(user.ID), bodyEmail.StudentNis, bodyEmail.StudentName, transactionID, message, "payment_success", "")
		if err != nil {
			return err
		}
	case "expire":
		emailTemplate := utilities.GenerateEmailBodyTransactionFailedMidtrans()

		cleanedTotalPayment := strings.ReplaceAll(bodyEmail.TotalPayment, "Rp", "")
		cleanedTotalPayment = strings.ReplaceAll(cleanedTotalPayment, ".00", "")

		totalPaymentBigInt, success := new(big.Int).SetString(cleanedTotalPayment, 10)
		if !success {
			return fmt.Errorf(constants.MessageErrorConvertTotalPaymenr, cleanedTotalPayment)
		}

		bodyEmail.TotalPayment = utilities.RupiahFormat(totalPaymentBigInt)
		message := fmt.Sprintf("Hai %s, Waktu pembayaran telah habis, sehingga transaksi tidak dapat diproses. Untuk informasi lebih lanjut, cek detail tagihan Anda di aplikasi", user.Username)

		err = SendEmailPaymentWaiting(user.Email, "Pembayaran Anda Dibatalkan", emailTemplate, bodyEmail)
		if err != nil {
			return err
		}

		scheduleSvc := &scheduleService{}

		err := scheduleSvc.SendPaymentNotif(int(user.ID), bodyEmail.StudentNis, bodyEmail.StudentName, transactionID, message, "payment_failed", "")
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

type BodyEmailTransaction struct {
	StudentName    string
	StudentNis     string
	InvoiceNumber  string
	PaymentDate    string
	TotalPayment   string
	SchoolName     string
	Year           int
	SchoolLogo     string
	BankName       string
	VirtualAccount string
	ExpireDate     string
}

func formatToIDR(amount int64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%d", amount)
}

// prepareEmailData retrieves and formats data for email notifications
func (transactionService *TransactionService) prepareEmailData(payload request.WebhookPayload) (BodyEmailTransaction, models.User, string, error) {
	transactionBilling, err := repositories.GetBillingByOrderId(payload.OrderID)
	if err != nil {
		return BodyEmailTransaction{}, models.User{}, "", err
	}

	transactionID := strconv.Itoa(int(transactionBilling.ID))

	student, err := repositories.GetStudentByIDOnlyStudent(transactionBilling.StudentID)
	if err != nil {
		return BodyEmailTransaction{}, models.User{}, "", err
	}

	school, err := repositories.GetSchoolByStudentId(transactionBilling.StudentID)
	if err != nil {
		return BodyEmailTransaction{}, models.User{}, "", err
	}

	user, _ := repositories.GetEmailParentById(int(transactionBilling.StudentID))

	// Determine VA details based on the bank
	var vaNumber, bankName string
	if len(payload.VANumbers) > 0 {
		vaNumber = payload.VANumbers[0].VANumber
		bankName = payload.VANumbers[0].Bank
	} else if payload.PermataVaNumber != "" {
		vaNumber = payload.PermataVaNumber
		bankName = "Permata"
	} else if payload.PaymentType == "echannel" && payload.BillerCode == "70012" {
		vaNumber = payload.BillerCode + payload.BillKey
		bankName = "Mandiri"
	}

	transactionTime := time.Now().Format(constants.DateFormatDDMMMYYYhhmm)
	expireTimeParse, _ := time.Parse("2006-01-02 15:04:05", payload.ExpiryTime)
	expireTime := expireTimeParse.Format(constants.DateFormatDDMMMYYYhhmm)
	schoolLogo := utilities.ConvertPath(school.SchoolLogo)
	stringAmount := strconv.Itoa(transactionBilling.TotalAmount)
	bodyEmail := BodyEmailTransaction{
		StudentName:    student.FullName,
		StudentNis:     student.Nis,
		InvoiceNumber:  transactionBilling.InvoiceNumber,
		PaymentDate:    transactionTime,
		TotalPayment:   stringAmount,
		SchoolName:     school.SchoolName,
		Year:           time.Now().Year(),
		SchoolLogo:     schoolLogo,
		VirtualAccount: vaNumber,
		BankName:       strings.ToUpper(bankName),
		ExpireDate:     expireTime,
	}

	return bodyEmail, user, transactionID, nil
}

func SendEmailPaymentSuccess(email string, subject string, emailTemplate string, bodyEmail BodyEmailTransaction, filename string) (string, error) {
	// Use html/template to parse the email template
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse email template: %w", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, bodyEmail); err != nil {
		return "", fmt.Errorf("Failed to execute email template: %w", err)
	}
	emailBody := bodyBuffer.String()
	err = utilities.SendVerificationEmailWithAttachment(email, "", subject, emailBody, "./"+filename)
	if err != nil {
		return "", fmt.Errorf("Failed to send verification email.")
	}

	go func() {
		err := os.Remove(filename)
		if err != nil {
			fmt.Printf("Error deleting temporary file: %v\n", err)
		}
	}()

	return "", nil
}

func SendEmailPaymentWaiting(email string, subject string, emailTemplate string, bodyEmail BodyEmailTransaction) error {
	// Use html/template to parse the email template
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("Failed to parse email template: %w", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, bodyEmail); err != nil {
		return fmt.Errorf("Failed to execute email template: %w", err)
	}

	emailBody := bodyBuffer.String()

	err = utilities.SendEmail(email, subject, emailBody)
	if err != nil {
		return fmt.Errorf("Failed to send verification email.")
	}

	return nil
}
