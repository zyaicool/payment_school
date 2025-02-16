package models

import "time"

type Announcements struct {
	Master
	SchoolID    uint       `json:"schoolId"`
	HeroImage   string     `json:"heroImage"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	EventDate   *time.Time `json:"eventDate,omitempty"` // Add the EventDate field
	School      *School    `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
}
