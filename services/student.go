package services

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"schoolPayment/configs"
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"
)

var (
	studentServiceInstance *StudentService
	studentServiceOnce     sync.Once
)

type StudentService struct {
	studentRepository     repositories.StudentRepositoryInteface
	userRepository        repositories.UserRepository
	schoolClassRepository repositories.SchoolClassRepositoryInterface
	schoolYearRepository  repositories.SchoolYearRepository
	schoolGradeRepository  repositories.SchoolGradeRepositoryInterface
	userService           UserService
}

func NewStudentService(
	studentRepo repositories.StudentRepositoryInteface,
	userRepo repositories.UserRepository,
	schoolClassRepo repositories.SchoolClassRepositoryInterface,
	schoolYearRepo repositories.SchoolYearRepository,
	schoolGredeRepo repositories.SchoolGradeRepositoryInterface,
	userService UserService,
) *StudentService {
	return &StudentService{
		studentRepository:     studentRepo,
		userRepository:        userRepo,
		schoolClassRepository: schoolClassRepo,
		schoolYearRepository:  schoolYearRepo,
		schoolGradeRepository: schoolGredeRepo,
		userService:           userService,
	}
}

func GetStudentService() StudentService {
	studentServiceOnce.Do(func() {
		studentRepository := repositories.NewStudentRepository(configs.DB)
		userRepository := repositories.NewUserRepository()
		schoolClassRepository := repositories.NewSchoolClassRepository()
		schoolYearRepository := repositories.NewSchoolYearRepository(configs.DB)
		userService := GetUserService()

		studentServiceInstance = &StudentService{
			studentRepository:     studentRepository,
			userRepository:        userRepository,
			schoolClassRepository: schoolClassRepository,
			schoolYearRepository:  schoolYearRepository,
			userService:           userService,
		}
	})
	return *studentServiceInstance
}

func (studentService *StudentService) GetAllStudent(page int, limit int, search string, userID int, status string, gradeID int, yearId int, schoolID int, searchNis string, classID int, sortBy string, sortOrder string, studentId int, isActive *bool) (response.StudentListResponse, error) {
	var resp response.StudentListResponse
	resp.Limit = limit
	resp.Page = page
	resp.TotalData = 0
	resp.TotalPage = 0

	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		resp.Data = []response.DetailStudentResponse{}
		return resp, nil
	}

	sortBy = utilities.ChangeStringSortByStudent(sortBy)
	sortBy = utilities.ToSnakeCase(sortBy)

	// Get students
	dataStudent, totalPage, totalData, err := repositories.GetAllStudent(page, limit, search, user, status, gradeID, yearId, schoolID, searchNis, classID, sortBy, sortOrder, studentId, isActive)
	if err != nil {
		resp.Data = []response.DetailStudentResponse{}
		return resp, nil
	}

	resp.Data = dataStudent
	resp.TotalData = totalData
	resp.TotalPage = totalPage
	return resp, nil
}

func (studentService *StudentService) GetStudentByID(id uint, userID int) (*response.DetailStudentResponse, error) {
	schoolClassName := ""
	schoolGradeName := ""
	schoolName := ""
	schoolClassID := uint(0)
	schoolGradeID := uint(0)
	schoolID := uint(0)
	var emailUser string

	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	getStudent, err := studentService.studentRepository.GetStudentByID(id, user)
	if err != nil {
		return nil, err
	}

	getSchoolGrade, err := studentService.schoolGradeRepository.GetSchoolGradeByID(getStudent.SchoolGradeID)
	if err == nil {
		schoolGradeName = getSchoolGrade.SchoolGradeName
		schoolGradeID = getSchoolGrade.ID
	}

	// Get School Class
	getSchoolClass, err := studentService.schoolClassRepository.GetSchoolClassByID(getStudent.SchoolClassID)
	if err == nil {
		schoolClassName = getSchoolClass.SchoolClassName
		schoolClassID = getSchoolClass.ID
	}

	GetDataSchool, err := repositories.GetSchoolByStudentId(id)
	if err == nil {
		schoolName = GetDataSchool.SchoolName
		schoolID = GetDataSchool.ID
	}

	getEmailUser, err := repositories.GetUserByStudentId(getStudent.ID)
	if err == nil {
		emailUser = getEmailUser
	}

	// Get School Year Name
	schoolYearName := ""
	if getStudent.SchoolYearID != 0 {
		schoolYear, err := studentService.schoolYearRepository.GetSchoolYearByID(uint(getStudent.SchoolYearID))
		if err == nil {
			schoolYearName = schoolYear.SchoolYearName
		}
	}

	detailStudent := response.DetailStudentResponse{
		ID:                 getStudent.ID,
		Nisn:               getStudent.Nisn,
		RegistrationNumber: getStudent.RegistrationNumber,
		Nis:                getStudent.Nis,
		Nik:                getStudent.Nik,
		FullName:           getStudent.FullName,
		Gender:             getStudent.Gender,
		Religion:           getStudent.Religion,
		Citizenship:        getStudent.Citizenship,
		BirthPlace:         getStudent.BirthPlace,
		BirthDate:          getStudent.BirthDate,
		Address:            getStudent.Address,
		SchoolGradeID:      uint(schoolGradeID),
		SchoolGrade:        schoolGradeName,
		SchoolClassID:      schoolClassID,
		SchoolClass:        schoolClassName,
		SchoolYearID:       getStudent.SchoolYearID,
		NoHandphone:        getStudent.NoHandphone,
		Height:             getStudent.Height,
		Weight:             getStudent.Weight,
		MedicalHistory:     getStudent.MedicalHistory,
		DistanceToSchool:   getStudent.DistanceToSchool,
		Sibling:            getStudent.Sibling,
		NickName:           getStudent.NickName,
		Email:              emailUser,
		EntryYear:          getStudent.EntryYear,
		Status:             getStudent.Status,
		Image:              getStudent.Image,
		SchoolID:           schoolID,
		SchoolName:         schoolName,
		SchoolYearName:     schoolYearName,
	}

	return &detailStudent, nil
}

