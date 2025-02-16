package response

import (
	"time"
)

type BillingStudentListResponse struct {
	LatestBilling []LatestBillingStudent `json:"latestBilling"`
	ListBilling   ListBillingPerStudent  `jaon:"listBilling"`
}

type ListBillingPerStudent struct {
	Page      int                         `json:"page"`
	Limit     int                         `json:"limit"`
	TotalPage int                         `json:"totalPage"`
	TotalData int64                       `json:"totalData"`
	Data      []DataListBillingPerStudent `json:"data"`
}

type LatestBillingStudent struct {
	Nis             string     `json:"nis"`
	StudentID       int        `json:"studentId"`
	StudentName     string     `json:"studentName"`
	BillingNumber   string     `json:"billingNumber"`
	SchoolClassName string     `json:"schoolClassName"`
	BillingName     string     `json:"billingName"`
	CreatedDate     time.Time  `json:"createdDate"`
	DueDate         *time.Time `json:"dueDate"`
	Semester        string     `json:"semester"`
	BillingAmount   int        `json:"billingAmount"`
	BillingStatus   string     `json:"billingStatus"`
	BillingID       int        `json:"billingId"`
}

type DataListBillingPerStudent struct {
	BillingStudentID  int    `json:"billingStudentId"`
	DetailBillingName string `json:"billingDetailName"`
	BillingType       string `json:"billingType"`
	SchoolGradeName   string `json:"schoolGrade"`
	SchoolClassName   string `json:"schoolClass"`
	StudentName       string `json:"studentName"`
	Amount            int    `json:"amount"`
}

type DetailBillingStudentResponse struct {
	ID                 uint                `json:"id"`
	Nis                string              `json:"nis"`
	RegistrationNumber string              `json:"registrationNumber"`
	StudentName        string              `json:"studentName"`
	SchoolClassName    string              `json:"schoolClassName"`
	CreatedDate        *time.Time          `json:"createdDate"`
	DueDate            *time.Time          `json:"dueDate"`
	BillingId          int                 `json:"billingId"`
	BillingName        string              `json:"billingName"`
	BillingNumber      string              `json:"billingNumber"`
	BillingTypeName    string              `json:"billingTypeName"`
	BillingAmount      int                 `json:"billingAmount"`
	Discount           int                 `json:"discount"`
	Ppn                int                 `json:"ppn"`
	Semester           string              `json:"semester"`
	PaymentType        string              `json:"paymentType"`
	SubTotal           int                 `json:"subTotal"`
	BillingStatus      string              `json:"billingStatus"`
	InstallmentDetails []InstallmentDetail `json:"installmentDetails"`
	Total              int                 `json:"total"`
	StudentID          int                 `json:"studentId"`
}

type BillingStudentByStudentIDBillingID struct {
	ID                uint       `json:"id"`
	StudentName       string     `json:"studentName"`
	SchoolClass       string     `json:"schoolClass"`
	InvoiceNumber     string     `json:"invoiceNumber"`
	TransactionStatus string     `json:"transactionStatus"`
	ChangeAmount      int64      `json:"changeAmount"`
	DiscountType      string     `json:"discountType"`
	DiscountAmount    int64      `json:"discountAmount"`
	TotalAmount       int64      `json:"totalAmount"`
	PaymentDate       *time.Time `json:"paymentDate"`
	BillingStudentIds string     `json:"billingStudentIds"`
	PaymentMethodId   int        `json:"paymentMethodId"`
	TransactionType   string     `json:"transactionType"`
}

type InstallmentDetail struct {
	BillingAmount    int     `json:"billingAmount"`
	NominalPayment   float64 `json:"nominalPayment"`
	RemainingPayment float64 `json:"remainingPayment"`
	DueDate          string  `json:"dueDate"`
}

type BillingStudentByStudentIDResponse struct {
	ID               uint                                      `json:"id"`
	Nis              string                                    `json:"nis"`
	StudentName      string                                    `json:"studentName"`
	SchoolClassName  string                                    `json:"schoolClassName"`
	SchoolGradeName  string                                    `json:"schoolGradeName"`
	LatestSchoolYear string                                    `json:"latestSchoolYear"`
	ListBilling      []BillingStudentByStudentIDDetailResponse `json:"listBilling"`
	ListDonation     []DonationBillingResponse                 `json:"listDonation"`
}

