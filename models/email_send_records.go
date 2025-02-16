package models

import (
	"time"
)

type EmailSendRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Timestamp time.Time `gorm:"type:timestamp;not null"`
}
