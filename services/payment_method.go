package services

import (
	"fmt"
	"strconv"

	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
	"schoolPayment/utilities"
)

type PaymentMethodService interface {
	CreatePaymentMethod(paymentRequest *request.PaymentMethodCreateRequest, userID int) (*models.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethodID int, updateRequest *request.PaymentMethodCreateRequest) (*models.PaymentMethod, error)
	GetAllPaymentMethod(search string) ([]response.PaymentMethodResponse, error)
	GetPaymentMethodDetail(id int) (*response.PaymentMethodResponse, error)
}

type paymentMethodService struct {
	paymentMethodRepository repositories.PaymentMethodRepository
}

// Constructor untuk membuat instance service baru
func NewPaymentMethodService(paymentMethodRepository repositories.PaymentMethodRepository) PaymentMethodService {
	return &paymentMethodService{paymentMethodRepository: paymentMethodRepository}
}

func (s *paymentMethodService) CreatePaymentMethod(paymentRequest *request.PaymentMethodCreateRequest, userID int) (*models.PaymentMethod, error) {
	// Parsing nilai boolean untuk isPercentage
	isPercentage, err := strconv.ParseBool(paymentRequest.IsPercentage)
	if err != nil {
		return nil, fmt.Errorf("invalid value for is_percentage: %v", err)
	}

	// Membuat model PaymentMethod dari request
	paymentMethod := models.PaymentMethod{
		PaymentMethod:      paymentRequest.PaymentMethod,
		BankCode:           paymentRequest.BankCode,
		BankName:           paymentRequest.BankName,
		AdminFee:           paymentRequest.AdminFee,
		MethodLogo:         paymentRequest.MethodLogo,
		IsPercentage:       isPercentage,
		AdminFeePercentage: paymentRequest.AdminFeePercentage,
	}
	paymentMethod.Master.CreatedBy = userID
	paymentMethod.Master.UpdatedBy = userID

	// Menyimpan data ke repository
	dataPaymentMethod, err := s.paymentMethodRepository.CreatePaymentMethod(&paymentMethod)
	if err != nil {
		return nil, err
	}

	return dataPaymentMethod, nil
}

func (s *paymentMethodService) UpdatePaymentMethod(paymentMethodID int, updateRequest *request.PaymentMethodCreateRequest) (*models.PaymentMethod, error) {
	// Mengambil data payment method yang akan di-update
	paymentMethod, err := s.paymentMethodRepository.GetPaymentMethodByID(paymentMethodID)
	if err != nil {
		return nil, err
	}

	// Memperbarui field yang disediakan dalam request
	if updateRequest.PaymentMethod != "" {
		paymentMethod.PaymentMethod = updateRequest.PaymentMethod
	}
	if updateRequest.BankCode != "" {
		paymentMethod.BankCode = updateRequest.BankCode
	}
	if updateRequest.BankName != "" {
		paymentMethod.BankName = updateRequest.BankName
	}
	if updateRequest.MethodLogo != "" {
		paymentMethod.MethodLogo = updateRequest.MethodLogo
	}
	if updateRequest.IsPercentage != "" {
		isPercentage, err := strconv.ParseBool(updateRequest.IsPercentage)
		if err != nil {
			return nil, fmt.Errorf("invalid value for is_percentage: %v", err)
		}
		paymentMethod.IsPercentage = isPercentage
	}
	if updateRequest.AdminFeePercentage != "" {
		paymentMethod.AdminFeePercentage = updateRequest.AdminFeePercentage
	}
	paymentMethod.AdminFee = updateRequest.AdminFee

	// Menyimpan perubahan ke repository
	updatedPaymentMethod, err := s.paymentMethodRepository.UpdatePaymentMethod(paymentMethod)
	if err != nil {
		return nil, err
	}

	return updatedPaymentMethod, nil
}

func (s *paymentMethodService) GetAllPaymentMethod(search string) ([]response.PaymentMethodResponse, error) {
	// Mengambil semua metode pembayaran dari repository
	paymentMethods, err := s.paymentMethodRepository.GetAllPaymentMethod(search)
	if err != nil {
		return nil, err
	}

	// Mapping data dari model ke response
	var responseData []response.PaymentMethodResponse
	for _, paymentMethod := range paymentMethods {
		convertedLogo := utilities.ConvertPath(paymentMethod.MethodLogo)

		responseData = append(responseData, response.PaymentMethodResponse{
			ID:                 paymentMethod.ID,
			PaymentMethod:      paymentMethod.PaymentMethod,
			BankCode:           paymentMethod.BankCode,
			BankName:           paymentMethod.BankName,
			AdminFee:           paymentMethod.AdminFee,
			MethodLogo:         convertedLogo,
			IsPercentage:       paymentMethod.IsPercentage,
			AdminFeePercentage: paymentMethod.AdminFeePercentage,
		})
	}

	return responseData, nil
}

func (s *paymentMethodService) GetPaymentMethodDetail(id int) (*response.PaymentMethodResponse, error) {
	// Mengambil detail metode pembayaran berdasarkan ID
	paymentMethod, err := s.paymentMethodRepository.GetPaymentMethodByID(id)
	if err != nil {
		return nil, err
	}

	// Mapping data dari model ke response
	response := &response.PaymentMethodResponse{
		ID:                 paymentMethod.ID,
		PaymentMethod:      paymentMethod.PaymentMethod,
		BankCode:           paymentMethod.BankCode,
		BankName:           paymentMethod.BankName,
		AdminFee:           paymentMethod.AdminFee,
		MethodLogo:         utilities.ConvertPath(paymentMethod.MethodLogo),
		IsPercentage:       paymentMethod.IsPercentage,
		AdminFeePercentage: paymentMethod.AdminFeePercentage,
	}

	return response, nil
}