func (studentService *StudentService) CreateStudent(createStudentRequest request.CreateStudentRequest, userID int) (*models.Student, error) {
	// Validasi manual required fields
	var phoneNumberFormat string
	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}
	if err := validateCreateStudentRequest(createStudentRequest, user, studentService); err != nil {
		return nil, err
	}

	if createStudentRequest.EmailParent != "" {
		err := utilities.ValidateEmail(createStudentRequest.EmailParent)
		if err != nil {
			return nil, err
		}
	}

	if createStudentRequest.EmailParent == "" {
		return nil, fmt.Errorf("Please Fill Email Parent")
	}

	if strings.HasPrefix(createStudentRequest.NoHandphone, "+") {
		phoneNumberFormat = utilities.FormatPhoneNumber(createStudentRequest.NoHandphone)
	} else {
		phoneNumberFormat = createStudentRequest.NoHandphone
	}

	// Cek apakah NIS sudah ada
	if err := validateDuplicateNis(createStudentRequest.Nis, user, studentService); err != nil {
		return nil, err
	}

	userIdFromCheck, valid, errCheckEmail := repositories.GetUserByEmailAndFilterSchoolId(createStudentRequest.EmailParent, createStudentRequest.SchoolID)

	if !valid {
		return nil, fmt.Errorf("Email Terdaftar di tempat lain, silahkan pakai email yang berbeda. ")
	}

	// Map DTO to Model
	student := models.Student{
		Nis:          createStudentRequest.Nis,
		FullName:     strings.ToUpper(createStudentRequest.FullName),
		Gender:       createStudentRequest.Gender,
		Religion:     createStudentRequest.Religion,
		BirthPlace:   strings.ToUpper(createStudentRequest.BirthPlace),
		SchoolYearID: createStudentRequest.SchoolYearID,
		// Parse BirthDate string to time.Time
		BirthDate: func() *time.Time {
			t, _ := time.Parse(constants.DateFormatYYYYMMDD, createStudentRequest.BirthDate)
			return &t
		}(),
		Address:       createStudentRequest.Address,
		NoHandphone:   phoneNumberFormat,
		Email:         createStudentRequest.EmailParent,
		SchoolGradeID: createStudentRequest.SchoolGradeID,
		SchoolClassID: createStudentRequest.SchoolClassID,
		Status:        createStudentRequest.Status,
	}

	// ini logic untuk upload image, nanti di pakai ketika di butuhkan
	// file, err := c.FormFile("image")
	// if err != nil {
	// 	// Log or return the error if image is missing or not uploaded correctly
	// 	fmt.Println("Error retrieving file:", err)
	// } else {
	// 	epoch := time.Now().Unix()

	// 	extension := filepath.Ext(file.Filename)

	// 	newFileName := fmt.Sprintf("%d%s", epoch, extension)

	// 	// Ensure the upload directory exists
	// 	uploadDir := "./upload/"
	// 	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
	// 		err = os.MkdirAll(uploadDir, 0755) // Create the directory with necessary permissions
	// 		if err != nil {
	// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 				"error": "Unable to create upload directory",
	// 			})
	// 		}
	// 	}

	// 	filePath := filepath.Join(uploadDir, newFileName)
	// 	if err := c.SaveFile(file, filePath); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Unable to save image",
	// 		})
	// 	}

	// 	student.Image = newFileName
	// }

	dataStudent, err := studentService.studentRepository.CreateStudent(&student)
	if err != nil {
		return nil, err
	}

	username, err := utilities.GenerateRandomUsername(5)
	if err != nil {
		return nil, err
	}

	userRequest := request.UserCreateRequest{
		SchoolID: createStudentRequest.SchoolID,
		RoleID:   uint(2),
		Username: username,
		Email:    createStudentRequest.EmailParent,
	}

	if errCheckEmail != nil {
		userCreate, err := studentService.userService.CreateUser(&userRequest, userID)
		if err != nil {
			return nil, err
		}
		userIdFromCheck = int(userCreate.ID)
	}

	_, err = studentService.CreateUserStudentService(userIdFromCheck, int(student.ID), userID)
	if err != nil {
		return nil, err
	}

	newData, err := json.Marshal(dataStudent)
	if err != nil {
		return nil, err
	}

	studentHistory := models.StudentHistory{
		Master: models.Master{
			CreatedBy: userID,
		},
		StudentID: student.ID,
		NewData:   string(newData),
		OldData:   "",
		Action:    "create",
	}

	_, err = studentService.studentRepository.CreateStudentHistory(&studentHistory)
	if err != nil {
		return nil, err
	}

	return dataStudent, err
}

