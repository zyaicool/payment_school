package repositories

import (
	"fmt"
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/constants"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"

	"gorm.io/gorm"
)

type StudentRepositoryInteface interface {
	GetAllStudentForBilling(user models.User, schoolGradeId int, schoolClassIds []int) ([]models.Student, error)
	CreateStudent(student *models.Student) (*models.Student, error)
	CreateStudentHistory(studentHistory *models.StudentHistory) (*models.StudentHistory, error)
	UpdateStudent(student *models.Student) (*models.Student, error)
	CreateUserStudentRepository(userStudent *models.UserStudent) (*models.UserStudent, error)
	GetStudentByUserIdRepository(id uint) ([]models.Student, error)
	BulkCreateStudents(students []models.Student) error
	BulkUpdateStudents(students []models.Student) error
	GetStudentByNis(nis string, user models.User) (models.Student, error)
	GetStudentsByNIS(nisNumbers []string, user models.User) ([]models.Student, error)
	BulkCheckUserStudentExists(pairs []struct {
		UserID    uint
		StudentID uint
	}) (map[string]bool, error)
	BulkCreateUserStudents(pairs []models.UserStudent) error
	BulkCreateStudentHistory(histories []models.StudentHistory) error
	CreateImageStudentRepository(student *models.Student, id int) (*models.Student, error)
	GetStudentByID(id uint, user models.User) (models.Student, error)
}

