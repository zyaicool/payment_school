package response

import (
	"time"

	models "schoolPayment/models"
)

type DataInvoice struct {
	InvoiceNumber     string     `json:"invoiceNumber"`
	PaymentDate       *time.Time `json:"paymentDate"`
	PrintDate         *time.Time `json:"printDate"`
	TransactionType   string     `json:"transctionType"`
	Nis               string     `json:"nis"`
	StudentName       string     `json:"studentName"`
	SchoolClassName   string     `json:"schoolClassName"`
	BillingStudentIds string     `json:"billingStudentIds"`
	SchoolID          uint       `json:"schoolId"`
	SubTotal		  int64		 `json:"subTotal"`
	Discount			  int64		 `json:"discount"`
	TotalAmount       int64      `json:"totalAmount"`
}

type RespDataInvoice struct {
	InvoiceNumber   string                  `json:"invoiceNumber"`
	PaymentDate     *time.Time              `json:"paymentDate"`
	PrintDate       *time.Time              `json:"printDate"`
	TransactionType string                  `json:"transctionType"`
	Nis             string                  `json:"nis"`
	StudentName     string                  `json:"studentName"`
	SchoolClassName string                  `json:"schoolClassName"`
	SchoolID        uint                    `json:"schoolId"`
	SubTotal		int64					`json:"subTotal"`
	Discount		int64					`json:"diskon"`
	TotalAmount     int64                   `json:"totalAmount"`
	BillingStudents []models.BillingStudent `json:"billingStudents"`
}
