package request

import (
	"fmt"
	"time"
)

type PaymentReportRequest struct {
	Page          int        `query:"page"`
	Limit         int        `query:"limit"`
	SortBy        string     `query:"sortBy"`
	SortOrder     string     `query:"sortOrder"`
	PaymentTypeId int        `query:"paymentTypeId"`
	UserId        int        `query:"userId"`
	StartDate     *time.Time `query:"startDate"`
	EndDate       *time.Time `query:"endDate"`
	StudentId     int        `query:"studentId"`
}

func (r *PaymentReportRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.SortOrder == "" {
		r.SortOrder = "asc"
	}
}

func (r *PaymentReportRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}
	return nil
}
