package services

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"schoolPayment/configs"
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	userServiceInstance *UserService
	userServiceOnce     sync.Once
)

func GetUserService() UserService {
	userServiceOnce.Do(func() {
		userRepository := repositories.NewUserRepository()
		roleRepository := repositories.NewRoleRepository()
		schoolRepository := repositories.NewSchoolRepository(configs.DB)

		userServiceInstance = &UserService{
			userRepository:   userRepository,
			roleRepository:   roleRepository,
			schoolRepository: schoolRepository,
		}
	})
	return *userServiceInstance
}

type UserService struct {
	userRepository   repositories.UserRepository
	roleRepository   repositories.RoleRepository
	schoolRepository repositories.SchoolRepository
}

func NewUserService(userRepository repositories.UserRepository, roleRepository repositories.RoleRepository,
	schoolRepository repositories.SchoolRepository) UserService {
	return UserService{
		userRepository:   userRepository,
		roleRepository:   roleRepository,
		schoolRepository: schoolRepository,
	}
}

func (userService *UserService) GetAllUser(page int, limit int, search string, roleID []int, userID int, schoolID int, sortBy string, sortOrder string, status *bool) (response.UserListResponse, error) {
	var resp response.UserListResponse
	var listUserResponse []response.ListDataUserForIndex
	var schoolIDRes uint
	var schoolName string

	// Initialize response
	resp.Limit = limit
	resp.Page = page
	resp.TotalData = 0
	resp.TotalPage = 0

	// Get user information to determine schoolID
	getUser, err := userService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		resp.Data = []response.ListDataUserForIndex{}
		return resp, nil
	}

	if getUser.UserSchool != nil {
		schoolID = int(getUser.UserSchool.School.ID)
	}

	// Adjust sorting field
	if sortBy != "" {
		sortBy = utilities.ChangeStringSortByUser(sortBy)
	}

	// Call repository to fetch users
	dataUser, totalPage, totalData, err := userService.userRepository.GetAllUser(page, limit, search, roleID, schoolID, sortBy, sortOrder, status)
	if err != nil {
		resp.Data = []response.ListDataUserForIndex{}
		return resp, nil
	}

	// Process fetched data
	for _, user := range dataUser {

		if user.UserSchool != nil {
			schoolIDRes = user.UserSchool.School.ID
			schoolName = user.UserSchool.School.SchoolName
		}

		responseDataUser := response.ListDataUserForIndex{
			ID:                int(user.ID),
			RoleID:            int(user.RoleID),
			RoleName:          user.Role.Name,
			Username:          user.Username,
			SchoolID:          int(schoolIDRes),
			SchooolName:       schoolName,
			Email:             user.Email,
			CreatedDate:       user.CreatedAt,
			Status:            mapStatus(user.IsBlock),
			VerificationEmail: user.IsVerification,
		}

		listUserResponse = append(listUserResponse, responseDataUser)
	}

	// Update response
	if len(listUserResponse) > 0 {
		resp.Data = listUserResponse
	} else {
		resp.Data = []response.ListDataUserForIndex{}
	}

	resp.TotalData = totalData
	resp.TotalPage = totalPage

	return resp, nil
}

func (userService *UserService) GetUserByID(id uint) (models.User, error) {
	return userService.userRepository.GetUserByID(id)
}

func (userService *UserService) GetUserDetail(id uint) (*response.DetailUserResponse, error) {
	// Ambil data user lengkap dari repository

	user, err := userService.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Mapping data ke struct response baru
	image := utilities.ConvertPath(user.Image)
	detailResponse := response.DetailUserResponse{
		ID:                int(user.ID),
		Username:          user.Username,
		Email:             user.Email,
		RoleID:            int(user.RoleID),
		RoleName:          user.Role.Name,
		SchoolID:          int(user.UserSchool.School.ID),
		SchoolName:        user.UserSchool.School.SchoolName,
		VerificationEmail: user.IsVerification,
		Status:            mapStatus(user.IsBlock),
		Image:             image,
	}

	return &detailResponse, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	getUser, err := repositories.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &getUser, nil
}

