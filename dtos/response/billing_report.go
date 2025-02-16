package response

import "math/big"

type BillingReport struct {
	DetailBillingName string `json:"detailBillingName"`
	BillingType       string `json:"billingType"`
	StudentName       string `json:"studentName"`
	SchoolGradeName   string `json:"schoolGradeName"`
	SchoolClassName   string `json:"schoolClassName"`
	SchoolYearName    string `json:"schoolYearName"`
	Amount            int64  `json:"amount"`
	BankAccountId     int    `json:"bankAccountId"`
	BankAccountName   string `json:"bankAccountName"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	PaymentStatus     string `json:"payment_status"`
}

type BillingReportDetail struct {
	DetailBillingName string `json:"detailBillingName"`
	BillingType       string `json:"billingType"`
	StudentName       string `json:"studentName"`
	SchoolGradeName   string `json:"schoolGradeName"`
	SchoolClassName   string `json:"schoolClassName"`
	SchoolYearName    string `json:"schoolYearName"`
	Amount            int64  `json:"amount"`
	BankAccountName   string `json:"bankAccountName"`
	PaymentStatus     string `json:"paymentStatus"`
}

type ListReportBilling struct {
	Page      int                   `json:"page"`
	Limit     int                   `json:"limit"`
	TotalPage int                   `json:"totalPage"`
	TotalData int64                 `json:"totalData"`
	Data      []BillingReportDetail `json:"data"`
}

type BillingReportResponse struct {
	// TotalBillingAmount is the total billing amount (string representation of big.Int)
	// @swagger:strfmt string
	TotalBillingAmount *big.Int `json:"totalBillingAmount"`
	// TotalPayAmount is the total pay amount (string representation of big.Int)
	// @swagger:strfmt string
	TotalPayAmount *big.Int `json:"totalPayAmount"`
	// TotalNotPayAmount is the total not pay amount (string representation of big.Int)
	// @swagger:strfmt string
	TotalNotPayAmount  *big.Int          `json:"totalNotPayAmount"`
	TotalBillingPay    int               `json:"totalBillingPay"`
	TotalBillingNotPay int               `json:"totalBillingNotPay"`
	TotalStudent       int               `json:"totalStudent"`
	ListBillingReport  ListReportBilling `json:"listBillingReport"`
}

type BillingReportSummary struct {
	TotalBillingAmount string `json:"totalBillingAmount"`
	TotalPayAmount     string `json:"totalPayAmount"`
	TotalNotPayAmount  string `json:"totalNotPayAmount"`
	TotalBillingPay    int    `json:"totalBillingPay"`
	TotalBillingNotPay int    `json:"totalBillingNotPay"`
	TotalStudent       int    `json:"totalStudent"`
}

// for export Excel
type BillingReportExportDTO struct {
	TotalBillingAmount float64          `json:"totalBillingAmount"`
	TotalPaidAmount    float64          `json:"totalPaidAmount"`
	TotalUnpaidAmount  float64          `json:"totalUnpaidAmount"`
	TotalPaidCount     int              `json:"totalPaidCount"`
	TotalUnpaidCount   int              `json:"totalUnpaidCount"`
	TotalStudents      int              `json:"totalStudents"`
	BillingData        []BillingDataDTO `json:"billingData"`
}

type BillingDataDTO struct {
	StudentName       string  `json:"studentName"`
	SchoolClass       string  `json:"schoolClass"`
	InvoiceNumber     string  `json:"invoiceNumber"`
	TransactionStatus string  `json:"transactionStatus"`
	Amount            float64 `json:"amount"`
	DueDate           string  `json:"dueDate"`
}
