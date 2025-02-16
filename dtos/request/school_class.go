package request

type SchoolClassCreateUpdateRequest struct {
	SchoolID        int    `json:"schoolId"`
	SchoolGradeID   int    `json:"schoolGradeId"`
	PrefixClassID   int    `json:"prefixClassId"`
	SchoolMajorID   int    `json:"schoolMajorId"`
	Suffix          string `json:"suffix"`
	SchoolClassName string `json:"schoolClassName"`
}