func (userService *UserService) CreateUser(userRequest *request.UserCreateRequest, userID int) (*models.User, error) {
	// check existing registered email
	_, err := repositories.GetUserByEmail(userRequest.Email)
	if err == nil {
		return nil, fmt.Errorf("Email has already been registered.")
	}

	// check format email
	err = utilities.ValidateEmail(userRequest.Email)
	if err != nil {
		return nil, err
	}

	// validate password, use 1 number and upper+lowe char and password length
	// // err = utilities.ValidatePassword(userRequest.Password)
	// // if err != nil {
	// // 	return nil, err
	// // }

	if userRequest.Username != "" {
		// validate username
		err = utilities.ValidateUsername(userRequest.Username)
		if err != nil {
			return nil, err
		}

		valid := repositories.ValidateDuplicateUsername(userRequest.Username)
		if !valid {
			return nil, fmt.Errorf("Username duplicated.")
		}
	}

	// set role if userID 0, because that from case self register. that process has not set role and assummed role id set to 2
	if userID == 0 {
		userRequest.RoleID = 2
	}

	user := models.User{
		RoleID:         userRequest.RoleID,
		Username:       userRequest.Username,
		Email:          userRequest.Email,
		IsVerification: false,
		IsBlock:        true,
	}

	// set user id to created by and updated by
	user.Master.CreatedBy = userID
	user.Master.UpdatedBy = userID

	dataUser, err := repositories.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	if userRequest.SchoolID != 0 {
		userSchool := models.UserSchool{
			UserID:   dataUser.Master.ID,
			SchoolID: userRequest.SchoolID,
		}

		_, err := repositories.CreateUserSchool(&userSchool)
		if err != nil {
			return nil, err
		}
	}

	school, err := userService.schoolRepository.GetSchoolByID(userRequest.SchoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve school data: %v", err)
	}

	// Generate a verification token (could be a JWT or a unique string)
	verificationLink := os.Getenv("VERIFICATION_LINK")
	subject := constants.EmailConfirmationText
	emailTemplate := utilities.GenerateEmailBodyVerification()

	userName := dataUser.Username
	schoolLogo := utilities.ConvertPath(school.SchoolLogo)
	token, err := GenerateTokenAndSendEmail(dataUser.Email, userName, schoolLogo, verificationLink, subject, emailTemplate, "TK01")
	if err != nil {
		return nil, err
	}
	if verificationLink != "" {
		verif := models.TempVerificationEmail{
			UserID:           dataUser.Master.ID,
			VerificationLink: token,
			Type:             nil,
		}
		_, err := repositories.CreateVerificationLink(&verif)
		if err != nil {
			return nil, err
		}
	}

	return dataUser, nil
}

func (userService *UserService) UpdateUserService(id uint, userRequest *request.UserUpdateRequest, userID int) (*models.User, error) {
	// check id
	user, err := userService.userRepository.GetUserByIDPass(id)
	if err != nil {
		return nil, fmt.Errorf("User not found.")
	}

	// check username
	err = utilities.ValidateUsername(userRequest.Username)
	if err != nil {
		return nil, err
	}

	if user.Username != userRequest.Username {
		valid := repositories.ValidateDuplicateUsername(userRequest.Username)
		if !valid {
			return nil, fmt.Errorf("Username duplicated.")
		}
		user.Username = userRequest.Username
	}

	if user.IsVerification == false {
		return nil, fmt.Errorf("please verification")
	}

	if userRequest.Status != "" {
		if userRequest.Status == "Aktif" {
			user.IsBlock = false
		} else if userRequest.Status == "Tidak Aktif" {
			user.IsBlock = true
		}
	}

	user.RoleID = uint(userRequest.RoleID)
	user.Master.UpdatedBy = userID

	// Update user (hanya password yang diupdate) menggunakan pointer
	dataUser, err := userService.userRepository.UpdateUserRepository(user)
	if err != nil {
		return nil, err
	}

	return dataUser, nil
}

