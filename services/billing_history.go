package services

import (
	"fmt"
	"strconv"
	"strings"

	"schoolPayment/constants"
	response "schoolPayment/dtos/response"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

type BillingHistoryServiceInterface interface {
	GetAllBillingHistory(page int, limit int, search string, studentID int, roleID int, schoolYearId int, paymentTypeId int, paymentStatusCode string, userID int, userLoginID int, sortBy string, sortOrder string) (response.ListBillingHistory, error)
	GetDetailBillingHistoryIDService(transactionId int, userID int) (*response.DetailBillingHistory, error)
	GenerateInvoice(c *fiber.Ctx, invoiceNumber string, isPrint bool, userID int) (string, error)
}

type BillingHistoryService struct {
	billingHistoryRepository  repositories.BillingHistoryRepositoryInterface
	userRepository            repositories.UserRepository
	schoolRepositories        repositories.SchoolRepository
	paymentMethodRepositories repositories.PaymentMethodRepository
}

func NewBillingHistoryService(
	billingHistoryRepository repositories.BillingHistoryRepositoryInterface,
	userRepository repositories.UserRepository,
	schoolRepositories repositories.SchoolRepository,
	paymentMethodRepositories repositories.PaymentMethodRepository,
) BillingHistoryServiceInterface {
	return &BillingHistoryService{
		billingHistoryRepository:  billingHistoryRepository,
		userRepository:            userRepository,
		schoolRepositories:        schoolRepositories,
		paymentMethodRepositories: paymentMethodRepositories,
	}
}
func (billingHistoryService *BillingHistoryService) GetAllBillingHistory(page int, limit int, search string, studentID int, roleID int, schoolYearId int, paymentTypeId int, paymentStatusCode string, userID int, userLoginID int, sortBy string, sortOrder string) (response.ListBillingHistory, error) {
	var listBilling response.ListBillingHistory
	var schoolID int
	listBilling.Page = page
	listBilling.Limit = limit
	listBilling.TotalData = 0
	listBilling.TotalPage = 0

	user, err := billingHistoryService.userRepository.GetUserByID(uint(userLoginID))
	if err != nil {
		return listBilling, err
	}

	if user.UserSchool != nil {
		schoolID = int(user.UserSchool.SchoolID)
	}

	// Handle sorting logic
	if sortBy != "" {
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	listBillingPerStudent, totalPage, totalData, err := billingHistoryService.billingHistoryRepository.GetAllBillingHistory(page, limit, search, studentID, roleID, schoolYearId, paymentTypeId, schoolID, paymentStatusCode, sortBy, sortOrder, userID, userLoginID)
	if err != nil {
		listBilling.Data = []response.DataListBillingHistory{}
		return listBilling, err
	}

	// Map order IDs to retrieve redirect URLs
	orderIds := make([]string, len(listBillingPerStudent))
	for i, billing := range listBillingPerStudent {
		orderIds[i] = billing.OrderID
	}

	// Retrieve redirect URLs for the given order IDs
	redirectUrls, err := repositories.GetRedirectUrlsByOrderIds(orderIds)
	if err != nil {
		listBilling.Data = []response.DataListBillingHistory{}
		return listBilling, err
	}

	// Map redirect URLs to each billing record
	for i, billing := range listBillingPerStudent {
		totalAmount, err := billingHistoryService.billingHistoryRepository.TotalAmountBillingStudent(billing.ID)
		if err != nil {
			listBilling.Data = []response.DataListBillingHistory{}
			return listBilling, err
		}

		listBillingPerStudent[i].TotalAmount = totalAmount
		if redirectData, exists := redirectUrls[billing.OrderID]; exists {
			listBillingPerStudent[i].RedirectUrl = redirectData.RedirectUrl
			listBillingPerStudent[i].Token = redirectData.Token
		}
	}

	listBilling.Data = listBillingPerStudent
	listBilling.TotalData = totalData
	listBilling.TotalPage = totalPage

	return listBilling, nil
}

func (billingHistoryService *BillingHistoryService) GetDetailBillingHistoryIDService(transactionId int, userID int) (*response.DetailBillingHistory, error) {
	getBillingId, err := billingHistoryService.billingHistoryRepository.GetDetailBillingHistoryIDRepositories(transactionId)
	if err != nil {
		return nil, err
	}

	paymentMethod, _ := billingHistoryService.paymentMethodRepositories.GetPaymentMethodByID(getBillingId.PaymentMethodId)

	billingAmount, err := billingHistoryService.billingHistoryRepository.TotalAmountBillingStudent(int(getBillingId.ID))
	if err != nil {
		return nil, err
	}

	adminFee := int64(0)
	paymentMethodStr := "Kasir" // Default for PT01

	if getBillingId.TransactionType == "PT02" {
		adminFee, err = utilities.CalculateAdminFee(billingAmount, paymentMethod)
		if err != nil {
			return nil, err
		}
		paymentMethodStr = paymentMethod.PaymentMethod + "-" + paymentMethod.BankName
	}

	response := response.DetailBillingHistory{
		ID:                 getBillingId.ID,
		StudentName:        getBillingId.StudentName,
		SchoolClass:        getBillingId.SchoolClass,
		InvoiceNumber:      getBillingId.InvoiceNumber,
		TransactionStatus:  getBillingId.TransactionStatus,
		ChangeAmount:       getBillingId.ChangeAmount,
		TotalBillingAmount: 0,
		TotalPayAmount:     getBillingId.TotalAmount,
		PaymentDate:        getBillingId.PaymentDate,
		AdminFee:           adminFee,
		PaymentMethod:      paymentMethodStr,
	}

	ids := strings.Split(getBillingId.BillingStudentIds, ",")
	for _, id := range ids {
		idInt, _ := strconv.Atoi(id)
		listBilling, err := billingHistoryService.billingHistoryRepository.GetInstallmentHistoryDetails(idInt)
		if err != nil {
			return nil, err
		}
		response.TotalBillingAmount += listBilling.Amount
		response.ListBilling = append(response.ListBilling, listBilling)
		response.TotalBillingBeforeDiscount += listBilling.Amount
	}

	if getBillingId.DiscountAmount != 0 {
		var discountAmount int64
		discountAmount = 0
		if getBillingId.DiscountType == "%" {
			discountAmount = getBillingId.DiscountAmount
		} else {
			discountAmount = getBillingId.DiscountAmount
		}
		response.DiscountAmount = discountAmount
		response.TotalBillingAmount = response.TotalBillingAmount - discountAmount
	}

	return &response, nil
}

func (billingHistoryService *BillingHistoryService) GenerateInvoice(c *fiber.Ctx, invoiceNumber string, isPrint bool, userID int) (string, error) {
	user, err := billingHistoryService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return "", err
	}

	dataInvoice, err := billingHistoryService.billingHistoryRepository.GetDataForInvoice(invoiceNumber, int(user.UserSchool.SchoolID))
	if err != nil || len(dataInvoice) == 0 {
		return "", fmt.Errorf(constants.DataNotFoundMessage)
	}

	school, err := billingHistoryService.schoolRepositories.GetSchoolByID(dataInvoice[0].SchoolID)
	if err != nil {
		return "", err
	}

	filename, err := utilities.GeneratePDF(c, dataInvoice, school, isPrint)
	if err != nil {
		return "", err
	}

	return filename, nil
}
