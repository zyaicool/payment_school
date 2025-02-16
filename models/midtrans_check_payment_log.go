package models

type MidtransCheckPaymentLog struct {
	Master
	OrderID          string `json:"order_id"`
	StatusCode       string `json:"statusCode"`
	ResponseBodyJson string `json:"responseBodyJson"`
}
