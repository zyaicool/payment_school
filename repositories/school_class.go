package repositories

import (
	"strings"

	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
)

type SchoolClassRepositoryInterface interface {
	CreateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error)
	GetLastSequenceNumberSchoolClasss() (int, error)
	GetSchoolClassByID(id uint) (models.SchoolClass, error)
	UpdateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error)
}

type SchoolClassRepository struct{}

func NewSchoolClassRepository() SchoolClassRepositoryInterface {
	return &SchoolClassRepository{}
}

func GetAllSchoolClass(page int, limit int, search string, sortBy string, sortOrder string, showDeletedData bool, user models.User) ([]response.SchoolClassResponseRepo, int, int64, error) {
	var schoolClassList []response.SchoolClassResponseRepo
	var total int64
	totalPages := 0

	query := database.DB.Table("school_classes a").
		Select(`
			a.id,
			b.school_grade_name AS unit,
			a.school_class_name,
			a.created_by AS created_by,
			(CASE WHEN a.deleted_at IS NULL THEN true ELSE false END) AS status,
			CONCAT('Kelas - ', a.school_class_name) AS placeholder,
			sm.school_major_name AS school_major,
			d.prefix_name AS prefix_class,
			u.username AS created_by_username
		`).
		Joins("JOIN school_grades b ON a.school_grade_id = b.id").
		Joins("JOIN users u ON a.created_by = u.id").
		Joins("LEFT JOIN prefix_classes d ON a.prefix_class_id = d.id").
		Joins("LEFT JOIN school_majors sm ON a.school_major_id = sm.id")
		// Where("a.deleted_at IS NULL")
	if user.UserSchool != nil {
		query = query.Where("a.school_id = ?", user.UserSchool.School.ID)
	}

	if !showDeletedData {
		query = query.Where("a.deleted_at IS NULL")
	}

	if search != "" {
		query = query.Where("LOWER(a.school_class_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if limit != 0 {
		// Calculate total pages
		totalPages = int(total / int64(limit))
		if total%int64(limit) != 0 {
			totalPages++
		}
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	// Mengambil data
	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("CASE WHEN a.updated_at IS NOT NULL THEN 0 ELSE 1 END, a.updated_at DESC, a.created_at DESC")
	}
	err := query.Scan(&schoolClassList).Error
	return schoolClassList, totalPages, total, err
}

func (schoolClassRepository *SchoolClassRepository) GetSchoolClassByID(id uint) (models.SchoolClass, error) {
	var schoolClass models.SchoolClass
	result := database.DB.Where("id = ?", id).First(&schoolClass)
	return schoolClass, result.Error
}

func (schoolClassRepository *SchoolClassRepository) CreateSchoolClass(SchoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	result := database.DB.Create(&SchoolClass)
	return SchoolClass, result.Error
}

func (schoolClassRepository *SchoolClassRepository) UpdateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	result := database.DB.Save(&schoolClass)
	return schoolClass, result.Error
}

func (schoolClassRepository *SchoolClassRepository) GetLastSequenceNumberSchoolClasss() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("school_classes").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

func GetAllSchoolClassesBulk(classNames []string, user models.User) ([]models.SchoolClass, error) {
	var classes []models.SchoolClass
	result := database.DB.Where("school_class_name IN ? AND school_id = ?", classNames, user.UserSchool.SchoolID).Find(&classes)
	return classes, result.Error
}
