package response

type DashboardResponse struct {
	DataDashboard []DetailDashboardResponse `json:"dataDashboard"`
}

type DetailDashboardResponse struct {
	DetailDataStudent StudentResponse             `json:"detailDataStudent"`
	ListLatestBilling []ListLatestBillingResponse `json:"listLatestBilling"`
	DonationName      string                      `json:"donationName"`
}

type StudentResponse struct {
	ID          uint   `json:"id"`
	Nis         string `json:"nis"`
	FullName    string `json:"fullName"`
	Status      string `json:"status"`
	SchoolClass string `json:"schoolClass"`
	SchoolName  string `json:"schoolName"`
}

type ListLatestBillingResponse struct {
	BillingName string `json:"billingName"`
	StudentName string `json:"studentName"`
	DueDate     string `json:"dueDate"` // Use string if you want formatted date, otherwise `time.Time`
	Amount      int    `json:"amount"`
	Status      string `json:"status"`
}

type DashboardAdminResponse struct {
	ActiveStudent    int `json:"activeStudent"`
	TotalTransaction int `json:"totalTransaction"`
	TotalBilling     int `json:"totalBilling"`
}
