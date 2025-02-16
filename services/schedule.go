package services

import (
	"bytes"
	"fmt"
	"html/template"
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/repositories"
	"schoolPayment/utilities"
	"strconv"
	"time"
)

type ScheduleService interface {
    GetCheckPaymentFailedUsingScheduleService() (response.CheckingPaymentStatusResponse, error)
    DummySendNotif(userId int, schedule *request.DummyNotifRequest) error
    SendReminderDueDate() (string, error)
}


type scheduleService struct {
	scheduleRepository repositories.ScheduleRepository
	userRepositories   repositories.UserRepository
	schoolRepositories repositories.SchoolRepository
	billingRepository  repositories.BillingRepositoryInterface
}

func NewScheduleService(scheduleRepository repositories.ScheduleRepository, userRepositories repositories.UserRepository, billingRepository repositories.BillingRepositoryInterface) ScheduleService {
	return &scheduleService{scheduleRepository: scheduleRepository, userRepositories: userRepositories, billingRepository: billingRepository}
}


func (s *scheduleService) GetCheckPaymentFailedUsingScheduleService() (response.CheckingPaymentStatusResponse, error) {
	fmt.Println("start: ", time.Now())
	var resp response.CheckingPaymentStatusResponse
	var listDetail []response.CheckingPaymentStatusDetailResponse

	dataTransactionStatus, err := repositories.GetStatusPayment()
	if err != nil {
		return resp, err
	}

	var checkStatus []response.CheckPaymentStatusFailed
	if dataTransactionStatus != nil {
		checkStatus, err = repositories.GetCheckStatusPayment()
		if err != nil {
			return resp, err
		}

		// Update status transaksi untuk setiap item dalam checkStatus
		for _, transaction := range checkStatus {
			if err := repositories.UpdateTransactionStatus(transaction.ID); err != nil {
				return resp, err
			}
			paymentMethod, err := repositories.GetPaymentMethodByID(transaction.MasterPaymentMethodID)
			if err != nil {
				return resp, err
			}

			emailParent, _ := repositories.GetEmailParentById(int(transaction.StudentID))

			school, err := repositories.GetSchoolLogoSendEmailFailed(transaction.ID)
			if err != nil {
				return resp, fmt.Errorf("failed to retrieve school data")
			}
			schoolLogo := utilities.ConvertPath(school.SchoolLogo)

			// Mengonversi totalPayment dari int ke string
			totalPayment := int(transaction.TotalAmount)
			formattedTotalPayment := strconv.Itoa(totalPayment)

			formattedDate := time.Now().Format("02-01-2006")

			bodyEmail := BodyEmailTransaction{
				StudentName:    transaction.StudentName,
				InvoiceNumber:  transaction.InvoiceNumber,
				PaymentDate:    formattedDate,
				TotalPayment:   formattedTotalPayment,
				SchoolName:     school.SchoolName,
				Year:           time.Now().Year(),
				SchoolLogo:     schoolLogo,
				BankName:       paymentMethod.BankName,
				VirtualAccount: transaction.VirtualAccountNumber,
			}

			emailTemplate := utilities.GenerateEmailBodyTransactionFailed()
			err = s.SendEmailPaymentFailed(emailParent.Email, "Pembayaran Anda Dibatalkan", emailTemplate, bodyEmail)
			if err != nil {
				return resp, err
			}

			var listFirebaseToken []string
			dataAuditTrails, _ := repositories.GetDataAuditTrailByUserId(transaction.CreatedBy)
			for _, dataAuditTrail := range dataAuditTrails {
				if dataAuditTrail.FirebaseID != "" {
					exists := false
					for _, firebaseID := range listFirebaseToken {
						if firebaseID == dataAuditTrail.FirebaseID {
							exists = true
							break
						}
					}

					if !exists {
						listFirebaseToken = append(listFirebaseToken, dataAuditTrail.FirebaseID)
					}
				}
			}

			if len(listFirebaseToken) != 0 {
				for _, token := range listFirebaseToken {
					fmt.Println("send notif to ", token)
					client, err := utilities.NewFirebaseClient(constants.JsonFirebaseConfigFile)
					if err != nil {
						return resp, fmt.Errorf(constants.ErrorFirebaseClientMessage, err)
					}

					err = client.SendNotificationPayment(token, "School Payment", "Segera Melakukan Pembayaran Pembayaran Anda dibawah ini akan segera <b>kadaluwarsa</b> pada Senin, 28 Oktober 2024, 16.49", transaction.Nis, transaction.StudentName, strconv.Itoa(int(transaction.ID)), "payment_failed ", "", strconv.Itoa(transaction.CreatedBy))
					if err != nil {
						return resp, fmt.Errorf("Error sending notification: %v", err)
					}
				}
			}
		}
	}

	// Iterasi melalui dataTransactionStatus dan buat list detail respons
	for _, transaction := range dataTransactionStatus {
		listDetail = append(listDetail, response.CheckingPaymentStatusDetailResponse{
			ID:                   transaction.ID,
			TransactionType:      transaction.TransactionType,
			VirtualAccountNumber: transaction.VirtualAccountNumber,
			TotalAmount:          transaction.TotalAmount,
			ReferenceNumber:      transaction.ReferenceNumber,
			Description:          transaction.Description,
			StudentID:            int(transaction.StudentID),
			OrderID:              transaction.OrderID,
			TransactionStatus:    transaction.TransactionStatus,
			InvoiceNumber:        transaction.InvoiceNumber,
			BillingStudentIds:    transaction.BillingStudentIds,
		})
	}

	resp.Data = listDetail
	fmt.Println("finish 123: ", len(checkStatus))

	// logic send notif
	return resp, nil
}