func (userService *UserService) DeleteUserService(user *models.User, userID int) (*models.User, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime

	userSchool, err := repositories.GetUserSchoolByUserId(user.Master.ID)
	if err != nil {
		return nil, err
	}

	userSchool.Master.DeletedAt = currentTimePointer
	userSchool.Master.DeletedBy = &userID
	_, err = repositories.UpdateUserSchool(&userSchool)
	if err != nil {
		return nil, err
	}

	user.Master.DeletedAt = currentTimePointer
	user.Master.DeletedBy = &userID
	dataUser, err := userService.userRepository.DeleteUserRepository(user)
	if err != nil {
		return nil, err
	}

	return dataUser, nil
}

func (userService *UserService) ValidateEmail(token string) error {
	if strings.Contains(token, "%") {
		// Decode the token if needed
		decodedToken, err := url.QueryUnescape(token)
		if err != nil {
			return err
		}
		token = decodedToken
	}

	email, time, _, err := utilities.DecodeToken(token)
	if err != nil {
		return err
	}

	timeStr, _ := strconv.ParseInt(time, 10, 64)
	_, hasPassed := utilities.ValidateTime(timeStr)
	if hasPassed {
		return fmt.Errorf(constants.MessageErrorLinkExpired)
	}

	getUser, err := repositories.GetUserByEmailIsBlock(email)
	if err != nil {
		return err
	}

	getUser.UpdatedBy = int(getUser.Master.ID)
	getUser.IsVerification = true
	getUser.IsBlock = false

	_, err = userService.userRepository.DeleteUserRepository(&getUser)
	if err != nil {
		return err
	}

	err = userService.userRepository.DeleteAllTempVerificationEmails(int(getUser.ID))
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) GenerateTokenChangePassword(email string) error {

	getUser, err := repositories.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("User Not Found")
	}

	// Generate the token and send the email
	verificationLink := os.Getenv("CHANGE_PASSWORD_LINK")
	subject := "Ubah Password"
	emailTemplate := utilities.GenerateEmailBodyChangePassword()

	schoolLogo := utilities.ConvertPath(getUser.UserSchool.School.SchoolLogo)

	token, err := GenerateTokenAndSendEmail(email, getUser.Username, schoolLogo, verificationLink, subject, emailTemplate, "TK02") // Remove := here
	if strings.Contains(token, "%") {
		// Decode the token if needed
		decodedToken, err := url.QueryUnescape(token)
		if err != nil {
			return err
		}
		token = decodedToken
	}

	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	typeChangePassword := "change_password"

	if verificationLink != "" {
		verif := models.TempVerificationEmail{
			UserID:           getUser.Master.ID,
			VerificationLink: token,
			IsValid:          true,
			Type:             &typeChangePassword,
		}
		_, err := repositories.CreateVerificationLink(&verif)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (userService *UserService) VerifyTokenChangePassword(token string) error {
	if strings.Contains(token, "%") {
		// Decode the token if needed
		decodedToken, err := url.QueryUnescape(token)
		if err != nil {
			return err
		}
		token = decodedToken
	}

	email, time, _, err := utilities.DecodeToken(token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	getUser, err := repositories.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("User Not Found")
	}

	objectVerify, err := userService.userRepository.GetEmailVerification(int(getUser.ID), "change_password")
	if err != nil {
		// Jika tidak ditemukan atau ada error
		return fmt.Errorf("invalid token: %v", err)
	}

	if objectVerify.VerificationLink == token && !objectVerify.IsValid {
		return fmt.Errorf(constants.MessageErrorLinkExpired)
	}

	timeStr, _ := strconv.ParseInt(time, 10, 64)
	_, hasPassed := utilities.ValidateTime(timeStr)
	if hasPassed {
		return fmt.Errorf(constants.MessageErrorLinkExpired)
	}

	_, err = repositories.GetUserByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) ChangePassword(token string, newPassword string) (*models.User, error) {
	// check id
	var objectVerify models.TempVerificationEmail
	if strings.Contains(token, "%") {
		// Decode the token if needed
		decodedToken, err := url.QueryUnescape(token)
		if err != nil {
			return nil, err
		}
		token = decodedToken
	}
	email, _, contentCode, err := utilities.DecodeToken(token)
	if err != nil {
		return nil, err
	}

	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User not found.")
	}

	if contentCode == "TK02" {
		objectVerify, err := userService.userRepository.GetEmailVerification(int(user.ID), "change_password")
		if err != nil {
			// Jika tidak ditemukan atau ada error
			return nil, fmt.Errorf("invalid token: %v", err)
		}

		if objectVerify.VerificationLink == token && !objectVerify.IsValid {
			return nil, fmt.Errorf(constants.MessageErrorLinkExpired)
		}
	}

	// check old password and new password identic or not
	err = utilities.ComparePassword(user.Password, newPassword)
	if err != nil {
		return nil, err
	}

	// Hash password baru
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	user.Master.UpdatedBy = int(user.Master.ID)

	// Update user (hanya password yang diupdate) menggunakan pointer
	dataUser, err := userService.userRepository.UpdateUserRepository(&user)
	if err != nil {
		return nil, err
	}

	if contentCode == "TK02" {
		objectVerify.IsValid = false
		_ = userService.userRepository.UpdateTempVerificationEmail(&objectVerify)
	}

	return dataUser, nil
}

