package services

import (
	"database/sql"
	"fmt"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"
	"strconv"
	"time"
)

type AnnouncementService interface {
	GetAnnouncementTypes() ([]repositories.AnnouncementType, error)
	CreateAnnouncement(announcementRequest *request.AnnouncementCreateRequest, userID int) (*models.Announcements, error)
	DeleteAnnouncement(id uint, userID int) error
	UpdateAnnouncement(announcementID int, updateRequest *request.AnnouncementUpdateRequest, userID int) (*models.Announcements, error)
	GetListAnnouncement(page, limit int, search, sortBy, sortOrder, announcementType string, userID int) (response.GetAnnouncementListResponse, error)
	GetAnnouncementByID(id int) (response.AnnouncementResponse, error)
}

type announcementService struct {
	repo     repositories.AnnouncementRepository
	userRepo repositories.UserRepository
	db       *sql.DB
}

func NewAnnouncementService(repo repositories.AnnouncementRepository, userRepo repositories.UserRepository) AnnouncementService {
	return &announcementService{repo: repo, userRepo: userRepo}
}

func (s *announcementService) GetAnnouncementTypes() ([]repositories.AnnouncementType, error) {
	return s.repo.GetAllAnnouncement()
}

func (announcementService *announcementService) CreateAnnouncement(announcementRequest *request.AnnouncementCreateRequest, userID int) (*models.Announcements, error) {
	var eventDate time.Time
	var err error
	announcement := models.Announcements{
		SchoolID:    uint(announcementRequest.SchoolID),
		HeroImage:   announcementRequest.HeroImage,
		Title:       announcementRequest.Title,
		Description: announcementRequest.Description,
		Type:        announcementRequest.Type,
		EventDate:   nil,
	}
	if announcementRequest.EventDate != "" {
		eventDate, err = time.Parse("2006-01-02 15:04", announcementRequest.EventDate) // Expected format: YYYY-MM-DD hh:mm
		if err != nil {
			return nil, fmt.Errorf("invalid eventDate format. Expected YYYY-MM-DD: %v", err)
		}
		announcement.EventDate = &eventDate
	}

	// Validate the event date based on the type
	// if announcementRequest.Type == "AT01" && announcementRequest.EventDate == "" {
	//  // EventDate is required for AT01 type
	//  return nil, fmt.Errorf("eventDate is required for type AT01 (Event)")
	// }

	announcement.Master.CreatedBy = userID
	announcement.Master.UpdatedBy = userID

	// Call repository to save the announcement
	dataAnnouncement, err := announcementService.repo.CreateAnnouncement(&announcement)
	if err != nil {
		return nil, err
	}

	image := ""
	if dataAnnouncement.HeroImage != "" {
		image = utilities.ConvertPath(dataAnnouncement.HeroImage)
	}

	typeAnnouncement := "Pengumuman"
	if dataAnnouncement.Type == "AT01" {
		typeAnnouncement = "Event"
	}

	SendAnnouncementNotif(
		announcementRequest.SchoolID,
		dataAnnouncement.Title,
		dataAnnouncement.Description,
		image, typeAnnouncement,
		strconv.Itoa(int(dataAnnouncement.ID)))

	return dataAnnouncement, nil
}

func SendAnnouncementNotif(schoolId int, title, body, image, typeAnnouncement, announcementID string) error {
	// Use a map to track unique Firebase tokens
	tokenSet := make(map[string]bool)
	dataAuditTrails, _ := repositories.GetAllAuditrailsNotificationMobile(schoolId)

	for _, dataAuditTrail := range dataAuditTrails {
		if dataAuditTrail.FirebaseID != "" && !tokenSet[dataAuditTrail.FirebaseID] {
			tokenSet[dataAuditTrail.FirebaseID] = true
		}
	}

	// Extract unique tokens into a slice
	var listFirebaseToken []string
	for token := range tokenSet {
		listFirebaseToken = append(listFirebaseToken, token)
	}

	// Return early if no tokens are available
	if len(listFirebaseToken) == 0 {
		return nil
	}

	// Initialize Firebase client once
	client, err := utilities.NewFirebaseClient("./data/firebase_config_file.json")
	if err != nil {
		return fmt.Errorf("Error initializing Firebase client: %v", err)
	}

	// Send notifications
	for _, token := range listFirebaseToken {
		fmt.Println("Sending notification to:", token)
		err = client.SendingAnouncementNotification(token, title, body, image, typeAnnouncement, announcementID)
		if err != nil {
			// Log error and continue sending to other tokens
			fmt.Printf("Error sending notification to %s: %v\n", token, err)
			continue
		}
		fmt.Println("Notification sent successfully to:", token)
	}

	return nil
}

