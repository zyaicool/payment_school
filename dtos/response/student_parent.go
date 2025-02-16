package response

import (
	models "schoolPayment/models"
)

type StudentParentResponse struct {
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
	Data  []models.StudentParent `json:"data"`
}