func (userService *UserService) ChangePasswordWithoutToken(userID uint, oldPassword, newPassword string) (*models.User, error) {

	user, err := repositories.NewUserRepository().GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("User not found: %v", err)
	}

	// Validasi apakah oldPassword sesuai dengan password yang ada di database
	if err := utilities.CompareOldPassword(user.Password, string(oldPassword)); err != nil {
		return nil, err
	}

	err = utilities.ComparePassword(user.Password, newPassword)
	if err != nil {
		return nil, err
	}

	// Hash password baru
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	user.Master.UpdatedBy = int(user.Master.ID)

	// Update user (hanya password yang diupdate) menggunakan pointer
	dataUser, err := userService.userRepository.UpdateUserRepository(&user)
	if err != nil {
		return nil, err
	}

	return dataUser, nil
}

func GenerateTokenAndSendEmail(email, userName, schoolLogo, verificationLink, subject, emailTemplate string, contentCode string) (string, error) {
	// Generate a verification token (could be a JWT or a unique string)
	verificationToken, err := utilities.GenerateVerificationToken(email, contentCode)
	if err != nil {
		return "", fmt.Errorf("Failed to generate verification token.")
	}

	baseUrl := os.Getenv("BASE_URL")
	// Encode the verification token to prevent issues with special characters in URLs
	fullVerificationLink := fmt.Sprintf("%s%s?token=%s", baseUrl, verificationLink, url.QueryEscape(verificationToken))

	// Prepare the data for the template
	emailData := request.EmailData{
		ConfirmationLink: fullVerificationLink,
		SchoolLogo:       schoolLogo,
		UserName:         userName,
	}

	// Use html/template to parse the email template
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse email template: %w", err)
	}

	// Execute the template with data and capture the result
	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, emailData); err != nil {
		return "", fmt.Errorf("Failed to execute email template: %w", err)
	}

	// Get the final email body as a string
	emailBody := bodyBuffer.String()

	// Send verification email
	// err = utilities.SendVerificationEmail(email, verificationToken, subject, emailBody)
	err = utilities.SendVerificationEmail(email, "", subject, emailBody)
	if err != nil {
		return "", fmt.Errorf("Failed to send verification email - %w", err)
	}

	return verificationToken, nil
}