func (s *announcementService) DeleteAnnouncement(id uint, userID int) error {
	return s.repo.DeleteAnnouncement(id, userID)
}

func (announcementService *announcementService) UpdateAnnouncement(announcementID int, updateRequest *request.AnnouncementUpdateRequest, userID int) (*models.Announcements, error) {
	// Fetch the existing announcement
	existingAnnouncement, err := announcementService.repo.FindByID(announcementID)
	if err != nil {
		return nil, fmt.Errorf("announcement not found: %v", err)
	}

	// Update fields only if they are provided
	if updateRequest.SchoolID != nil {
		existingAnnouncement.SchoolID = uint(*updateRequest.SchoolID)
	}

	if updateRequest.HeroImage != "" {
		existingAnnouncement.HeroImage = updateRequest.HeroImage
	}

	if updateRequest.Title != "" {
		existingAnnouncement.Title = updateRequest.Title
	}

	if updateRequest.Description != "" {
		existingAnnouncement.Description = updateRequest.Description
	}

	if updateRequest.EventDate != "" {
		// Parse and validate eventDate
		eventDate, err := time.Parse("2006-01-02 15:04", updateRequest.EventDate)
		if err != nil {
			return nil, fmt.Errorf("invalid eventDate format. Expected YYYY-MM-DD hh:mm: %v", err)
		}
		existingAnnouncement.EventDate = &eventDate
	}

	// Update audit fields
	existingAnnouncement.Master.UpdatedBy = userID

	// Save the updated announcement
	updatedAnnouncement, err := announcementService.repo.UpdateAnnouncement(existingAnnouncement)
	if err != nil {
		return nil, err
	}

	return updatedAnnouncement, nil
}

func (s *announcementService) GetListAnnouncement(page, limit int, search, sortBy, sortOrder, announcementType string, userID int) (response.GetAnnouncementListResponse, error) {
	if sortBy != "" {
		sortBy = utilities.ChangeStringSortByAnnouncement(sortBy)
	}

	user, err := s.userRepo.GetUserByID(uint(userID))
	if err != nil {
		return response.GetAnnouncementListResponse{}, err
	}

	announcements, totalPage, totalData, err := s.repo.GetListAnnouncement(page, limit, search, sortBy, sortOrder, announcementType, user)
	if err != nil {
		return response.GetAnnouncementListResponse{}, err
	}

	for i := range announcements {

		if announcements[i].HeroImage != nil {
			tempHeroImage := utilities.ConvertPath(*announcements[i].HeroImage)
			announcements[i].HeroImage = &tempHeroImage
		}

		// if announcements[i].HeroImage == nil || *announcements[i].HeroImage == "" || *announcements[i].HeroImage == "https:////api/assets/" {
		// 	announcements[i].HeroImage = nil
		// }

		if announcements[i].Type == "AT01" {
			announcements[i].Type = "Event"
		} else {
			announcements[i].Type = "Pengumuman"
		}
	}

	responseDTO := response.GetAnnouncementListResponse{
		Page:      page,
		Limit:     limit,
		TotalPage: totalPage,
		TotalData: totalData,
		Data:      announcements,
	}

	return responseDTO, nil
}

func (s *announcementService) GetAnnouncementByID(id int) (response.AnnouncementResponse, error) {

	announcement, err := s.repo.GetAnnouncementByID(id)
	if err != nil {
		return response.AnnouncementResponse{}, err
	}

	if announcement.HeroImage != nil {
		tempHeroImage := utilities.ConvertPath(*announcement.HeroImage)
		announcement.HeroImage = &tempHeroImage
	}

	// if announcement.HeroImage == nil || *announcement.HeroImage == "" || *announcement.HeroImage == "https:////api/assets/" {
	//     announcement.HeroImage = nil
	// }

	if announcement.Type == "AT01" {
		announcement.Type = "Event"
	} else {
		announcement.Type = "Pengumuman"
	}
	return announcement, nil
}
