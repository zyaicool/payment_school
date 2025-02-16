package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	database "schoolPayment/configs"
	"schoolPayment/constants"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUser(page int, limit int, search string, roleID []int, schoolID int, sortBy string, sortOrder string, status *bool) ([]models.User, int, int64, error)
	GetUserByIDPass(id uint) (*models.User, error)
	UpdateUserRepository(user *models.User) (*models.User, error)
	DeleteUserRepository(user *models.User) (*models.User, error)
	GetEmailSendCount(email string, now time.Time) (int, error)
	RecordEmailSend(email string) error
	GetUserByID(id uint) (models.User, error)
	DeleteAllTempVerificationEmails(userID int) error
	GetEmailVerification(userID int, typeFilter string) (models.TempVerificationEmail, error)
	UpdateTempVerificationEmail(verificationEmail *models.TempVerificationEmail) error
	BulkValidateAndCreateUsers(users []models.User, userSchools []models.UserSchool) ([]models.User, []response.ResponseErrorUploadUser, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (userRepository *userRepository) GetAllUser(page int, limit int, search string, roleID []int, schoolID int, sortBy string, sortOrder string, status *bool) ([]models.User, int, int64, error) {
	var users []models.User
	var total int64
	totalPages := 0

	// Base query without pagination
	query := database.DB.Model(&models.User{})
	query = query.Joins("JOIN roles on roles.id = users.role_id")

	query = query.Where("users.deleted_at IS NULL")

	// Add school filter if schoolID is not 0
	if schoolID != 0 {
		query = query.Joins("JOIN user_schools ON user_schools.user_id = users.id AND user_schools.deleted_at IS NULL").
			Joins("JOIN schools ON schools.id = user_schools.school_id AND schools.deleted_at IS NULL").
			Where("schools.id = ?", schoolID)
	}

	// Add search filter if provided
	if search != "" {
		query = query.Where("LOWER(users.username) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if sortBy != "" && sortOrder != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		query = query.Order("CASE WHEN users.updated_at IS NOT NULL THEN 0 ELSE 1 END, users.updated_at DESC, users.created_at DESC")
	}

	if status != nil {
		query = query.Where("users.is_block != ?", *status)
	}

	// Add role filter if roleID is not 0
	if len(roleID) > 0 {
		query = query.Where("users.role_id IN ?", roleID)
	}

	// Count total records without pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if limit != 0 {
		// Calculate total pages
		totalPages = int(total / int64(limit))
		if total%int64(limit) != 0 {
			totalPages++
		}

		// Add pagination (offset and limit) for the actual query
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)

	}

	// Execute the query to get paginated users
	result := query.Preload("Role").Preload("UserSchool").Preload(constants.PreloadUserSchoolToSchool).Find(&users)

	// Return users, total pages, total records, and error
	return users, totalPages, total, result.Error
}

func (userRepository *userRepository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	result := database.DB.Where("id = ? AND deleted_at IS NULL ", id).
		Preload("Role").Preload(constants.PreloadRoleToRoleMatrix).Preload("UserSchool").Preload(constants.PreloadUserSchoolToSchool).Preload("UserSchool.School.SchoolGrade").Preload("UserStudents").First(&user)
	return user, result.Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	result := database.DB.Where("email = ? AND deleted_at IS NULL", email).
		Preload("Role").Preload(constants.PreloadRoleToRoleMatrix).Preload("UserSchool").Preload(constants.PreloadUserSchoolToSchool).First(&user)
	return user, result.Error
}

func GetUserByEmailIsBlock(email string) (models.User, error) {
	var user models.User

	result := database.DB.Where("email = ? AND deleted_at IS NULL AND is_block = TRUE", email).
		Preload("Role").Preload(constants.PreloadRoleToRoleMatrix).Preload("UserSchool").Preload(constants.PreloadUserSchoolToSchool).First(&user)
	return user, result.Error
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	result := database.DB.Where("username = ? AND deleted_at IS NULL", username).First(&user)
	return user, result.Error
}

