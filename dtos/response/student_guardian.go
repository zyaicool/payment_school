package response

import (
	models "schoolPayment/models"
)

type StudentGuardianResponse struct {
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
	Data  []models.StudentGuardian `json:"data"`
}
