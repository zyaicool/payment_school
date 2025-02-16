package models

type PrefixClass struct {
	Master
	PrefixName string `json:"prefixName"`
	SchoolID   uint   `jsom:"schoolId"`
}