func (s *scheduleService) SendReminderDueDate() (string, error) {
	fmt.Println("start: ", time.Now())

	billingStudents, err := repositories.GetBillingStudentReminder()

	if err != nil {
		return "", err
	}

	currentDate := time.Now().Truncate(24 * time.Hour)

	// Looping melalui setiap billing student
	for _, billingStudent := range billingStudents {
		// Menghitung selisih antara due_date dan current_date
		dueDate := billingStudent.DueDate // Asumsi student.DueDate adalah tipe time.Time
		daysRemaining := dueDate.Sub(currentDate).Hours() / 24

		emailParent, _ := repositories.GetEmailParentById(int(billingStudent.StudentID))

		// Menampilkan hasil selisih hari
		fmt.Printf("Sending reminder to %v billing student %v with due date: %v, days remaining: %.0f days\n", emailParent.Email , billingStudent.ID,  billingStudent.DueDate, daysRemaining)
		
		// kirim reminder ke student (via email)
		totalPayment := int(billingStudent.Amount)
		formattedTotalPayment := strconv.Itoa(totalPayment)

		schoolLogo := utilities.ConvertPath(billingStudent.SchoolLogo)

		bodyEmail := BodyEmailTransaction{
			StudentName:    billingStudent.FullName,
			TotalPayment:   formattedTotalPayment,
			SchoolName:     billingStudent.SchoolName,
			Year:           time.Now().Year(),
			SchoolLogo:     schoolLogo,
			ExpireDate: 	s.formatDateToIndonesian(*dueDate),
		}

		emailTemplate := utilities.GenerateEmailBodyBillingReminder()

		err = s.SendEmailPaymentFailed(emailParent.Email, "Pengingat untuk melakukan pembayaran", emailTemplate, bodyEmail)
		if err != nil {
			return "", err
		}

		err = s.SendPaymentNotif(int(emailParent.ID), bodyEmail.StudentNis, bodyEmail.StudentName, "", "Pengingat untuk melakukan pembayaran", "payment_waiting", "")
	}

	return "Email Send", nil
}

func (s *scheduleService) SendEmailPaymentFailed(email string, subject string, emailTemplate string, bodyEmail BodyEmailTransaction) error {
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

func (s *scheduleService) DummySendNotif(userId int, schedule *request.DummyNotifRequest) error {
	// Fetch unique Firebase tokens
	var listFirebaseToken []string
	dataAuditTrails, _ := repositories.GetDataAuditTrailByUserId(userId)
	tokenSet := make(map[string]bool) // Use a map to avoid duplicates

	for _, dataAuditTrail := range dataAuditTrails {
		if dataAuditTrail.FirebaseID != "" && !tokenSet[dataAuditTrail.FirebaseID] {
			tokenSet[dataAuditTrail.FirebaseID] = true
			listFirebaseToken = append(listFirebaseToken, dataAuditTrail.FirebaseID)
		}
	}

	// Return early if no tokens are available
	if len(listFirebaseToken) == 0 {
		return nil
	}

	// Initialize Firebase client once
	client, err := utilities.NewFirebaseClient(constants.JsonFirebaseConfigFile)
	if err != nil {
		return fmt.Errorf(constants.ErrorFirebaseClientMessage, err)
	}

	// Send notifications
	for _, token := range listFirebaseToken {
		fmt.Printf("Debug: Preparing to send notification - Token: %s, Title: %s, Body: %s\n",
			token, schedule.Title, schedule.Body)
		err = client.SendNotificationDummy(token, schedule, strconv.Itoa(userId))
		if err != nil {
			fmt.Printf("Error while sending notification: %v\n", err)
			// Log error but continue sending notifications to other tokens
			continue
		}
		fmt.Println("Notification sent successfully to:", token)
	}

	return nil
}

func (s *scheduleService) SendPaymentNotif(userId int, nis, studentName, transactionID, message, notificationType, redirectUrl string) error {
	// Use a map to track unique Firebase tokens
	tokenSet := make(map[string]bool)
	dataAuditTrails, _ := repositories.GetDataAuditTrailByUserId(userId)

	for _, dataAuditTrail := range dataAuditTrails {
		if dataAuditTrail.FirebaseID != "" && !tokenSet[dataAuditTrail.FirebaseID] {
			tokenSet[dataAuditTrail.FirebaseID] = true
		}
	}

	// Extract unique tokens into a slice
	var listFirebaseToken []string
	for token := range tokenSet {
		listFirebaseToken = append(listFirebaseToken, token)
	}

	// Return early if no tokens are available
	if len(listFirebaseToken) == 0 {
		return nil
	}

	// Initialize Firebase client once
	client, err := utilities.NewFirebaseClient(constants.JsonFirebaseConfigFile)
	if err != nil {
		return fmt.Errorf(constants.ErrorFirebaseClientMessage, err)
	}

	// Send notifications
	for _, token := range listFirebaseToken {
		err = client.SendNotificationPayment(token, "School Payment", message, nis, studentName, transactionID, notificationType, redirectUrl, strconv.Itoa(userId))
		if err != nil {
			// Log the error but continue sending to other tokens
			fmt.Printf("Error sending notification to %s: %v\n", token, err)
			continue
		}
		fmt.Println("Notification sent successfully to:", token)
	}

	return nil
}

func (s *scheduleService) formatDateToIndonesian(date time.Time) string {
	// Mapping nama bulan ke bahasa Indonesia
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	// Format: 1 Januari 2006
	day := date.Day()
	month := months[date.Month()-1] // time.Month dimulai dari 1 (Januari = 1)
	year := date.Year()

	return fmt.Sprintf("%d %s %d", day, month, year)
}