func (userService *UserService) GenerateFileExcelForUser(c *fiber.Ctx, userID int) (*bytes.Buffer, error) {
	filename := "attachment;filename=import_User_data.xlsx"
	headers := []string{
		"Role", "Nama Sekolah", "Username", "Email",
	}

	userSchool, err := repositories.GetUserSchoolByUserId(uint(userID))
	if err != nil {
		return nil, err
	}

	school, err := userService.schoolRepository.GetSchoolByID(userSchool.SchoolID)
	if err != nil {
		return nil, err
	}

	// Example data for the upload format
	exampleData := [][]string{
		{"Admin", school.SchoolName, "username1", "email1@example.com"},
	}

	roleOptions, err := repositories.GetRoleNames()
	if err != nil {
		return nil, err
	}

	buffer, err := utilities.GenerateFileExcelUser(c, headers, filename, exampleData, roleOptions)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (userService *UserService) HandleExcelFileUser(file *multipart.FileHeader, c *fiber.Ctx, userID int) error {

	userSchool, err := repositories.GetUserSchoolByUserId(uint(userID))
	if err != nil {
		return fmt.Errorf("failed to get user school: %v", err)
	}

	school, err := userService.schoolRepository.GetSchoolByID(userSchool.SchoolID)
	if err != nil {
		return fmt.Errorf("failed to get school by ID: %v", err)
	}

	// Open the uploaded file with buffered reading
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Use buffered reader for better memory efficiency
	reader := bufio.NewReader(src)

	// Configure excelize with optimized settings
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return fmt.Errorf("failed to read Excel file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("error closing excel file: %v\n", err)
		}
	}()

	// Get sheet name and validate
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("invalid excel file: no sheets found")
	}

	// Read rows with error handling
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows: %v", err)
	}

	// Validate minimum rows
	if len(rows) < 2 { // At least header + 1 data row
		return fmt.Errorf("file must contain at least one data row")
	}

	// schoolname validation
	SchoolName := rows[1][1]
	if SchoolName != school.SchoolName {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "School name in the file does not match the current user's school",
		})
	}

	// Process rows in batches for better memory management
	const batchSize = 500
	var allErrors []response.ResponseErrorUploadUser

	for i := 1; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		errors := userService.processRowsUser(rows[i:end], userID)
		allErrors = append(allErrors, errors...)
	}

	// Prepare response
	message := "Success upload data"
	if len(allErrors) > 0 {
		message = "Failed upload data"
	}
	response := fiber.Map{
		"message": message,
		"data":    allErrors,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (userService *UserService) HandleCSVFileUser(file *multipart.FileHeader, c *fiber.Ctx, userID int) error {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("Failed to open file")
	}
	defer src.Close()

	// Create a CSV reader
	reader := csv.NewReader(src)

	// Read all rows
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to read CSV file")
	}

	// Process the rows (skipping the header row)
	students := userService.processRowsUser(rows[1:], userID)

	// Return a response with imported data
	responseMessage := ""
	if len(students) > 0 {
		responseMessage = "Failed upload data"
	} else {
		responseMessage = "Success upload data"
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":         responseMessage,
		"totalDataFailed": len(students),
		"data":            students,
	})
}