func (studentService *StudentService) UpdateStudent(id uint, updateStudentRequest request.UpdateStudentRequest, userID int) (*models.Student, error) {
	var phoneNumberFormat string
	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	existingStudent, err := studentService.studentRepository.GetStudentByID(uint(id), user)
	if err != nil {
		return nil, err
	}

	// Validate NIS uniqueness if it's being changed
	if updateStudentRequest.Nis != existingStudent.Nis {
		if err := validateDuplicateNis(updateStudentRequest.Nis, user, studentService, id); err != nil {
			return nil, err
		}
	}

	OldData, err := json.Marshal(existingStudent)
	if err != nil {
		return nil, err
	}

	// file, err := c.FormFile("image")
	// if err == nil {
	// 	if existingStudent.Image != "" {
	// 		oldImagePath := filepath.Join("./upload", existingStudent.Image)
	// 		if err := os.Remove(oldImagePath); err != nil {
	// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 				"error": "Failed to delete old image",
	// 			})
	// 		}
	// 	}

	// 	newImageName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Ext(file.Filename))
	// 	if err := c.SaveFile(file, fmt.Sprintf("./upload/%s", newImageName)); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Failed to upload new image",
	// 		})
	// 	}
	// 	existingStudent.Image = newImageName
	// } else {
	// 	existingStudent.Image = existingStudent.Image
	// }

	// Buat objek student dengan data yang diperbarui
	if strings.HasPrefix(updateStudentRequest.NoHandphone, "+") {
		phoneNumberFormat = utilities.FormatPhoneNumber(updateStudentRequest.NoHandphone)
	} else {
		phoneNumberFormat = updateStudentRequest.NoHandphone
	}

	existingStudent.FullName = strings.ToUpper(updateStudentRequest.FullName)
	existingStudent.Gender = updateStudentRequest.Gender
	existingStudent.Religion = updateStudentRequest.Religion
	existingStudent.BirthPlace = strings.ToUpper(updateStudentRequest.BirthPlace)
	existingStudent.BirthDate = func() *time.Time {
		t, _ := time.Parse(constants.DateFormatYYYYMMDD, updateStudentRequest.BirthDate)
		return &t
	}()
	existingStudent.Address = updateStudentRequest.Address
	existingStudent.NoHandphone = phoneNumberFormat

	existingStudent.UpdatedBy = userID
	existingStudent.Status = updateStudentRequest.Status
	existingStudent.Nis = updateStudentRequest.Nis

	if updateStudentRequest.SchoolGradeID != 0 {
		schoolGrade, err := studentService.schoolGradeRepository.GetSchoolGradeByID(updateStudentRequest.SchoolGradeID)
		if err != nil {
			return nil, err
		}

		existingStudent.SchoolGradeID = schoolGrade.ID
	}

	if updateStudentRequest.SchoolClassID != 0 {
		schoolClass, err := studentService.schoolClassRepository.GetSchoolClassByID(updateStudentRequest.SchoolClassID)
		if err != nil {
			return nil, err
		}

		existingStudent.SchoolClassID = schoolClass.ID
	}

	if updateStudentRequest.SchoolYearID != 0 {
		schoolYear, err := studentService.schoolYearRepository.GetSchoolYearByID(uint(updateStudentRequest.SchoolYearID))
		if err != nil {
			return nil, err
		}

		existingStudent.SchoolYearID = schoolYear.ID
	}

	updatedStudent, err := studentService.studentRepository.UpdateStudent(&existingStudent)
	if err != nil {
		return nil, err
	}

	newData, err := json.Marshal(updatedStudent)
	if err != nil {
		return nil, err
	}

	studentHistory := models.StudentHistory{
		Master: models.Master{
			CreatedBy: userID,
		},
		StudentID: existingStudent.ID,
		NewData:   string(newData),
		OldData:   string(OldData),
		Action:    "update",
	}

	_, err = studentService.studentRepository.CreateStudentHistory(&studentHistory)
	if err != nil {
		return nil, err
	}

	return updatedStudent, nil
}

func (studentService *StudentService) DeleteStudentService(id uint, userID int) (*models.Student, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime

	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	student, err := studentService.studentRepository.GetStudentByID(id, user)
	if err != nil {
		return nil, fmt.Errorf("Student not found")
	}

	student.Master.DeletedAt = currentTimePointer
	student.Master.DeletedBy = &userID
	dataStudent, err := studentService.studentRepository.UpdateStudent(&student)
	if err != nil {
		return nil, err
	}

	newData, err := json.Marshal(dataStudent)
	if err != nil {
		return nil, err
	}

	oldData, _ := json.Marshal(student)

	studentHistory := models.StudentHistory{
		Master: models.Master{
			CreatedBy: userID,
		},
		StudentID: student.ID,
		NewData:   string(newData),
		OldData:   string(oldData),
		Action:    "delete",
	}

	_, err = studentService.studentRepository.CreateStudentHistory(&studentHistory)
	if err != nil {
		return nil, err
	}

	return dataStudent, nil
}

func (studentService *StudentService) CreateUserStudentService(id_user int, studentID int, auth int) (*models.UserStudent, error) {
	userStudent := &models.UserStudent{
		Master: models.Master{
			CreatedBy: auth,
		},
		UserID:    uint(id_user),
		StudentID: uint(studentID),
	}

	dataUserStudent, err := studentService.studentRepository.CreateUserStudentRepository(userStudent)
	if err != nil {
		return nil, err
	}

	return dataUserStudent, nil
}

func (studentService *StudentService) GetStudentByUserIDService(id uint) ([]models.Student, error) {
	return studentService.studentRepository.GetStudentByUserIdRepository(id)
}

func validateCreateStudentRequest(req request.CreateStudentRequest, user models.User, studentService *StudentService) error {
	// Check for duplicate NIS first
	if err := validateDuplicateNis(req.Nis, user, studentService); err != nil {
		return err
	}

	// Validate other required fields
	if req.FullName == "" {
		return fmt.Errorf("fullName is required")
	}
	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}
	if req.Religion == "" {
		return fmt.Errorf("religion is required")
	}
	if req.BirthPlace == "" {
		return fmt.Errorf("birthPlace is required")
	}
	if req.BirthDate == "" {
		return fmt.Errorf("birthDate is required")
	}
	if req.Address == "" {
		return fmt.Errorf("address is required")
	}
	if req.NoHandphone == "" {
		return fmt.Errorf("noHandphone is required")
	}
	return nil
}

