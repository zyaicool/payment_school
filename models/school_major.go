package models

type SchoolMajor struct {
	Master
	SchoolMajorName string `json:"schoolMajorName"`
	SchoolID        uint   `jsom:"schoolId"`
}
