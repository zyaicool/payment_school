package request

type PrefixClassCreate struct {
	PrefixName string `json:"prefixName"`
	SchoolID   uint   `json:"schoolId"`
}