func GenerateFileExcelForStudent(c *fiber.Ctx, db *gorm.DB) (*bytes.Buffer, error) {
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse user claims")
	}

	// Ambil userID dari userClaims
	userIDFloat, ok := userClaims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing user_id in user claims")
	}

	userID := uint(userIDFloat)

	// Panggil repository untuk mendapatkan jenjang sekolah
	schoolGrades, err := repositories.GetSchoolGradeByUser(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school grades: %w", err)
	}

	// Konversi hasil query ke format []string
	var schoolGradeNames string
	for _, grade := range schoolGrades {
		schoolGradeNames = grade.SchoolGradeName
	}

	// Panggil repository untuk mendapatkan sekolah
	schools, err := repositories.GetSchoolByUser(db, userID) // Pastikan kamu punya fungsi ini di repository
	if err != nil {
		return nil, fmt.Errorf("failed to get schools: %w", err)
	}

	// Konversi hasil query ke format []string
	var schoolNames string
	for _, school := range schools {
		schoolNames = school.SchoolName
	}

	// Pangil repository untuk dapetin kelas
	page := 1
	limit := 10
	search := ""
	sortBy := ""
	sortOrder := "asc"
	showDeletedData := false

	// Correct the struct initialization with userID as uint
	user, _ := repositories.GetUserByID2(uint(userID))

	schoolClass, _, _, err := repositories.GetAllSchoolClass(page, limit, search, sortBy, sortOrder, showDeletedData, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get school class: %w", err)
	}

	// Collect the school class names
	var schoolClassNames []string
	for _, class := range schoolClass {
		schoolClassNames = append(schoolClassNames, class.SchoolClassName)
	}

	// Inisialisasi repository dengan objek db
	schoolYearRepo := repositories.NewSchoolYearRepository(db)

	schoolYears, _, _, err := schoolYearRepo.GetAllSchoolYear(page, limit, search, sortBy, sortOrder, user.UserSchool.SchoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school years: %w", err)
	}

	// Ambil nama tahun ajaran
	var schoolYearNames []string
	for _, year := range schoolYears {
		schoolYearNames = append(schoolYearNames, year.SchoolYearName)
	}

	// Headers dan contoh data
	headers := []string{
		"NIS", "Nama Lengkap", "Tempat Lahir", "Tanggal Lahir", "Jenis Kelamin",
		"Agama", "Alamat", "Jenjang Sekolah", "Kelas", "No HP", "Email Orang Tua",
		"Status", "Sekolah", "Tahun Ajaran",
	}
	exampleData := [][]interface{}{
		{"442342424", "DAAM LARAVEL", "JAKARTA", "31-01-2000", "Laki-laki",
			"Lainnya", "Jl jalan ke mana", schoolGradeNames, "Pilih kelas", "08563142342", constants.DummyEmail,
			"Aktif", schoolNames, "Pilih tahun ajaran"},
	}

	// Dropdown configurations
	dropdowns := map[string][]utilities.DropdownOption{
		"E": {
			{Label: "Laki-laki", Value: "laki-laki"},
			{Label: "Perempuan", Value: "perempuan"},
		},
		"F": {
			{Label: "Islam", Value: "islam"},
			{Label: "Kristen", Value: "kristen"},
			{Label: "Katolik", Value: "katolik"},
			{Label: "Hindu", Value: "hindu"},
			{Label: "Budha", Value: "budha"},
			{Label: "Khonghucu", Value: "khonghucu"},
			{Label: "Lainnya", Value: "other"},
		},
		"L": {
			{
				Label: "Aktif",
				Value: "aktif",
			},
			{
				Label: "Tamat",
				Value: "tamat",
			},
			{
				Label: "Pindah Sekolah",
				Value: "pindah_sekolah",
			},
			{
				Label: "Dropout",
				Value: "dropout",
			},
		},
	}

	// For dynamic options like schoolClassNames and schoolYearNames
	for _, className := range schoolClassNames {
		dropdowns["I"] = append(dropdowns["I"], utilities.DropdownOption{
			Label: className,
			Value: className, // You might want to map this to an ID instead
		})
	}

	for _, yearName := range schoolYearNames {
		dropdowns["N"] = append(dropdowns["N"], utilities.DropdownOption{
			Label: yearName,
			Value: yearName, // You might want to map this to an ID instead
		})
	}

	// Define formats for specific columns
	formats := map[string]string{
		"D": "dd-mm-yyyy", // Column D is "Tanggal Lahir"
	}

	// Generate file Excel
	filename := "attachment;filename=import_student_data.xlsx"
	buffer, err := utilities.GenerateFileExcelStudent(c, headers, filename, exampleData, dropdowns, formats)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Function to handle Excel file (.xlsx or .xls)
func (studentService *StudentService) HandleExcelFileStudent(file *multipart.FileHeader, c *fiber.Ctx, userID int) error {
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

	// Process rows in batches for better memory management
	const batchSize = 500
	var allErrors []request.ResponseErrorUpload

	for i := 1; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		errors := studentService.processRowsBatch(rows[i:end], userID)
		allErrors = append(allErrors, errors...)
	}

	// Prepare response
	message := "Success upload data"
	if len(allErrors) > 0 {
		message = "Failed upload data"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"data":    allErrors,
	})
}

