package request

type AnnouncementCreateRequest struct {
	SchoolID    int    `json:"schoolId"`
	HeroImage   string `json:"heroImage"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	EventDate   string `json:"eventDate"`
}
type AnnouncementUpdateRequest struct {
	SchoolID    *int   `json:"schoolId,omitempty"`
	HeroImage   string `json:"heroImage,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	EventDate   string `json:"eventDate,omitempty"` // EventDate remains a string
}