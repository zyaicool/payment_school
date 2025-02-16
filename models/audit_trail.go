package models

import "time"

type AuditTrail struct {
	Master
	UserID     uint      `json:"userId"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	UserAction string    `json:"userAction"`
	ApiPath    string    `json:"apiPath"`
	LogTime    time.Time `json:"logTime"`
	Platform   string    `json:"platfrom"`
	FirebaseID string    `json:"firebasId"`
	IsValid    bool      `json:"isValid"`
}
