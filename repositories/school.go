package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type SchoolRepository interface {
	GetSchoolByID(id uint) (models.School, error)
	CheckNpsn(npsn uint) (models.School, error)
	CreateSchool(school *models.School) (*models.School, error)
	UpdateSchool(school *models.School) (*models.School, error)
	GetLastSequenceNumberSchool() (int, error)
	CheckNpsnExistsExcept(npsn uint, schoolID int) (models.School, error)
	GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) ([]models.SchoolList, int, error)
	GetAllOnboardingSchools(search string) ([]models.School, error)
	GetSchoolsByNames(schoolNames []string) (map[string]uint, error)
}

type schoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

func (schoolRepository *schoolRepository) GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) ([]models.SchoolList, int, error) {
	var schoolList []models.SchoolList
	var count int64

	// Join users table to get user information
	query := schoolRepository.db.Model(&models.School{}).
		Select("schools.*, case when users.username is null then 'system' else users.username end AS created_by_username").
		Joins("LEFT JOIN users ON users.id = schools.created_by").
		Where("schools.deleted_at IS NULL")

	if search != "" {
		query = query.Where("lower(schools.school_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	// Count total records
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and order by created_at DESC
	if limit != 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("CASE WHEN schools.updated_at IS NOT NULL THEN 0 ELSE 1 END, schools.updated_at DESC, schools.created_at DESC")
	}

	// Execute query to fetch school list
	result := query.Find(&schoolList)
	return schoolList, int(count), result.Error
}

func GetAllSchool(page int, limit int, search string) ([]models.School, error) {
	var schoolList []models.School
	query := database.DB.Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("lower(school_name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if limit != 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&schoolList)
	return schoolList, result.Error
}

func (schoolRepository *schoolRepository) GetSchoolByID(id uint) (models.School, error) {
	var school models.School
	// result := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&school)
	result := schoolRepository.db.Where("id = ? AND deleted_at IS NULL", id).Preload("SchoolGrade").First(&school)
	return school, result.Error
}

func (schoolRepository *schoolRepository) CheckNpsn(npsn uint) (models.School, error) {
	var school models.School
	result := database.DB.Where("npsn = ? AND deleted_at IS NULL", npsn).First(&school)
	return school, result.Error
}

func (schoolRepository *schoolRepository) CheckNpsnExistsExcept(npsn uint, schoolID int) (models.School, error) {
	var school models.School
	result := database.DB.Where("npsn = ? AND id != ? AND deleted_at IS NULL", npsn, schoolID).First(&school)
	return school, result.Error
}

func (schoolRepository *schoolRepository) CreateSchool(school *models.School) (*models.School, error) {
	result := schoolRepository.db.Create(&school)
	return school, result.Error
}

func (schoolRepository *schoolRepository) UpdateSchool(school *models.School) (*models.School, error) {
	result := schoolRepository.db.Save(&school)
	return school, result.Error
}

func (schoolRepository *schoolRepository) GetLastSequenceNumberSchool() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("schools").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

func GetEmailParentById(studentId int) (models.User, error) {
	var user models.User
	query := database.DB.Joins("JOIN user_students us on us.user_id = users.id").
		Where("us.student_id = ?", studentId)

	result := query.First(&user)
	return user, result.Error
}

func (schoolRepository *schoolRepository) GetAllOnboardingSchools(search string) ([]models.School, error) {
	var schools []models.School

	query := `SELECT 
		id, 
		school_name, 
		npsn, 
		school_logo
	FROM schools 
	WHERE school_name LIKE ?
	`
	searchTerm := "%" + search + "%"
	err := database.DB.Raw(query, searchTerm).Scan(&schools).Error
	if err != nil {
		return nil, err
	}

	// err := database.DB.Raw(query).Scan(&schools).Error
	// if err != nil {
	// 	return nil, err
	// }

	return schools, nil
}

func (schoolRepository *schoolRepository) GetSchoolsByNames(schoolNames []string) (map[string]uint, error) {
	var schools []models.School
	schoolMap := make(map[string]uint)

	err := schoolRepository.db.Where("school_name IN ? AND deleted_at IS NULL", schoolNames).Find(&schools).Error
	if err != nil {
		return nil, err
	}

	for _, school := range schools {
		schoolMap[school.SchoolName] = school.ID
	}

	return schoolMap, nil
}