func (userService *UserService) processRowsUser(rows [][]string, userID int) []response.ResponseErrorUploadUser {
	var users []models.User
	var userSchools []models.UserSchool
	var dataUserErrors []response.ResponseErrorUploadUser
	var emailTasks []struct {
		user   models.User
		school models.School
		token  string
		verif  models.TempVerificationEmail
	}

	// Pre-fetch roles and schools in bulk
	roleMap, err := userService.roleRepository.GetRolesByNames(getRoleNames(rows))
	if err != nil {
		return []response.ResponseErrorUploadUser{{Reason: "Failed to fetch roles"}}
	}

	schoolMap, err := userService.schoolRepository.GetSchoolsByNames(getSchoolNames(rows))
	if err != nil {
		return []response.ResponseErrorUploadUser{{Reason: "Failed to fetch schools"}}
	}

	// Validasi apakah seluruh row[3] berisi "email1@example.com"
	allEmailsAreDummy := true
	for _, row := range rows {
		if row[3] == "email1@example.com" {
			// Skip this row from processing (we'll show a dummy error later)
			dataUserErrors = append(dataUserErrors, response.ResponseErrorUploadUser{
				Username: row[2],
				Email:    row[3],
				Reason:   "Dummy Email - No Data",
			})
			continue
		}

		if row[3] != "email1@example.com" {
			allEmailsAreDummy = false
			break
		}
	}

	// Jika seluruh email adalah "email1@example.com", kembalikan error
	if allEmailsAreDummy {
		var errors []response.ResponseErrorUploadUser
		for _, row := range rows {
			errors = append(errors, response.ResponseErrorUploadUser{
				Username: row[2],
				Email:    row[3],
				Reason:   "No data (all emails are dummy)",
			})
		}
		return errors
	}

	// Prepare bulk data
	for _, row := range rows {
		if row[0] == "" || row[1] == "" || row[2] == "" || row[3] == "" {
			continue
		}
		if len(row) < 4 || row[3] == "email1@example.com" {
			continue
		}

		roleID := roleMap[row[0]]
		schoolID := schoolMap[row[1]]

		if roleID == 0 || schoolID == 0 {
			reason := ""
			if roleID == 0 && schoolID > 0 {
				reason = "Role Tidak Tersedia"
			} else if roleID > 0 && schoolID == 0 {
				reason = "Data Sekolah Tidak Ada"
			} else {
				reason = "Data Role dan Sekolah Tidak Ada"
			}

			dataUserErrors = append(dataUserErrors, response.ResponseErrorUploadUser{
				Username: row[2],
				Email:    row[3],
				Reason:   reason,
			})
			continue
		}

		user := models.User{
			RoleID:         roleID,
			Username:       row[2],
			Email:          row[3],
			IsVerification: false,
			IsBlock:        true,
			Master: models.Master{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
		}
		users = append(users, user)

		userSchool := models.UserSchool{
			SchoolID: schoolID,
		}
		userSchools = append(userSchools, userSchool)
	}

	// Bulk create users and user schools
	createdUsers, dataUserErrorsNew, err := userService.userRepository.BulkValidateAndCreateUsers(users, userSchools)
	if err != nil {
		fmt.Println("Bulk Validate and Create Users Error: ", err)
		return nil
	}

	dataUserErrors = append(dataUserErrors, dataUserErrorsNew...)

	// Prepare email verification tasks for successfully created users
	for _, user := range createdUsers {
		// Generate verification token
		token, err := utilities.GenerateVerificationToken(user.Email, "TK01")
		if err != nil {
			continue
		}

		if user.UserSchool == nil {
			log.Printf("Warning: UserSchool is nil for user ID: %d", user.ID)
			continue
		}

		school, err := userService.schoolRepository.GetSchoolByID(user.UserSchool.SchoolID)
		if err != nil {
			log.Printf("Error getting school data for user ID %d: %v", user.ID, err)
			continue
		}

		verif := models.TempVerificationEmail{
			UserID:           user.Master.ID,
			VerificationLink: token,
			Type:             nil,
		}

		emailTasks = append(emailTasks, struct {
			user   models.User
			school models.School
			token  string
			verif  models.TempVerificationEmail
		}{user, school, token, verif})
	}

	// Convert to EmailTask type
	var tasks []EmailTask
	for _, task := range emailTasks {
		tasks = append(tasks, EmailTask{
			User:   task.user,
			School: task.school,
			Token:  task.token,
			Verif:  task.verif,
		})
	}

	// Launch async processing
	go func(tasks []EmailTask) {
		workers := 2
		emailQueue := NewEmailQueue(workers)
		emailQueue.ProcessEmails(tasks)
	}(tasks)

	return dataUserErrors
}

func (userService *UserService) CheckEmail(email string) bool {
	_, err := repositories.GetUserByEmail(email)
	if err != nil {
		return true
	}
	return false
}

func CheckUsername(username string) bool {
	_, err := repositories.GetUserByUsername(username)
	if err != nil {
		return true
	}
	return false
}

func (userService *UserService) validateRow(row []string) *response.ResponseErrorUploadUser {
	// Validate email
	if err := utilities.ValidateEmail(row[3]); err != nil {
		return &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Format Email Salah",
		}
	}

	// Validate username
	if err := utilities.ValidateUsername(row[2]); err != nil {
		return &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Format Username Salah",
		}
	}

	// Check for duplicate username
	if !repositories.ValidateDuplicateUsername(row[2]) {
		return &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Username Sudah Terpakai",
		}
	}

	// Check for duplicate email
	if _, err := repositories.GetUserByEmail(row[3]); err == nil {
		return &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Email Sudah Terpakai",
		}
	}

	return nil
}

