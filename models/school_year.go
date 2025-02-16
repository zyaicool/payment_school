package models

import "time"

type SchoolYear struct {
	Master
	SchoolYearCode string     `json:"schoolYearCode"`
	SchoolYearName string     `json:"schoolYearName"`
	SchoolId       int        `json:"schoolId"`
	StartDate      *time.Time `json:"startDate"`
	EndDate        *time.Time `json:"endDate"`
}
type SchoolYearList struct {
	SchoolYear                 // Embeds the original Billing struct
	CreateByUsername string `gorm:"column:create_by_username"`
}

