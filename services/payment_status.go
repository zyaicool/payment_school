package services

import (
	"encoding/json"
	"os"
)

type PaymentStatus struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type PaymentStatusService interface {
	GetPaymentStatus(filename string) ([]PaymentStatus, error)
}

type PaymentStatusServiceImpl struct{}

func NewPaymentStatusService() PaymentStatusService {
	return &PaymentStatusServiceImpl{}
}

func (service *PaymentStatusServiceImpl) GetPaymentStatus(filename string) ([]PaymentStatus, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paymentStatuses []PaymentStatus
	if err := json.NewDecoder(file).Decode(&paymentStatuses); err != nil {
		return nil, err
	}
	return paymentStatuses, nil
}
