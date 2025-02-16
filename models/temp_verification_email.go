package models

type TempVerificationEmail struct {
	Master
	UserID           uint    `json:"userId"`
	VerificationLink string  `json:"verificationLink"`
	IsValid          bool    `json:"isValid"`
	Type             *string `json:"type"` // Pointer for nullable type
}
