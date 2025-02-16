package request

type BillingCreateRequest struct {
	BillingType    string           `json:"billingType"`
	SchoolGradeID  int              `json:"schoolGradeId"`
	SchoolYearId   int              `json:"schoolYearId"`
	BillingName    string           `json:"billingName"`
	BillingAmount  int64            `json:"billingAmount"`
	Description    string           `json:"description"`
	BillingCode    string           `json:"billingCode"`
	SchoolClassIds []string         `json:"schoolClassIds"`
	BankAccountId  int              `json:"bankAccountId"`
	DetailBillings []DetailBillings `json:"detailBillings"`
}

type BillingUpdateRequest struct {
	BillingName   string `json:"billingName"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	PaymentType   string `json:"paymentType"`
	Tenor         string `json:"tenor"`
	BillingAmount int    `json:"billingAmount"`
	Ppn           int    `json:"ppn"`
	Discount      int    `json:"discount"`
}

type DetailBillings struct {
	DetailBillingName string `json:"detailBillingName"`
	Amount            int64  `json:"amount"`
	DueDate           string `json:"dueDate"`
}

type UpdateBillingStudentRequest struct {
	DetailBillingName string `json:"detailBillingName"`
	DueDate           string `json:"dueDate"`
	Amount            int64  `json:"amount"`
}

type CreateBillingStudentRequest struct {
	StudentID         int                        `json:"studentId"`
	DetailDataBilling []DetailDataBillingRequest `json:"detailDataBilling"`
}

type DetailDataBillingRequest struct {
	BillingID        int    `json:"billingId"`
	BillingDetailIds string `json:"billingDetailIds"`
}