type StudentRepository struct{
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepositoryInteface {
	return &StudentRepository{db: db}
}

func GetAllStudent(page int, limit int, search string, user models.User, status string, gradeID int, yearId int, schoolID int, searchNis string, classID int, sortBy string, sortOrder string, studentId int, isActive *bool) ([]response.DetailStudentResponse, int, int64, error) {
	var students []response.DetailStudentResponse
	var total int64
	totalPages := 0

	// Base query without pagination
	// query := database.DB.Model(&models.Student{}).Where("students.deleted_at IS NULL")
	query := database.DB.Table("students as s")
	query = query.Select("s.*, sg.school_grade_name as school_grade, sc.school_class_name as school_class, sy.school_year_name, concat(s.nis, ' - ', s.full_name, ', ', sc.school_class_name, ', ', sg.school_grade_name) as placeholder ")
	query = query.Joins("JOIN school_grades sg ON sg.id = s.school_grade_id")
	query = query.Joins("JOIN school_classes sc ON sc.id = s.school_class_id")
	query = query.Joins("JOIN school_years sy ON sy.id = s.school_year_id")
	query = query.Where("s.deleted_at IS NULL")

	// Search filter
	if search != "" {
		query = query.Where("LOWER(s.full_name) like LOWER(?)", "%"+search+"%")
	}

	if searchNis != "" {
		query = query.Where("s.nis like ?", "%"+searchNis+"%")
	}

	if studentId != 0 {
		query = query.Where("s.id = ?", studentId)
	}

	// Role-based filters
	if user.RoleID == 2 {
		query = query.Joins(constants.JoinUserStudentsToStudents).
			Joins(constants.JoinUsersToUserStudents).
			Where(constants.FilterByUsersId, user.ID)
	} else if user.RoleID == 5 || user.RoleID == 4 || user.RoleID == 3 {
		query = query.Joins(constants.JoinUserStudentsToStudents).
			Joins(constants.JoinUsersToUserStudents).
			Joins("JOIN user_schools ON user_schools.user_id = users.id ").
			Where("user_schools.school_id = ?", user.UserSchool.School.ID)
	} else if user.RoleID == 1 && schoolID != 0 {
		query = query.Joins(constants.JoinUserStudentsToStudents).
			Joins(constants.JoinUsersToUserStudents).
			Joins("JOIN user_schools ON user_schools.user_id = users.id ").
			Where("user_schools.school_id = ?", schoolID)
	}

	// Additional filters
	if status != "" {
		query = query.Where("LOWER(s.status) = ?", strings.ToLower(status))
	}

	if isActive != nil {
		if *isActive {
			//query = query.Where("s.status = ?", "Aktif")
			query = query.Where("LOWER(s.status) IN (?)", []string{"aktif", "Aktif"})
		} else {
			//query = query.Where("s.status != ?", "Aktif")
			query = query.Where("LOWER(s.status) NOT IN (?)", []string{"aktif", "Aktif"})
		}
	}

	if gradeID != 0 {
		query = query.Where("s.school_grade_id = ?", gradeID)
	}

	if yearId != 0 {
		query = query.Where("s.school_year_id = ?", yearId)
	}

	if classID != 0 {
		query = query.Where("s.school_class_id = ?", classID)
	}

	// Get total count of records (without pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Apply pagination (Offset and Limit)
	if limit != 0 {
		// Calculate total pages
		totalPages = int(total / int64(limit))
		if total%int64(limit) != 0 {
			totalPages++
		}
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if sortBy == "" {
		query = query.Order("CASE WHEN s.updated_at IS NOT NULL THEN 0 ELSE 1 END, s.updated_at DESC, s.created_at DESC")
	} else {
		if sortOrder == "" {
			sortOrder = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}

	// Fetch students
	result := query.Find(&students)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	// Return students, total pages, and total records
	return students, totalPages, total, nil
}

func (studentRepository *StudentRepository) GetStudentByID(id uint, user models.User) (models.Student, error) {
	var student models.Student
	query := database.DB.Where("students.id = ? AND students.deleted_at IS NULL", id)
	if user.RoleID == 5 && user.UserSchool != nil || user.RoleID == 4 && user.UserSchool != nil {
		query = query.Joins(constants.JoinUserStudentsToStudentsAndFilterDeletedAt).
			Joins(constants.JoinUsersToUserStudentsAndFilterDeletedAt).
			Joins(constants.JoinUserSchoolsToUsersAndFilterDeletedAt).
			Joins(constants.JoinSchoolsToUserSChoolsAndFilterDeletedAt).
			Where("schools.id = ?", user.UserSchool.School.ID)
	} else if user.RoleID == 2 && user.UserSchool != nil {
		query = query.Joins(constants.JoinUserStudentsToStudentsAndFilterDeletedAt).
			Joins(constants.JoinUsersToUserStudentsAndFilterDeletedAt).
			Where(constants.FilterByUsersId, user.ID)
	}
	result := query.First(&student)
	return student, result.Error
}

func (studentRepository *StudentRepository) CreateStudent(student *models.Student) (*models.Student, error) {
	result := database.DB.Create(&student)
	return student, result.Error
}

func (studentRepository *StudentRepository) UpdateStudent(student *models.Student) (*models.Student, error) {
	result := database.DB.Save(student)

	return student, result.Error
}

func (studentRepository *StudentRepository) CreateUserStudentRepository(userStudent *models.UserStudent) (*models.UserStudent, error) {
	result := database.DB.Save(userStudent)

	return userStudent, result.Error
}

func (studentRepository *StudentRepository) GetStudentByUserIdRepository(id uint) ([]models.Student, error) {
	var student []models.Student
	result := database.DB.Joins("JOIN user_students ON user_students.student_id = students.id").
		Where("user_students.user_id = ?", id).Find(&student)
	return student, result.Error
}

func (studentRepository *StudentRepository) CreateStudentHistory(studentHistory *models.StudentHistory) (*models.StudentHistory, error) {
	result := database.DB.Save(studentHistory)

	return studentHistory, result.Error
}

func (studentRepository *StudentRepository) GetAllStudentForBilling(user models.User, schoolGradeId int, schoolClassIds []int) ([]models.Student, error) {
	var students []models.Student

	query := studentRepository.db.Joins(constants.JoinUserStudentsToStudentsAndFilterDeletedAt).
		Joins(constants.JoinUsersToUserStudentsAndFilterDeletedAt).
		Joins(constants.JoinUserSchoolsToUsersAndFilterDeletedAt).
		Joins(constants.JoinSchoolsToUserSChoolsAndFilterDeletedAt).
		Where("students.deleted_at IS NULL AND schools.id = ? AND students.school_grade_id = ? AND LOWER(students.status) = ?",
			user.UserSchool.School.ID,
			schoolGradeId,
			"aktif")

	if len(schoolClassIds) > 0 {
		query = query.Where("students.school_class_id IN ?", schoolClassIds)
	}

	result := query.Find(&students)
	return students, result.Error
}

func (studentRepository *StudentRepository) GetStudentByNis(nis string, user models.User) (models.Student, error) {
	var student models.Student
	query := database.DB.Where("students.deleted_at is null and LOWER(students.nis) like ?", strings.ToLower(nis))
	if user.RoleID == 5 || user.RoleID == 4 {
		query = query.Joins(constants.JoinUserStudentsToStudentsAndFilterDeletedAt).
			Joins(constants.JoinUsersToUserStudentsAndFilterDeletedAt).
			Joins(constants.JoinUserSchoolsToUsersAndFilterDeletedAt).
			Joins(constants.JoinSchoolsToUserSChoolsAndFilterDeletedAt).
			Where("schools.id = ?", user.UserSchool.School.ID)
	} else if user.RoleID == 2 {
		query = query.Joins(constants.JoinUserStudentsToStudentsAndFilterDeletedAt).
			Joins(constants.JoinUsersToUserStudentsAndFilterDeletedAt).
			Where(constants.FilterByUsersId, user.ID)
	}
	result := query.First(&student)
	return student, result.Error
}

func (studentRepository *StudentRepository) CreateImageStudentRepository(student *models.Student, id int) (*models.Student, error) {
	result := database.DB.Model(student).Where("id = ?", id).Update("image", student.Image)
	if result.Error != nil {
		return nil, result.Error
	}

	return student, nil
}

func GetSchoolByStudentId(studentID uint) (*models.School, error) {
	var school *models.School
	query := database.DB.Joins("JOIN user_schools us on us.school_id = schools.id").
		Joins("JOIN user_students st on st.user_id = us.user_id").
		Where("st.student_id = ?", studentID)
	result := query.First(&school)
	return school, result.Error
}

func GetStudentByIDOnlyStudent(studentID uint) (*models.Student, error) {
	var student *models.Student

	result := database.DB.Model(&models.Student{}).Where("id = ?", studentID).First(&student)

	return student, result.Error
}

func GetStudentByUserIdDashboard(id uint) ([]models.Student, error) {
	var student []models.Student
	result := database.DB.Joins("JOIN user_students ON user_students.student_id = students.id").
		Where("user_students.user_id = ?", id).Find(&student)
	return student, result.Error
}

func GetSchoolGradeByUser(db *gorm.DB, userID uint) ([]models.SchoolGrade, error) {
	var schoolGrades []models.SchoolGrade

	err := db.Joins("JOIN schools s ON s.school_grade_id = school_grades.id").
		Joins("JOIN user_schools us ON us.school_id = s.id").
		Joins("JOIN school_grades sg ON sg.id = s.school_grade_id").
		Where("us.user_id = ?", userID).
		Select("school_grades.school_grade_name").
		Find(&schoolGrades).Error

	if err != nil {
		return nil, err
	}

	return schoolGrades, nil
}

func GetSchoolByUser(db *gorm.DB, userID uint) ([]models.School, error) {
	var schools []models.School
	err := db.Joins("JOIN user_schools us ON us.school_id = schools.id").
		Where("us.user_id = ?", userID).
		Select("schools.school_name").
		Find(&schools).Error
	if err != nil {
		return nil, err
	}
	return schools, nil
}

func (r *StudentRepository) BulkCreateStudents(students []models.Student) error {
	return database.DB.CreateInBatches(students, 500).Error
}

func (r *StudentRepository) BulkUpdateStudents(students []models.Student) error {
	for i := 0; i < len(students); i += 500 {
		end := i + 500
		if end > len(students) {
			end = len(students)
		}

		batch := students[i:end]
		if err := database.DB.Save(&batch).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *StudentRepository) GetStudentsByNIS(nisNumbers []string, user models.User) ([]models.Student, error) {
	var students []models.Student
	query := database.DB.Where("nis IN ?", nisNumbers)

	if user.RoleID != 1 {
		query = query.Joins("JOIN user_students us ON us.student_id = students.id").
			Joins("JOIN user_schools usc ON usc.user_id = us.user_id").
			Where("usc.school_id = ?", user.UserSchool.SchoolID)
	}

	err := query.Find(&students).Error
	return students, err
}

func (studentRepository *StudentRepository) BulkCreateStudentHistory(histories []models.StudentHistory) error {
	return database.DB.CreateInBatches(histories, 500).Error
}

// Add this function to studentRepository
func (r *StudentRepository) BulkCreateUserStudents(pairs []models.UserStudent) error {
	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create maps for batch checking
	type userStudentKey struct {
		UserID    uint
		StudentID uint
	}
	pairsMap := make(map[userStudentKey]models.UserStudent)
	var userIDs []uint
	var studentIDs []uint

	// Collect all IDs for batch checking
	for _, pair := range pairs {
		key := userStudentKey{UserID: pair.UserID, StudentID: pair.StudentID}
		pairsMap[key] = pair
		userIDs = append(userIDs, pair.UserID)
		studentIDs = append(studentIDs, pair.StudentID)
	}

	// Batch check for existing relationships
	var existingPairs []models.UserStudent
	if err := tx.Where("user_id IN ? AND student_id IN ? AND deleted_at IS NULL",
		userIDs, studentIDs).Find(&existingPairs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check existing relationships: %v", err)
	}

	// Remove existing pairs from the map
	for _, existing := range existingPairs {
		key := userStudentKey{UserID: existing.UserID, StudentID: existing.StudentID}
		delete(pairsMap, key)
	}

	// If no new pairs to create, return early
	if len(pairsMap) == 0 {
		tx.Commit()
		return nil
	}

	// Convert map back to slice for batch insert
	var pairsToCreate []models.UserStudent
	for _, pair := range pairsMap {
		pairsToCreate = append(pairsToCreate, pair)
	}

	// Batch insert new relationships
	const batchSize = 500
	for i := 0; i < len(pairsToCreate); i += batchSize {
		end := i + batchSize
		if end > len(pairsToCreate) {
			end = len(pairsToCreate)
		}

		if err := tx.CreateInBatches(pairsToCreate[i:end], batchSize).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create user-student relationships batch: %v", err)
		}
	}

	return tx.Commit().Error
}

// In student repository
func (r *StudentRepository) CheckUserStudentExists(userID, studentID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.UserStudent{}).
		Where("user_id = ? AND student_id = ? AND deleted_at IS NULL", userID, studentID).
		Count(&count).Error
	return count > 0, err
}

func (studentRepository *StudentRepository) BulkCheckUserStudentExists(pairs []struct {
	UserID    uint
	StudentID uint
}) (map[string]bool, error) {
	results := make(map[string]bool)
	var existingRelations []models.UserStudent

	// Create slices to hold all userIDs and studentIDs
	var userIDs []uint
	var studentIDs []uint
	for _, pair := range pairs {
		userIDs = append(userIDs, pair.UserID)
		studentIDs = append(studentIDs, pair.StudentID)
	}

	// Single query to get all existing relationships
	err := database.DB.Where("user_id IN ? AND student_id IN ? AND deleted_at IS NULL",
		userIDs, studentIDs).Find(&existingRelations).Error
	if err != nil {
		return nil, err
	}

	// Create map of existing relationships for quick lookup
	existingMap := make(map[string]bool)
	for _, relation := range existingRelations {
		key := fmt.Sprintf("%d-%d", relation.UserID, relation.StudentID)
		existingMap[key] = true
	}

	// Check each pair against the map
	for _, pair := range pairs {
		key := fmt.Sprintf("%d-%d", pair.UserID, pair.StudentID)
		results[key] = existingMap[key]
	}

	return results, nil
}