func (studentService *StudentService) processRowsBatch(rows [][]string, userID int) []request.ResponseErrorUpload {
	log.Printf("[DEBUG] Starting to process batch of %d rows", len(rows))

	var dataStudentErrors []request.ResponseErrorUpload
	var userStudentPairs []models.UserStudent
	var usersToCreate []models.User
	var studentHistories []models.StudentHistory

	// Get user for validation
	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		log.Printf("[ERROR] Failed to fetch user information: %v", err)
		return []request.ResponseErrorUpload{{
			Reason: "Failed to fetch user information",
		}}
	}
	log.Printf("[DEBUG] Successfully fetched user with ID: %d", userID)

	// Process all rows in bulk
	studentsToProcess, errors := studentService.validateAndPrepareStudentsBulk(rows, userID, user)
	if len(errors) > 0 {
		log.Printf("[WARN] Found %d validation errors during bulk preparation", len(errors))
		dataStudentErrors = append(dataStudentErrors, errors...)
	}
	log.Printf("[DEBUG] Prepared %d students for processing", len(studentsToProcess))

	// Collect all parent emails for bulk lookup
	parentEmails := make([]string, 0)
	emailToStudentMap := make(map[string][]models.Student)
	for _, student := range studentsToProcess {
		if student.Email != "" {
			parentEmails = append(parentEmails, student.Email)
			emailToStudentMap[student.Email] = append(emailToStudentMap[student.Email], student)
		}
	}
	log.Printf("[DEBUG] Collected %d unique parent emails", len(parentEmails))

	// Bulk fetch existing users by email
	existingUsers, err := repositories.GetUsersByEmails(parentEmails)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch existing users by email: %v", err)
		return append(dataStudentErrors, request.ResponseErrorUpload{
			Reason: "Failed to check existing users",
		})
	}
	log.Printf("[DEBUG] Found %d existing users", len(existingUsers))

	// Create maps for existing data and track emails from other schools
	existingUserMap := make(map[string]models.User)
	emailsFromOtherSchools := make(map[string]bool)

	for _, u := range existingUsers {
		if u.UserSchool != nil {
			if u.UserSchool.SchoolID != user.UserSchool.SchoolID {
				// Mark this email as belonging to another school
				emailsFromOtherSchools[u.Email] = true
				continue
			}
			existingUserMap[u.Email] = u
		}
	}

	// Split into create and update operations
	var studentsToCreate, studentsToUpdate []models.Student
	existingStudents, err := studentService.studentRepository.GetStudentsByNIS(getStudentNISList(studentsToProcess), user)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch existing students by NIS: %v", err)
		return append(dataStudentErrors, request.ResponseErrorUpload{
			Reason: "Failed to check existing students",
		})
	}
	log.Printf("[DEBUG] Found %d existing students", len(existingStudents))

	existingStudentMap := utilities.MakeExistingStudentMap(existingStudents)

	// First, update existing students
	for _, student := range studentsToProcess {
		if existing, exists := existingStudentMap[student.Nis]; exists {
			student.ID = existing.ID
			student.FullName = strings.ToUpper(existing.FullName)
			student.BirthPlace = strings.ToUpper(existing.BirthPlace)
			studentsToUpdate = append(studentsToUpdate, student)
			studentHistories = append(studentHistories, models.StudentHistory{
				Master: models.Master{
					CreatedBy: userID,
				},
				StudentID: student.ID,
			})
		} else {
			student.FullName = strings.ToUpper(existing.FullName)
			student.BirthPlace = strings.ToUpper(existing.BirthPlace)
			studentsToCreate = append(studentsToCreate, student)
		}
	}
	log.Printf("[DEBUG] Identified %d students to update and %d to create",
		len(studentsToUpdate), len(studentsToCreate))

	// Update existing students first
	if len(studentsToUpdate) > 0 {
		if err := studentService.studentRepository.BulkUpdateStudents(studentsToUpdate); err != nil {
			log.Printf("[ERROR] Failed to update students in bulk: %v", err)
			for _, student := range studentsToUpdate {
				dataStudentErrors = append(dataStudentErrors, request.ResponseErrorUpload{
					Nis:         student.Nis,
					StudentName: student.FullName,
					Reason:      fmt.Sprintf("Failed to update student: %v", err),
				})
			}
		} else {
			log.Printf("[INFO] Successfully updated %d existing students", len(studentsToUpdate))
		}
	}

	// Create new students in bulk
	if len(studentsToCreate) > 0 {
		if err := studentService.studentRepository.BulkCreateStudents(studentsToCreate); err != nil {
			log.Printf("[ERROR] Failed to create students in bulk: %v", err)
			for _, student := range studentsToCreate {
				dataStudentErrors = append(dataStudentErrors, request.ResponseErrorUpload{
					Nis:         student.Nis,
					StudentName: student.FullName,
					Reason:      fmt.Sprintf("Failed to create student: %v", err),
				})
			}
		} else {
			log.Printf("[INFO] Successfully created %d new students", len(studentsToCreate))
		}
	}

	// Create new users with transaction
	emailToUserIDMap := make(map[string]uint)
	if len(parentEmails) > 0 {
		var userSchools []models.UserSchool
		for _, student := range studentsToProcess {
			if student.Email != "" && existingUserMap[student.Email].ID == 0 {
				username, err := utilities.GenerateRandomUsername(5)
				if err != nil {
					log.Printf("[ERROR] Failed to generate username for email %s: %v",
						student.Email, err)
					continue
				}

				newUser := models.User{
					Master: models.Master{
						CreatedBy: userID,
					},
					Username: username,
					Email:    student.Email,
					RoleID:   uint(2),
				}
				usersToCreate = append(usersToCreate, newUser)

				userSchool := models.UserSchool{
					SchoolID: user.UserSchool.SchoolID,
				}
				userSchools = append(userSchools, userSchool)
			}
		}
		log.Printf("[DEBUG] Prepared %d new users to create", len(usersToCreate))

		if len(usersToCreate) > 0 {
			createdUsers, errorResponses, err := studentService.userRepository.BulkValidateAndCreateUsers(usersToCreate, userSchools)
			log.Printf("[DEBUG] Bulk create users response - Created: %d, Errors: %d",
				len(createdUsers), len(errorResponses))

			if err != nil {
				log.Printf("[ERROR] Failed to create users in bulk: %v", err)
			} else {
				log.Printf("[INFO] Successfully created %d new users", len(createdUsers))

				// Prepare verification records for bulk insert
				var verificationLinks []*models.TempVerificationEmail
				var emailTasks []EmailTask

				for _, u := range createdUsers {
					emailToUserIDMap[u.Email] = u.ID
					existingUserMap[u.Email] = u

					// Create verification record
					verif := &models.TempVerificationEmail{
						UserID: u.ID,
						Type:   nil,
					}
					verificationLinks = append(verificationLinks, verif)

					// Add to email tasks
					emailTasks = append(emailTasks, EmailTask{
						User:   u,
						School: *user.UserSchool.School,
						Verif:  *verif,
					})
				}

				// Process email tasks asynchronously
				if len(emailTasks) > 0 {
					go func(tasks []EmailTask) {
						workers := 2
						emailQueue := NewEmailQueue(workers)
						emailQueue.ProcessEmails(tasks)
					}(emailTasks)
				}
			}
		}
	}

	// Create user-student relationships for both existing and new users
	var pairsToCheck []struct {
		UserID    uint
		StudentID uint
	}

	for _, student := range studentsToProcess {
		if student.Email == "" {
			continue
		}

		// Check if this email belongs to another school
		if emailsFromOtherSchools[student.Email] {
			dataStudentErrors = append(dataStudentErrors, request.ResponseErrorUpload{
				Nis:         student.Nis,
				StudentName: student.FullName,
				Reason:      "Email Terdaftar di tempat lain, silahkan pakai email yang berbeda.",
			})
			continue
		}

		var userID uint
		if u, exists := existingUserMap[student.Email]; exists {
			userID = u.ID
		} else if id, exists := emailToUserIDMap[student.Email]; exists {
			userID = id
		}

		if userID == 0 {
			log.Printf("[WARN] No user found for email: %s", student.Email)
			continue
		}

		var studentID uint
		if existing, exists := existingStudentMap[student.Nis]; exists {
			studentID = existing.ID
		} else {
			for _, newStudent := range studentsToCreate {
				if newStudent.Nis == student.Nis {
					studentID = newStudent.ID
					break
				}
			}
		}

		if studentID == 0 {
			log.Printf("[WARN] No student ID found for NIS: %s", student.Nis)
			continue
		}

		pairsToCheck = append(pairsToCheck, struct {
			UserID    uint
			StudentID uint
		}{
			UserID:    userID,
			StudentID: studentID,
		})
	}

	// Bulk check existing relationships
	existingRelations, err := studentService.studentRepository.BulkCheckUserStudentExists(pairsToCheck)
	if err != nil {
		log.Printf("[ERROR] Failed to check existing relationships: %v", err)
		return dataStudentErrors
	}

	// Create new relationships for non-existing pairs
	for _, pair := range pairsToCheck {
		key := fmt.Sprintf("%d-%d", pair.UserID, pair.StudentID)
		if !existingRelations[key] {
			userStudentPairs = append(userStudentPairs, models.UserStudent{
				Master: models.Master{
					CreatedBy: int(pair.UserID),
				},
				UserID:    pair.UserID,
				StudentID: pair.StudentID,
			})
		}
	}

	// Bulk create user-student relationships with error handling
	if len(userStudentPairs) > 0 {
		if err := studentService.studentRepository.BulkCreateUserStudents(userStudentPairs); err != nil {
			log.Printf("[ERROR] Failed to create user-student relationships: %v", err)
			// Consider adding these failures to dataStudentErrors
			for _, pair := range userStudentPairs {
				dataStudentErrors = append(dataStudentErrors, request.ResponseErrorUpload{
					Reason: fmt.Sprintf("Failed to create user-student relationship for UserID: %d, StudentID: %d", pair.UserID, pair.StudentID),
				})
			}
		} else {
			log.Printf("[INFO] Successfully created %d user-student relationships", len(userStudentPairs))
		}
	}

	// Bulk create all histories
	if len(studentHistories) > 0 {
		if err := studentService.studentRepository.BulkCreateStudentHistory(studentHistories); err != nil {
			fmt.Printf("Error creating student histories: %v\n", err)
		} else {
			fmt.Printf("Successfully created %d student history records\n", len(studentHistories))
		}
	}

	if len(dataStudentErrors) > 0 {
		fmt.Printf("Total errors encountered: %d\n", len(dataStudentErrors))
	}

	log.Printf("[DEBUG] Batch processing completed. Total errors: %d", len(dataStudentErrors))
	return dataStudentErrors
}

