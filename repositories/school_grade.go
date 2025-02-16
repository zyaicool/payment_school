package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type SchoolGradeRepositoryInterface interface {
	GetSchoolGradeByID(id uint) (models.SchoolGrade, error)
	CreateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error)
	UpdateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error)
	GetLastSequenceNumberSchoolGrades() (int, error)
}

type SchoolGradeRepository struct{
	db *gorm.DB
}

func NewSchoolGradeRepository(db *gorm.DB) SchoolGradeRepositoryInterface {
	return &SchoolGradeRepository{db: db}
}

func GetAllSchoolGrade(page int, limit int, search string) ([]models.SchoolGrade, error) {
	var schoolGradeList []models.SchoolGrade
	query := database.DB.Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("LOWER(school_grade_name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if limit != 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&schoolGradeList)
	return schoolGradeList, result.Error
}

func (schoolGradeRepository *SchoolGradeRepository) GetSchoolGradeByID(id uint) (models.SchoolGrade, error) {
	var schoolGrade models.SchoolGrade
	result := schoolGradeRepository.db.Where("id = ? AND deleted_at IS NULL", id).First(&schoolGrade)
	return schoolGrade, result.Error
}

func (schoolGradeRepository *SchoolGradeRepository) CreateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error) {
	result := database.DB.Create(&schoolGrade)
	return schoolGrade, result.Error
}

func (schoolGradeRepository *SchoolGradeRepository) UpdateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error) {
	result := database.DB.Save(&schoolGrade)
	return schoolGrade, result.Error
}

func (schoolGradeRepository *SchoolGradeRepository) GetLastSequenceNumberSchoolGrades() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("school_grades").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

func GetAllSchoolGradesBulk(gradeNames []string) ([]models.SchoolGrade, error) {
	var grades []models.SchoolGrade
	result := database.DB.Where("school_grade_name IN ?", gradeNames).Find(&grades)
	return grades, result.Error
}
