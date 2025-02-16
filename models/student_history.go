package models

type StudentHistory struct {
	Master
	StudentID uint   `json:"studentId"`
	NewData   string `json:"newData"`
	OldData   string `json:"oldData"`
	Action    string `json:"action"`
}