func getStudentNISList(students []models.Student) []string {
	nisList := make([]string, len(students))
	for i, student := range students {
		nisList[i] = student.Nis
	}
	return nisList
}

func (studentService *StudentService) CreateImageStudentService(request request.CreateStudentImageRequest, id uint, userID int) (*models.Student, error) {
	user, err := studentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	student, err := studentService.studentRepository.GetStudentByID(id, user)
	if err != nil {
		return nil, fmt.Errorf("Student not found")
	}
	student.Image = request.Image
	student.Master.UpdatedBy = userID
	data, err := studentService.studentRepository.CreateImageStudentRepository(&student, int(id))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func validateDuplicateNis(nis string, user models.User, studentService *StudentService, excludeID ...uint) error {
	if nis == "" {
		return nil
	}

	student, err := studentService.studentRepository.GetStudentByNis(nis, user)
	if err == nil {
		// If we're updating and the found student is the same as the one being updated, it's OK
		if len(excludeID) > 0 && student.ID == excludeID[0] {
			return nil
		}
		return fmt.Errorf("NIS already exists")
	}

	return nil
}

type ExcelExportData struct {
	*bytes.Buffer
	SchoolName string
}

func (studentService *StudentService) ExportStudentToExcel(search, searchNis string, gradeId, yearId int, status string, schoolId, userID int, isActive *bool) (*ExcelExportData, error) {
	// Get student data using existing GetAllStudent method
	studentData, err := studentService.GetAllStudent(0, 0, search, userID, status, gradeId, yearId, schoolId, searchNis, 0, "", "", 0, isActive)
	if err != nil {
		return nil, err
	}

	// Get school name
	schoolName := "data" // Default name
	if len(studentData.Data) > 0 && studentData.Data[0].SchoolName != "" {
		schoolName = studentData.Data[0].SchoolName
	}

	// Create Excel utility instance
	excelUtil := utilities.NewExcelUtility()
	defer excelUtil.Close()

	sheetName := "Student Data"
	excelUtil.File.SetSheetName("Sheet1", sheetName)

	// Create styles
	boldStyle, err := excelUtil.CreateBoldStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create bold style: %v", err)
	}

	centerStyle, err := excelUtil.CreateCenterStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create center style: %v", err)
	}

	// Define headers
	headers := []string{
		"No.", "NIS", "Nama Siswa",
		"Tanggal Lahir", "Status",
		"Tahun Ajaran",
	}

	if err := excelUtil.WriteHeaders(sheetName, headers, true); err != nil {
		return nil, fmt.Errorf("failed to write headers: %v", err)
	}

	// Write headers with bold style
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", strconv.Itoa('A'+i))
		excelUtil.SetCellValue(sheetName, cell, header)
		excelUtil.SetCellStyle(sheetName, cell, cell, boldStyle)
	}

	// Fill data
	for i, student := range studentData.Data {
		row := i + 2

		// Add index number with center alignment
		noCell := fmt.Sprintf("A%d", row)
		excelUtil.SetCellValue(sheetName, noCell, i+1)
		excelUtil.SetCellStyle(sheetName, noCell, noCell, centerStyle)

		// Set other cell values
		excelUtil.SetCellValue(sheetName, fmt.Sprintf("B%d", row), student.Nis)
		excelUtil.SetCellValue(sheetName, fmt.Sprintf("C%d", row), student.FullName)
		if student.BirthDate != nil {
			formattedDate := student.BirthDate.Format(constants.DateFormatYYYYMMDD)
			excelUtil.SetCellValue(sheetName, fmt.Sprintf("D%d", row), formattedDate)
		} else {
			excelUtil.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
		}
		excelUtil.SetCellValue(sheetName, fmt.Sprintf("E%d", row), student.Status)

		// Get school year name
		schoolYearName := "-"
		if student.SchoolYearID != 0 {
			schoolYear, err := studentService.schoolYearRepository.GetSchoolYearByID(student.SchoolYearID)
			if err == nil {
				schoolYearName = schoolYear.SchoolYearName
			}
		}

		// Set the school year name in the Excel cell
		excelUtil.SetCellValue(sheetName, fmt.Sprintf("F%d", row), schoolYearName)
	}

	// // Set column widths
	// columnWidths := map[string]float64{
	// 	"A": 5,  // No.
	// 	"B": 15, // NIS
	// 	"C": 30, // Nama Siswa
	// 	"D": 15, // Tanggal Lahir
	// 	"E": 15, // Status
	// 	"F": 15, // Tahun Ajaran
	// }
	// if err := excelUtil.SetColumnWidths(sheetName, columnWidths); err != nil {
	// 	return nil, fmt.Errorf("failed to set column widths: %v", err)
	// }

	// Adjust column widths automatically
	columns := []string{"A", "B", "C", "D", "E", "F"}
	for _, col := range columns {
		if err := excelUtil.AutoFitColumn(sheetName, col); err != nil {
			return nil, fmt.Errorf("failed to auto-fit column %s: %v", col, err)
		}
	}

	// Add auto filter
	lastRow := len(studentData.Data) + 1
	if err := excelUtil.AddAutoFilter(sheetName, "A1", fmt.Sprintf("F%d", lastRow)); err != nil {
		return nil, fmt.Errorf("failed to add auto filter: %v", err)
	}

	// Write to buffer
	buffer, err := excelUtil.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %v", err)
	}

	return &ExcelExportData{
		Buffer:     buffer,
		SchoolName: schoolName,
	}, nil
}

