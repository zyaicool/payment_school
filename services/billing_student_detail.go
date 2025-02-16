package services

import (
	"schoolPayment/repositories"
)

type BillingStudentDetailService struct {
	billingStudentRepository repositories.BillingStudentRepository
}

func NewBillingStudentDetailService(billingStudentRepository repositories.BillingStudentRepository) BillingStudentService {
	return BillingStudentService{
		billingStudentRepository: billingStudentRepository,
	}
}
