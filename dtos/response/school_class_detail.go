package response

type SchoolClassDetailResponse struct {
	ID              int    `json:"id"`
	SchoolGradeID   int    `json:"school_grade_id"`
	PrefixClassID   int    `json:"prefix_class_id"`
	SchoolMajorID   int    `json:"school_major_id"`
	Suffix          string `json:"suffix"`
	SchoolClassName string `json:"school_class_name"`
}