func (studentService *StudentService) validateAndPrepareStudentsBulk(rows [][]string, userID int, user models.User) ([]models.Student, []request.ResponseErrorUpload) {
	var studentsToProcess []models.Student
	var errors []request.ResponseErrorUpload

	// Collect all unique values for batch lookups
	schoolGradeNames := make(map[string]bool)
	schoolClassNames := make(map[string]bool)
	schoolYearNames := make(map[string]bool)
	nisNumbers := make(map[string]bool)

	// First pass - collect all unique values
	for _, row := range rows {
		// Skip dummy data
		if row[0] == "442342424" || row[9] == "08563142342" || row[10] == constants.DummyEmail {
			continue
		}

		if len(row) > 7 {
			schoolGradeNames[row[7]] = true
		}
		if len(row) > 8 {
			schoolClassNames[row[8]] = true
		}
		if len(row) > 13 {
			schoolYearNames[row[13]] = true
		}
		if len(row) > 0 {
			nisNumbers[row[0]] = true
		}
	}

	// Batch fetch all required data
	schoolGrades, err := repositories.GetAllSchoolGradesBulk(utilities.GetMapKeys(schoolGradeNames))
	if err != nil {
		return nil, []request.ResponseErrorUpload{{Reason: "Failed to fetch school grades"}}
	}

	schoolClasses, err := repositories.GetAllSchoolClassesBulk(utilities.GetMapKeys(schoolClassNames), user)
	if err != nil {
		return nil, []request.ResponseErrorUpload{{Reason: "Failed to fetch school classes"}}
	}

	schoolYears, err := studentService.schoolYearRepository.GetAllSchoolYearsBulk(utilities.GetMapKeys(schoolYearNames), user.UserSchool.SchoolID)
	if err != nil {
		return nil, []request.ResponseErrorUpload{{Reason: "Failed to fetch school years"}}
	}

	existingStudents, err := studentService.studentRepository.GetStudentsByNIS(utilities.GetMapKeys(nisNumbers), user)
	if err != nil {
		return nil, []request.ResponseErrorUpload{{Reason: "Failed to fetch existing students"}}
	}

	// Create lookup maps
	schoolGradeMap := utilities.MakeSchoolGradeMap(schoolGrades)
	schoolClassMap := utilities.MakeSchoolClassMap(schoolClasses)
	schoolYearMap := utilities.MakeSchoolYearMap(schoolYears)
	existingStudentMap := utilities.MakeExistingStudentMap(existingStudents)

	// Process all rows in bulk
	for _, row := range rows {
		student, err := validateAndPrepareStudentRow(
			row,
			userID,
			user,
			schoolGradeMap,
			schoolClassMap,
			schoolYearMap,
			existingStudentMap,
		)

		if err != nil {
			errors = append(errors, request.ResponseErrorUpload{
				Nis:         row[0],
				StudentName: row[1],
				Reason:      err.Error(),
			})
			continue
		}

		studentsToProcess = append(studentsToProcess, student)
	}

	return studentsToProcess, errors
}

