package models

type StudentParent struct {
	Master
	UserID            uint   `json:"userId"`
	ParentName        string `json:"parentName"`
	ParentAddress     string `json:"parentAddress"`
	ParentHandphone   int    `json:"parentHandphone"`
	ParentMail        string `json:"parentMail"`
	ParentCitizenship string `json:"parentCitizenship"`
	ParentSalary      string `json:"parentSalary"`
	ParentStatus      string `json:"parentStatus"`
}
