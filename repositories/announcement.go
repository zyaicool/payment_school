package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"time"

	"gorm.io/gorm"
)

type AnnouncementType struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type AnnouncementRepository interface {
	CreateAnnouncement(announcement *models.Announcements) (*models.Announcements, error)
	GetAllAnnouncement() ([]AnnouncementType, error)
	DeleteAnnouncement(id uint, userID int) error
	UpdateAnnouncement(announcement *models.Announcements) (*models.Announcements, error)
	FindByID(id int) (*models.Announcements, error)
	GetListAnnouncement(page int, limit int, search, sortBy, sortOrder, announcementType string, user models.User) ([]response.AnnouncementResponse, int, int, error)
	GetAnnouncementByID(id int) (response.AnnouncementResponse, error)
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (announcementRepository *announcementRepository) CreateAnnouncement(announcement *models.Announcements) (*models.Announcements, error) {
	result := announcementRepository.db.Create(&announcement)
	if result.Error != nil {
		return nil, result.Error
	}
	return announcement, nil
}

func (announcementRepository *announcementRepository) GetAllAnnouncement() ([]AnnouncementType, error) {
	data, err := os.ReadFile("data/anouncement_type.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return nil, err
	}

	var announcementTypes []AnnouncementType
	if err := json.Unmarshal(data, &announcementTypes); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	return announcementTypes, nil
}

func (announcementRepository *announcementRepository) DeleteAnnouncement(id uint, userID int) error {
	result := announcementRepository.db.Model(&models.Announcements{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"deleted_by": userID,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("announcement not found")
	}

	return nil
}

func (announcementRepository *announcementRepository) UpdateAnnouncement(announcement *models.Announcements) (*models.Announcements, error) {
	result := announcementRepository.db.Save(&announcement)
	return announcement, result.Error
}

func (announcementRepository *announcementRepository) FindByID(id int) (*models.Announcements, error) {
	var announcement models.Announcements
	result := announcementRepository.db.First(&announcement, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &announcement, nil
}

func (repositories *announcementRepository) GetListAnnouncement(page int, limit int, search, sortBy, sortOrder, announcementType string, user models.User) ([]response.AnnouncementResponse, int, int, error) {

	var announcements []response.AnnouncementResponse
	var totalData int64

	offset := (page - 1) * limit

	query := database.DB.Table("announcements as a").
		Select("a.id, a.title, a.description, case when a.hero_image = '' then null else a.hero_image end, a.type, a.event_date, a.created_at, (select username from users u where u.id = a.created_by ) as created_by, a.updated_at, (select username from users u where u.id = a.updated_by ) as updated_by")

	if user.UserSchool != nil {
		query = query.Joins("JOIN user_schools us ON us.user_id = a.created_by").
			Where("a.deleted_at IS NULL AND us.school_id = ?", user.UserSchool.SchoolID)
	}

	if search != "" {
		query = query.Where("a.title LIKE ?", "%"+search+"%")
	}

	if announcementType != "" {
		var conditionType string
		if announcementType == "Event" {
			conditionType = "AT01"
		} else {
			conditionType = "AT02"
		}

		query = query.Where("a.type LIKE ?", "%"+conditionType+"%")
	}

	var countQuery = query
	countQuery.Count(&totalData)

	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("a.created_at DESC")
	}

	result := query.Limit(limit).Offset(offset).Find(&announcements)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPage := (int(totalData) + limit - 1) / limit

	return announcements, totalPage, int(totalData), nil
}

func (repositories *announcementRepository) GetAnnouncementByID(id int) (response.AnnouncementResponse, error) {
	var announcement response.AnnouncementResponse

	query := database.DB.Table("announcements as a").
		Select("a.id, a.title, a.description, case when a.hero_image = '' then null else a.hero_image end, a.type, a.event_date, a.created_at, a.created_by, a.updated_at, a.updated_by").
		Where("a.id = ?", id).
		First(&announcement)

	if query.Error != nil {
		return response.AnnouncementResponse{}, query.Error
	}

	return announcement, nil
}
