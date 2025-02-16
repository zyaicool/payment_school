package services

import (
	"encoding/json"
	"fmt"
	"os"

	response "schoolPayment/dtos/response"
	"schoolPayment/repositories"
)

type PaymentTypeService interface {
	GetDataPaymentTypes(filename string) ([]response.PaymentType, error)
	GetPaymentTypeByID(filename string, paymentTypeID int) (*response.PaymentType, error)
}

// ProvinceServiceImpl is the implementation of the ProvinceService interface
type PaymentTypeServiceImpl struct{}

// NewProvinceService creates a new instance of ProvinceServiceImpl
func NewPaymentTypeService() PaymentTypeService {
	return &PaymentTypeServiceImpl{}
}

// GetDataPaymentTypes loads the payment type data from a JSON file
func (paymentTypeService *PaymentTypeServiceImpl) GetDataPaymentTypes(filename string) ([]response.PaymentType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paymentTypes []response.PaymentType
	var paymentTypesNew []response.PaymentType
	if err := json.NewDecoder(file).Decode(&paymentTypes); err != nil {
		return nil, err
	}

	for _, paymentType := range paymentTypes {
		if paymentType.ID == 99 {
			paymentTypesNew = append(paymentTypesNew,paymentType)
		}
	}

	paymentMethods, err := repositories.GetPaymentMethodForPaymentType()
	if err != nil {
		return nil, err
	}
	for _, paymentMethod := range paymentMethods {
		paymentTypesNew = append(paymentTypesNew, paymentMethod)
	}

	return paymentTypesNew, nil
}

// GetPaymentTypeByID retrieves a single payment type by its ID
func (paymentTypeService *PaymentTypeServiceImpl) GetPaymentTypeByID(filename string, paymentTypeID int) (*response.PaymentType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paymentTypes []response.PaymentType
	if err := json.NewDecoder(file).Decode(&paymentTypes); err != nil {
		return nil, err
	}

	for _, paymentType := range paymentTypes {
		if paymentType.ID == paymentTypeID {
			return &paymentType, nil
		}
	}
	return nil, fmt.Errorf("No payment type found with ID %d", paymentTypeID)
}
