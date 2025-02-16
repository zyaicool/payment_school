package request

type DummyNotifRequest struct {
	Title            string `json:"title"`
	Body             string `json:"body"`
	Nis              string `json:"nis"`
	StudentName      string `json:"studentName"`
	TransactionID    string `json:"transactionId"`
	Image            string `json:"image"`
	Type             string `json:"type"`
	AnnouncementId   string `json:"announcementId"`
	NotificationType string `json:"notificationType"`
	RedirectUrl      string `json:"redirectUrl"`
}