type BillingStudentIDResponse struct {
	ID              uint   `json:"id"`
	Nis             string `json:"nis"`
	StudentName     string `json:"studentName"`
	SchoolClassName string `json:"schoolClassName"`
	SchoolGradeName string `json:"schoolGradeName"`
	BillingId       int    `json:"billingId"`
	BillingName     string `json:"billingName"`
	BillingNumber   string `json:"billingNumber"`
	BillingAmount   int    `json:"billingAmount"`
	BankAccount     string `json:"bankAccount"`
}

type BillingStudentByStudentIDDetailResponse struct {
	BillingStudentID  uint       `json:"billingStudentId"`
	DetailBillingName string     `json:"detailBillingName"`
	Amount            int        `json:"amount"`
	DueDate           *time.Time `json:"dueDate"`
	BillingType       string     `json:"billingType"`
	PaymentStatus     string     `json:"paymentStatus"`
}

type BillingStudentDetailResponse struct {
	BillingStudentId  int    `json:"billingStudentId"`
	DetailBillingName string `json:"detailBillingName"`
	DueDate           string `json:"dueDate"`
	Amount            int    `json:"amount"`
}

type BillingStudentUpdateResponse struct {
	Data    BillingStudentDetailResponse `json:"data"`
	Message string                       `json:"message"`
}

type ListBillingHistory struct {
	Page      int                      `json:"page"`
	Limit     int                      `json:"limit"`
	TotalPage int                      `json:"totalPage"`
	TotalData int64                    `json:"totalData"`
	Data      []DataListBillingHistory `json:"data"`
}

type DataListBillingHistory struct {
	ID                int        `json:"id"`
	InvoiceNumber     string     `json:"invoiceNumber"`
	StudentName       string     `json:"studentName"`
	PaymentDate       *time.Time `json:"paymentDate"`
	PaymentMethod     string     `json:"paymentMethod"`
	Username          string     `json:"username"`
	TotalAmount       int64      `json:"totalAmount"`
	TransactionStatus string     `json:"transactionStatus"`
	OrderID           string     `json:"orderId"`
	Token             string     `json:"token"`
	RedirectUrl       string     `json:"redirectUrl"`
}

type BillingHistoryListResponse struct {
	Page      int                      `json:"page"`
	Limit     int                      `json:"limit"`
	TotalData int64                    `json:"total_data"`
	TotalPage int                      `json:"total_page"`
	Data      []DataListBillingHistory `json:"data"`
}

type DetailBillingHistory struct {
	ID                         uint                       `json:"id"`
	StudentName                string                     `json:"studentName"`
	SchoolClass                string                     `json:"schoolClass"`
	InvoiceNumber              string                     `json:"invoiceNumber"`
	TransactionStatus          string                     `json:"transactionStatus"`
	ChangeAmount               int64                      `json:"changeAmount"`
	DiscountAmount             int64                      `json:"discountAmount"`
	TotalBillingBeforeDiscount int64                      `json:"totalBillingBeforeDiscount"`
	TotalBillingAmount         int64                      `json:"totalBillingAmount"`
	TotalPayAmount             int64                      `json:"totalPayAmount"`
	PaymentDate                *time.Time                 `json:"paymentDate"`
	AdminFee                   int64                      `json:"adminFee"`
	PaymentMethod              string                     `json:"paymentMethod"`
	ListBilling                []BillingStudentForHistory `json:"listBilling"`
}

type BillingStudentForHistory struct {
	Amount            int64     `json:"amount"`
	DueDate           time.Time `json:"dueDate"`
	BillingType       string    `json:"billingType"`
	DetailBillingName string    `json:"detailBillingName"`
}

type RedirectUrlData struct {
	OrderId     string
	Token       string
	RedirectUrl string
}

type BillingByStudentResponse struct {
	BillingExist bool                 `json:"billingExist"`
	Data         []BillingStudentData `json:"data"`
}

type BillingStudentData struct {
	BillingID       uint                       `json:"billingId"`
	BillingName     string                     `json:"billingName"`
	BillingStudents []BillingStudentDetailData `json:"billingStudents"`
}

type BillingStudentDetailData struct {
	BillingDetailID   uint   `json:"billingDetailId"`
	BillingDetailName string `json:"billingDetailName"`
	Amount            int64  `json:"amount"`
	IsExist           bool   `json:"isExist"`
	Disabled          bool   `json:"disabled"`
}

type DonationBillingResponse struct {
	BillingID   uint   `json:"billingId"`
	BillingName string `json:"billingName"`
}