// Helper function to validate and prepare a single row with pre-fetched data
func validateAndPrepareStudentRow(
	row []string,
	userID int,
	user models.User,
	schoolGradeMap map[string]uint,
	schoolClassMap map[string]uint,
	schoolYearMap map[string]uint,
	existingStudentMap map[string]models.Student,
) (models.Student, error) {
	var student models.Student

	// Check for dummy data first
	if row[0] == "442342424" {
		return student, fmt.Errorf("NIS masih sama dengan contoh")
	}
	if row[9] == "08563142342" {
		return student, fmt.Errorf("No HP masih sama dengan contoh")
	}
	if row[10] == constants.DummyEmail {
		return student, fmt.Errorf("Email masih sama dengan contoh")
	}

	// Basic validation
	if user.UserSchool.School.SchoolName != row[12] {
		return student, fmt.Errorf("Nama Sekolah Tidak sama")
	}

	// Format phone number
	phoneNumberFormat := utilities.FormatPhoneNumber(row[9])

	// Parse birth date
	birthDate, err := utilities.ParseBirthDate(row[3])
	if err != nil {
		return student, fmt.Errorf("Format Tanggal Lahir Salah")
	}

	// Get IDs from pre-fetched maps
	schoolGradeID := schoolGradeMap[row[7]]
	if schoolGradeID == 0 {
		return student, fmt.Errorf("Jenjang Sekolah Tidak Ada")
	}

	schoolClassID := schoolClassMap[row[8]]
	if schoolClassID == 0 {
		return student, fmt.Errorf("Kelas Tidak Ada")
	}

	schoolYearID := schoolYearMap[row[13]]
	if schoolYearID == 0 {
		return student, fmt.Errorf("Data Tahun Ajaran Tidak Ada")
	}

	// Create student object
	studentData := models.Student{
		Nis:           row[0],
		FullName:      row[1],
		BirthPlace:    row[2],
		BirthDate:     &birthDate,
		Gender:        utilities.MapDisplayLabelToValueStudent("gender", row[4]),
		Religion:      utilities.MapDisplayLabelToValueStudent("religion", row[5]),
		Address:       row[6],
		SchoolGradeID: schoolGradeID,
		SchoolClassID: schoolClassID,
		NoHandphone:   phoneNumberFormat,
		Email:         row[10],
		Status:        utilities.MapDisplayLabelToValueStudent("status", row[11]),
		SchoolYearID:  schoolYearID,
	}

	return studentData, nil
}

func (studentService *StudentService) HandleCSVFileStudent(file *multipart.FileHeader, c *fiber.Ctx, userID int) error {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Create a buffered reader for better performance
	reader := bufio.NewReader(src)

	// Create a CSV reader
	csvReader := csv.NewReader(reader)

	// Read all rows
	rows, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Validate minimum rows (header + at least one data row)
	if len(rows) < 2 {
		return fmt.Errorf("file must contain at least one data row")
	}

	// Process rows in batches for better memory management
	const batchSize = 500
	var allErrors []request.ResponseErrorUpload

	for i := 1; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		errors := studentService.processRowsBatch(rows[i:end], userID)
		allErrors = append(allErrors, errors...)
	}

	// Prepare response
	message := "Success upload data"
	if len(allErrors) > 0 {
		message = "Failed upload data"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"data":    allErrors,
	})
}