func (userService *UserService) getRoleAndSchool(row []string) (uint, uint, *response.ResponseErrorUploadUser) {
	// Get role
	listRole, err := userService.roleRepository.GetAllRole(1, 1, row[0], 5)
	if err != nil || len(listRole) == 0 {
		return 0, 0, &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Role Tidak Tersedia",
		}
	}

	// Get school
	listSchool, err := repositories.GetAllSchool(1, 1, row[1])
	if err != nil || len(listSchool) == 0 {
		return 0, 0, &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Data Sekolah Tidak Ada",
		}
	}

	return listRole[0].ID, listSchool[0].ID, nil
}

func (userService *UserService) createUserUpload(row []string, roleId, schoolId uint, userID int) *response.ResponseErrorUploadUser {

	// Create user request
	userRequest := request.UserCreateRequest{
		SchoolID: schoolId,
		RoleID:   roleId,
		Username: row[2],
		Email:    row[3],
	}

	// Create user in repository
	if _, err := userService.CreateUser(&userRequest, userID); err != nil {
		return &response.ResponseErrorUploadUser{
			Username: row[2],
			Email:    row[3],
			Reason:   "Terjadi Kesalahan Ketika Proses Simpan Data",
		}
	}

	return nil
}

func (userService *UserService) ResendEmailVerification(userID int) error {
	emailVerification, err := userService.userRepository.GetEmailVerification(userID, "")
	if err != nil {
		return err
	}

	getUser, err := userService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return err
	}

	password, err := utilities.GeneratePassword()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Generate a verification token (could be a JWT or a unique string)
	verificationLink := os.Getenv("VERIFICATION_LINK")
	subject := constants.EmailConfirmationText
	emailTemplate := utilities.GenerateEmailBodyVerification()

	userName := getUser.Username
	schoolLogo := utilities.ConvertPath(getUser.UserSchool.School.SchoolLogo)
	token, err := GenerateTokenAndSendEmail(getUser.Email, userName, schoolLogo, verificationLink, subject, emailTemplate, "TK01")
	if err != nil {
		return err
	}

	getUser.Password = string(hashPassword)
	_, err = userService.userRepository.UpdateUserRepository(&getUser)
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	emailVerification.VerificationLink = token
	emailVerification.Type = nil
	err = userService.userRepository.UpdateTempVerificationEmail(&emailVerification)
	if err != nil {
		return err
	}

	return nil
}

func mapStatus(isBlocked bool) string {
	if isBlocked {
		return "Tidak Aktif"
	}
	return "Aktif"
}

