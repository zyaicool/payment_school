package response

import "time"

type BillingStudentReminderList struct {
	ID          int       	`json:"id"`
	DueDate 	*time.Time	`json:"dueDate"`
	Amount   	int64    	`json:"amount"` 
	StudentID   int64		`json:"studentId"`
	FullName 	string    	`json:"fullName"`
	SchoolName	string  	`json:"schoolName"`
	SchoolLogo  string 		`json:"schoolLogo"`
	UserID		int64		`json:"userId"`
}