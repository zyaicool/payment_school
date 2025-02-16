package models

type MidtransPaymentLog struct {
	Master
	OrderID          string `json:"orderId"`
	StatusCode       string `json:"statusCode"`
	Token            string `json:"token"`
	RedirectUrl      string `json:"redirectUrl"`
	RequestBodyJson  string `json:"requestBodyJson"`
	ResponseBodyJson string `json:"responseBodyJson"`
}