func (userService *UserService) UpdateUserPhotoService(request request.UpdateUserImageRequest, userID int) (*models.User, error) {
	// Get user by ID
	user, err := userService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("User not found")
	}

	// Delete old image if exists
	if user.Image != "" {
		oldImagePath := filepath.Join("./upload/user/image", user.Image)
		if err := os.Remove(oldImagePath); err != nil {
			// Log error but continue
			fmt.Printf("Error deleting old image: %v\n", err)
		}
	}

	// Update user image
	user.Image = request.Image
	user.Master.UpdatedBy = userID

	// Save to database
	updatedUser, err := userService.userRepository.UpdateUserRepository(&user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func getRoleNames(rows [][]string) []string {
	roleNames := make([]string, 0, len(rows))
	seen := make(map[string]bool)

	for _, row := range rows {
		if len(row) > 0 && row[0] != "" && !seen[row[0]] {
			roleNames = append(roleNames, row[0])
			seen[row[0]] = true
		}
	}

	return roleNames
}

func getSchoolNames(rows [][]string) []string {
	schoolNames := make([]string, 0, len(rows))
	seen := make(map[string]bool)

	for _, row := range rows {
		if len(row) > 1 && row[1] != "" && !seen[row[1]] {
			schoolNames = append(schoolNames, row[1])
			seen[row[1]] = true
		}
	}

	return schoolNames
}

type EmailTask struct {
	User   models.User
	School models.School
	Token  string
	Verif  models.TempVerificationEmail
}

type EmailQueue struct {
	workers   int
	workersWg sync.WaitGroup
	tasksChan chan EmailTask
}

func NewEmailQueue(workers int) *EmailQueue {
	return &EmailQueue{
		workers:   workers,
		tasksChan: make(chan EmailTask),
	}
}

func (eq *EmailQueue) Start() {
	for i := 0; i < eq.workers; i++ {
		eq.workersWg.Add(1)
		go eq.worker(i)
	}
}

func (eq *EmailQueue) worker(id int) {
	defer eq.workersWg.Done()

	verificationLink := os.Getenv("VERIFICATION_LINK")
	subject := constants.EmailConfirmationText
	emailTemplate := utilities.GenerateEmailBodyVerification()

	// Create a batch of verification links
	var verificationLinks []*models.TempVerificationEmail

	for task := range eq.tasksChan {
		// Add 1 second delay between emails
		time.Sleep(1 * time.Second)

		schoolLogo := utilities.ConvertPath(task.School.SchoolLogo)
		token, err := GenerateTokenAndSendEmail(
			task.User.Email,
			task.User.Username,
			schoolLogo,
			verificationLink,
			subject,
			emailTemplate,
			"TK01",
		)

		if err != nil {
			log.Printf("Worker %d: Error sending email to %s: %v", id, task.User.Email, err)
			continue
		}

		// Create verification record
		verif := &models.TempVerificationEmail{
			UserID:           task.User.ID,
			VerificationLink: token,
			Type:             nil,
			IsValid:          true,
		}
		verificationLinks = append(verificationLinks, verif)

		// Bulk insert when batch size is reached or channel is closed
		if len(verificationLinks) >= 500 {
			if err := repositories.BulkCreateVerificationLinks(verificationLinks); err != nil {
				log.Printf("Worker %d: Error creating verification links in bulk: %v", id, err)
			}
			verificationLinks = verificationLinks[:0] // Clear the slice while keeping capacity
		}
	}

	// Insert any remaining verification links
	if len(verificationLinks) > 0 {
		if err := repositories.BulkCreateVerificationLinks(verificationLinks); err != nil {
			log.Printf("Worker %d: Error creating final batch of verification links: %v", id, err)
		}
	}
}

func (eq *EmailQueue) ProcessEmails(tasks []EmailTask) {
	// Start workers
	eq.Start()

	// Send tasks to workers
	go func() {
		for _, task := range tasks {
			eq.tasksChan <- task
		}
		close(eq.tasksChan)
	}()

	// Wait for all workers to complete
	eq.workersWg.Wait()
}
