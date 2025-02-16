package response

import "time"

type AnnouncementResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	HeroImage   *string    `json:"heroImage"`
	Type        string     `json:"type"`
	EventDate   *time.Time `json:"eventDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UpdatedBy   string     `json:"updatedBy"`
}

type GetAnnouncementListResponse struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	TotalPage int                    `json:"totalPage"`
	TotalData int                    `json:"totalData"`
	Data      []AnnouncementResponse `json:"data"`
}
