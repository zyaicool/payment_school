package models

import "time"

type Master struct {
	ID        uint       `json:"id" gorm:"primaryKey;"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy int        `json:"createdBy"`
	UpdatedAt time.Time  `json:"updatedAt"`
	UpdatedBy int        `json:"updatedBy"`
	DeletedAt *time.Time `json:"deletedAt"`
	DeletedBy *int       `json:"deletedBy"`
}
