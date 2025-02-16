package response

import (
	models "schoolPayment/models"
)

type SchoolGradeListResponse struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Data  []models.SchoolGrade `json:"data"`
}