func CreateUser(user *models.User) (*models.User, error) {
	result := database.DB.Create(&user)
	return user, result.Error
}

func CreateVerificationLink(verif *models.TempVerificationEmail) (*models.TempVerificationEmail, error) {
	result := database.DB.Create(&verif)
	return verif, result.Error
}

func BulkCreateVerificationLinks(verifications []*models.TempVerificationEmail) error {
	if len(verifications) == 0 {
		return nil
	}

	// Using gorm's CreateInBatches for better performance with large datasets
	result := database.DB.CreateInBatches(verifications, 500)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk create verification links: %v", result.Error)
	}

	return nil
}

func (userRepository *userRepository) GetUserByIDPass(id uint) (*models.User, error) {
	var user models.User
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil // Pastikan ini mengembalikan pointer
}

func ValidateDuplicateUsername(username string) bool {
	var user models.User
	result := database.DB.Where("username = ? AND deleted_at IS NULL", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}

	return false
}

func (userRepository *userRepository) UpdateUserRepository(user *models.User) (*models.User, error) {
	// Update password user berdasarkan ID
	result := database.DB.Model(&user).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username": user.Username,
		"role_id":  user.RoleID,
		"is_block": user.IsBlock,
		"password": user.Password, // Assuming password update is still needed
		"image":    user.Image,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (userRepository *userRepository) DeleteUserRepository(user *models.User) (*models.User, error) {
	// this function can use to update data user, beside naming of function is delete user
	result := database.DB.Save(&user)
	return user, result.Error
}

func (userRepository *userRepository) GetEmailSendCount(email string, now time.Time) (int, error) {
	var emailSendRecord models.EmailSendRecord
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var count int64
	result := database.DB.Model(emailSendRecord).Where("email = ? AND timestamp BETWEEN ? AND ?", email, startOfDay, endOfDay).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func (userRepository *userRepository) RecordEmailSend(email string) error {
	record := models.EmailSendRecord{
		Email:     email,
		Timestamp: time.Now(),
	}
	result := database.DB.Create(&record)
	return result.Error
}

func GetUserByStudentId(studentID uint) (string, error) {
	var user models.User
	result := database.DB.
		Joins("JOIN user_students ON user_students.user_id = users.id").
		Where("user_students.student_id = ? AND users.deleted_at IS NULL ", studentID).First(&user)
	return user.Email, result.Error
}

func GetUserByID2(id uint) (models.User, error) {
	var user models.User
	result := database.DB.Select("id, Username, email, role_id").Where("id = ? AND deleted_at IS NULL ", id).Preload("UserSchool").Preload(constants.PreloadUserSchoolToSchool).First(&user)
	return user, result.Error
}

func (userRepository *userRepository) DeleteAllTempVerificationEmails(userID int) error {
	result := database.DB.Table("temp_verification_emails").
		Where("user_id = ?", userID).
		Delete(nil)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userRepository *userRepository) GetEmailVerification(userID int, typeFilter string) (models.TempVerificationEmail, error) {
	var emailVerification models.TempVerificationEmail

	result := database.DB.Table("temp_verification_emails").
		Where("user_id = ? AND deleted_at IS NULL ", userID)

	if typeFilter != "" {
		result = result.Where("type = ? ", typeFilter)
	} else {
		result = result.Where("type is null ")
	}

	result = result.Order("created_at desc").First(&emailVerification)

	return emailVerification, result.Error
}

func (userRepository *userRepository) UpdateTempVerificationEmail(verificationEmail *models.TempVerificationEmail) error {
	// Update password user berdasarkan ID
	result := database.DB.Save(&verificationEmail)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// userRepository untuk mendapatkan data school berdasarkan email pengguna.
func GetSchoolByEmail(email string) (*models.School, error) {
	var school models.School
	err := database.DB.
		Table("schools").
		Select("schools.school_logo").
		Joins("JOIN users ON users.school_id = schools.id").
		Where("users.email = ?", email).
		First(&school).Error

	if err != nil {
		return nil, err
	}

	return &school, nil
}

func GetValidateVerifyToken(token string) (models.TempVerificationEmail, error) {
	var emailVerification models.TempVerificationEmail

	result := database.DB.Table("temp_verification_emails").
		Where("verification_link = ? AND deleted_at IS NULL AND is_valid IS TRUE", token).
		First(&emailVerification)

	return emailVerification, result.Error
}

func (userRepository *userRepository) BulkValidateAndCreateUsers(users []models.User, userSchools []models.UserSchool) ([]models.User, []response.ResponseErrorUploadUser, error) {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check for duplicate emails and usernames
	var existingUsers []models.User
	err := tx.Where("email IN ? OR username IN ? AND deleted_at IS NULL",
		getUserEmails(users),
		getUserUsernames(users)).
		Find(&existingUsers).Error

	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	var errorResponses []response.ResponseErrorUploadUser
	var validUsers []models.User
	var validUserSchools []models.UserSchool

	// Filter out duplicate users and create error responses
	for i, user := range users {
		isDuplicate := false
		for _, existingUser := range existingUsers {
			if user.Email == existingUser.Email || user.Username == existingUser.Username {
				isDuplicate = true
				errorResponses = append(errorResponses, response.ResponseErrorUploadUser{
					Username: user.Username,
					Email:    user.Email,
					Reason:   "Username/Email Sudah Terpakai",
				})
				break
			}
		}
		if !isDuplicate {
			validUsers = append(validUsers, user)
			if len(userSchools) > i {
				validUserSchools = append(validUserSchools, userSchools[i])
			}
		}
	}

	// If no valid users to create, return early with error responses
	if len(validUsers) == 0 {
		tx.Rollback()
		return nil, errorResponses, nil
	}

	// Bulk insert valid users
	if err := tx.CreateInBatches(validUsers, 500).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// After users are created, get their IDs from the database
	for i := range validUsers {
		validUserSchools[i].UserID = validUsers[i].ID
	}

	// Bulk insert user schools for valid users
	if len(validUserSchools) > 0 {
		if err := tx.CreateInBatches(validUserSchools, 500).Error; err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	// Before returning, preload the UserSchool relationship for all valid users
	var usersWithSchools []models.User
	if err := tx.Where("id IN ?", getUserIDs(validUsers)).
		Preload("UserSchool").
		Preload("UserSchool.School").
		Find(&usersWithSchools).Error; err != nil {
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return usersWithSchools, errorResponses, nil
}

func getUserEmails(users []models.User) []string {
	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}
	return emails
}

func getUserUsernames(users []models.User) []string {
	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}
	return usernames
}

func getUserIDs(users []models.User) []uint {
	ids := make([]uint, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}

func GetUsersByEmails(emails []string) ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("email IN ? AND deleted_at IS NULL", emails).
		Preload("UserSchool").
		Preload("UserSchool.School").
		Find(&users).Error
	return users, err
}

func GetUserByEmailAndFilterSchoolId(email string, schoolId uint) (int, bool, error) {
	var rsp response.GetUserByEmailAndSchoolIdDto
	valid := false
	query := database.DB.Table("users as u").
		Select(`
			u.id as user_id,
			case when us.school_id = null then 0 else us.school_id end as school_id
		`).Joins("left join user_schools us on us.user_id = u.id")

	query = query.Where("u.email = ? AND us.school_id = ?  AND u.deleted_at IS NULL", email, schoolId)
	err := query.Find(&rsp).Error
	if err != nil {
		return 0, true, err
	}

	if rsp.UserId == 0 && rsp.SchoolId == 0 {
		return 0, true, fmt.Errorf("no record") // No record found
	}

	if rsp.SchoolId != 0 {
		valid = true
	}

	return rsp.UserId, valid, nil
}
