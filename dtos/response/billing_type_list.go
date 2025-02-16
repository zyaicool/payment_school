package response

import (
	models "schoolPayment/models"
)

type BillingTypeListResponse struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Data  []models.BillingType `json:"data"`
}
