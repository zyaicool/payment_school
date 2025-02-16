package services

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"

	utilities "schoolPayment/utilities"
)

type BillingServiceInterface interface {
	GetBillingByID(id uint) (response.BillingDetailResponse, error)
	GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, userID int) (response.BillingListResponse, error)
	CreateBilling(billingRequest *request.BillingCreateRequest, userID int) (*models.Billing, error)
	GetBillingStatuses(filePath string) ([]models.BillingStatus, error)
	CreateDonation(billingName string, schoolGradeId int, bankAccountId int, userId uint) (response.BillingResponse, error)
	GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) (*response.BillingByStudentResponse, error)
}

type BillingService struct {
	billingRepository     repositories.BillingRepositoryInterface
	userRepository        repositories.UserRepository
	schoolClassRepository repositories.SchoolClassRepositoryInterface
	schoolYearRepository  repositories.SchoolYearRepository
	schoolGradeRepository repositories.SchoolGradeRepositoryInterface
	studentRepository	  repositories.StudentRepositoryInteface
}

func NewBillingService(
	billingRepository repositories.BillingRepositoryInterface, 
	userRepository repositories.UserRepository, 
	schoolClassRepository repositories.SchoolClassRepositoryInterface, 
	schoolYearRepository repositories.SchoolYearRepository, 
	schoolGradeRepository repositories.SchoolGradeRepositoryInterface,
	studentRepository	  repositories.StudentRepositoryInteface,
	) BillingServiceInterface {
	return &BillingService{
		billingRepository:     billingRepository,
		userRepository:        userRepository,
		schoolClassRepository: schoolClassRepository,
		schoolYearRepository:  schoolYearRepository,
		schoolGradeRepository: schoolGradeRepository,
		studentRepository: studentRepository,
	}
}

func (billingService *BillingService) GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, userID int) (response.BillingListResponse, error) {
	var mapBilling response.BillingListResponse
	mapBilling.Limit = limit
	mapBilling.Page = page
	mapBilling.TotalData = 0
	mapBilling.TotalPage = 0
	mapBilling.Data = []response.BillingResponse{}

	if sortBy != "" {
		sortBy = billingService.ChangeStringSortBy(sortBy)
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	user, err := billingService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return response.BillingListResponse{}, err
	}

	listBilling, totalPage, totalData, err := billingService.billingRepository.GetAllBilling(page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, isDonation, user)
	if err != nil {
		mapBilling.Data = []response.BillingResponse{}
		return mapBilling, nil
	}

	// Populate the response with the specific fields
	for _, billing := range listBilling {
		bankName, _ := GetBankName(billing.BankAccount.BankName)
		bankAccountName := bankName + " - " + billing.BankAccount.AccountNumber
		schoolGradeName := billing.SchoolGrade.SchoolGradeName
		mapBilling.Data = append(mapBilling.Data, response.BillingResponse{
			ID:              int(billing.ID),
			BillingName:     billing.BillingName,
			BillingType:     billing.BillingType,
			BankAccountName: bankAccountName,
			SchoolGradeName: schoolGradeName,
			CreatedBy:       billing.CreateByUsername, // Ensure `CreatedBy` contains the user name from the join
			CreatedAt:       billing.CreatedAt,
		})
	}

	mapBilling.TotalData = totalData
	mapBilling.TotalPage = totalPage

	return mapBilling, nil
}

func (billingService *BillingService) GetBillingByID(id uint) (response.BillingDetailResponse, error) {
	billing, err := billingService.billingRepository.GetBillingByID(int(id))

	if err != nil {
		return response.BillingDetailResponse{}, err // Handle the error appropriately
	}

	bankName, _ := GetBankName(billing.BankAccount.BankName)
	bankAccountName := bankName + " - " + billing.BankAccount.AccountNumber

	// Prepare the response struct
	billingResponse := response.BillingDetailResponse{
		ID:              billing.ID,
		BillingName:     billing.BillingName,
		BillingCode:     billing.BillingCode,
		BillingType:     billing.BillingType,
		BankAccountName: bankAccountName,
		Description:     billing.Description,
		SchoolYear:      billing.SchoolYear.SchoolYearName, // Assuming SchoolYear has a Name field
		SchoolClassList: billing.SchoolClassIds,
		DetailBillings:  []response.DetailBilling{}, // This will be populated later
	}

	// Populate detail billings if you have a method for that
	detailBillings, err := billingService.billingRepository.GetDetailBillingsByBillingID(billing.ID)
	if err == nil {
		// Map the detailBillings to the response type
		billingResponse.DetailBillings = detailBillings
	}

	// Retrieve school class names
	schoolClassIDs := strings.Split(billing.SchoolClassIds, ",") // Assuming this is a comma-separated string
	var schoolClassNames []string
	for _, idStr := range schoolClassIDs {
		id, err := strconv.ParseUint(idStr, 10, 32) // Parse string to uint
		if err != nil {
			continue // Skip if there's an error
		}

		schoolClass, err := billingService.schoolClassRepository.GetSchoolClassByID(uint(id)) // Get the class by ID
		if err == nil {
			schoolClassNames = append(schoolClassNames, schoolClass.SchoolClassName) // Assume ClassName is the field containing the name
		}
	}

	// Join the class names into a single string
	billingResponse.SchoolClassList = strings.Join(schoolClassNames, ", ")

	return billingResponse, nil
}

func (billingService *BillingService) CreateBilling(billingRequest *request.BillingCreateRequest, userID int) (*models.Billing, error) {
	if err := utilities.ValidateBillingName(billingRequest.BillingName); err != nil {
		return nil, err
	}

	if err := utilities.ValidateFieldNotEmpty(billingRequest.BillingCode, "Biling Code"); err != nil {
		return nil, err
	}

	if err := utilities.ValidateFieldNotEmpty(billingRequest.BillingType, "Biling Type"); err != nil {
		return nil, err
	}

	if err := utilities.ValidateFieldNotEmpty(billingRequest.SchoolGradeID, "School Grade"); err != nil {
		return nil, err
	}

	if err := utilities.ValidateFieldNotEmpty(billingRequest.SchoolYearId, "School Year"); err != nil {
		return nil, err
	}

	// Add validation for detailBillings amount
	for _, detail := range billingRequest.DetailBillings {
		if detail.Amount <= 0 {
			return nil, fmt.Errorf("Detail billing amount tidak valid")
		}
	}

	// get data school grade
	schoolGrade, err := billingService.schoolGradeRepository.GetSchoolGradeByID(uint(billingRequest.SchoolGradeID))
	if err != nil {
		fmt.Printf("Error fetching school grade: %v", err)
		return nil, err
	}

	// get data school year
	schoolYear, err := billingService.schoolYearRepository.GetSchoolYearByID(uint(billingRequest.SchoolYearId))
	if err != nil {
		return nil, err
	}

	billingNumber, err := billingService.GenerateBillingNumber(billingRequest.BillingType, schoolGrade.SchoolGradeCode, schoolYear.SchoolYearName)
	if err != nil {
		return nil, err
	}

	user, err := billingService.userRepository.GetUserByID(uint(userID)) 
	if err != nil {
		return nil, err
	}

	if user.UserSchool == nil {
		return nil, fmt.Errorf("Silahkan add user ke data sekolah")
	}

	ok := billingService.billingRepository.CheckBillingCode(billingRequest.BillingCode)
	if !ok {
		return nil, fmt.Errorf("Billing code sudah di gunakan")
	}

	schoolClassIdsStr := strings.Join(billingRequest.SchoolClassIds, ",")

	billing := models.Billing{
		BillingNumber:  billingNumber,
		BillingName:    billingRequest.BillingName,
		BillingType:    billingRequest.BillingType,
		SchoolGradeID:  uint(billingRequest.SchoolGradeID),
		SchoolYearId:   uint(billingRequest.SchoolYearId),
		BillingAmount:  billingRequest.BillingAmount,
		Description:    billingRequest.Description,
		BillingCode:    billingRequest.BillingCode,
		BankAccountId:  billingRequest.BankAccountId,
		SchoolClassIds: schoolClassIdsStr,
	}
	billing.Master.CreatedBy = userID
	billing.Master.UpdatedBy = userID

	dataBilling, err := billingService.billingRepository.CreateBilling(&billing)
	if err != nil {
		return nil, err
	}

	err = billingService.GenerateBillingToStudent(user, dataBilling, billingRequest.SchoolClassIds, billingRequest.DetailBillings)

	return dataBilling, nil
}

func (billingService *BillingService) GenerateBillingToStudent(user models.User, billing *models.Billing, schoolClassIds []string, detailBillings []request.DetailBillings) error {
	var classIdsInt []int
	for _, id := range schoolClassIds {
		// Convert string to int
		intId, err := strconv.Atoi(id)
		if err != nil {
			// Handle error (e.g., log it or skip the invalid entry)
			fmt.Printf("Error converting %s to int: %v\n", id, err)
			continue
		}
		classIdsInt = append(classIdsInt, intId)
	}

	listStudent, err := billingService.studentRepository.GetAllStudentForBilling(user, int(billing.SchoolGradeID), classIdsInt)
	if err != nil {
		return err
	}

	for _, detail := range detailBillings {
		dueDate, err := time.Parse("2006-01-02", detail.DueDate)
		if err != nil {
			return fmt.Errorf("error parsing due date: %w", err)
		}

		billingDetailId, err := SaveDataBillingDetail(user, billing, detail, &dueDate)
		if err != nil {
			return err
		}

		for _, student := range listStudent {
			// Save billing data for each billing detail and student
			err = SaveDataBillingStudent(user, student, billing, &dueDate, detail, billingDetailId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateBilling(billingService *BillingService, id uint, billingRequest *request.BillingUpdateRequest, userID int) (*models.Billing, error) {
	getBilling, err := billingService.billingRepository.GetBillingByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getBilling.BillingName = billingRequest.BillingName
	getBilling.UpdatedBy = userID

	dataBilling, err := repositories.UpdateBilling(&getBilling)
	if err != nil {
		return nil, err
	}

	return dataBilling, nil
}

func DeleteBilling(billingService *BillingService, id uint, userID int) (*models.Billing, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getBilling, err := billingService.billingRepository.GetBillingByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getBilling.Master.DeletedAt = currentTimePointer
	getBilling.Master.DeletedBy = &userID
	dataBilling, err := repositories.UpdateBilling(&getBilling)
	if err != nil {
		return nil, err
	}

	return dataBilling, nil
}

func GenerateInstallment(id uint) ([]int, error) {
	var totalInstallment []int
	installment := 0
	i := 1

	billingType, err := repositories.GetBillingTypeByID(id)
	if err != nil {
		return totalInstallment, err
	}

	if strings.ToLower(billingType.BillingTypePeriod) != "bulanan" {
		if strings.ToLower(billingType.BillingTypePeriod) == "triwulan" {
			installment = 12 / 3
		} else {
			installment = 12 / 6
		}

		for {

			if i != 1 {
				totalInstallment = append(totalInstallment, i)
			}

			if i == installment {
				break
			}
			i++
		}

	}

	return totalInstallment, nil
}

func (billingService *BillingService) GenerateBillingNumber(billingType string, schoolGrade string, schoolYear string) (string, error) {
	fmt.Printf("BillingType: %s, SchoolGradeCode: %s, SchoolYearName: %s\n", billingType, schoolGrade, schoolYear)
	lastNumber, err := billingService.billingRepository.GetLastSequenceNumberBilling()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1

	// get date now
	currentTime := time.Now()
	currentYear := currentTime.Format("06") // get year 2 digit
	monthDay := currentTime.Format("0102")  // get month and day

	// format school year kasih log di sini
	years := strings.Split(schoolYear, "/")
	firstYear := years[0][2:]
	secondYear := years[1][2:]
	resultSchoolYear := firstYear + secondYear

	newCode := fmt.Sprintf("INV/"+billingType+"/"+schoolGrade+"/"+resultSchoolYear+"/"+currentYear+"/"+monthDay+"/%05d", newSequence)

	return newCode, nil
}

func GenerateBillingStudentInstallment(user models.User, student models.Student, billing *models.Billing, intTenor int, period int) error {
	i := 1
	for {
		nextMonth := time.Now().AddDate(0, i*period, 0) // Add one month
		dueDate := time.Date(nextMonth.Year(), nextMonth.Month(), 10, 0, 0, 0, 0, nextMonth.Location())
		err := SaveDataBillingStudent(user, student, billing, &dueDate, request.DetailBillings{}, 0)
		if err != nil {
			return err
		}

		if i == intTenor {
			break
		}
		i++
	}

	return nil
}

func SaveDataBillingStudent(user models.User, student models.Student, billing *models.Billing, dueDate *time.Time, detailBillings request.DetailBillings, billingDetailId uint) error {
	belumBayarCode, _ := BillingStatusCode("Belum bayar")
	strBelumBayarCode := strconv.Itoa(belumBayarCode)
	billingStudent := models.BillingStudent{
		BillingID:         billing.ID,
		StudentID:         student.ID,
		PaymentStatus:     strBelumBayarCode,
		DetailBillingName: detailBillings.DetailBillingName,
		Amount:            detailBillings.Amount,
		DueDate:           dueDate,
		BillingDetailID:   billingDetailId,
	}
	billingStudent.Master.CreatedBy = int(user.ID)
	billingStudent.Master.UpdatedBy = int(user.ID)

	_, err := repositories.CreateBillingStudent(&billingStudent, "", true)
	if err != nil {
		return err
	}

	return nil
}

func SaveDataBillingDetail(user models.User, billing *models.Billing, detailBillings request.DetailBillings, dueDate *time.Time) (uint, error) {
	billingDetail := models.BillingDetail{
		BillingID:         billing.ID,
		DetailBillingName: detailBillings.DetailBillingName,
		Amount:            detailBillings.Amount,
		DueDate:           dueDate,
	}
	billingDetail.Master.CreatedBy = int(user.ID)
	billingDetail.Master.UpdatedBy = int(user.ID)

	data, err := repositories.CreateBillingDetail(&billingDetail, "", true)
	if err != nil {
		return uint(0), err
	}

	return data.ID, nil
}

func ExtractTenorNumber(tenor string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(tenor)

	if match == "" {
		return 0, fmt.Errorf("no number found in tenor")
	}

	number, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}

	return number, nil
}

const (
	BelumBayarCode = 1
	LunasCode      = 2
)

// GetPaymentStatuses reads billing_status.json and returns a list of payment statuses
func (billingService *BillingService) GetBillingStatuses(filePath string) ([]models.BillingStatus, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// // Read the file contents
	// data, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	return nil, fmt.Errorf("error reading file: %v", err)
	// }

	// Parse JSON data into a slice of BillingStatus
	var billingStatuses []models.BillingStatus
	if err := json.NewDecoder(file).Decode(&billingStatuses); err != nil {
		return nil, err
	}

	// // Parse JSON data into a slice of BillingStatus
	// var billingStatuses []models.BillingStatus
	// err = json.Unmarshal(data, &billingStatuses)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing JSON: %v", err)
	// }

	return billingStatuses, nil
}

func BillingStatusCode(status string) (int, error) {
	// switch status {
	// case "Belum bayar":
	// 	return BelumBayarCode, nil
	// case "Lunas":
	// 	return LunasCode, nil
	// default:
	// 	return 0, fmt.Errorf("invalid billing status: %s", status)
	// }

	billingStatuses, err := models.LoadBillingStatus("data/billing_status.json") // Path to your JSON file
	if err != nil {
		return 0, fmt.Errorf("invalid billing status: %s", status)
	}

	for _, billingStatus := range billingStatuses {
		if billingStatus.Name == status {
			return billingStatus.Id, nil
		}
	}
	return 0, fmt.Errorf("invalid billing status: %s", status)
}

func (billingService *BillingService) ChangeStringSortBy(sortBy string) string {
	switch sortBy {
	case "bankAccountName":
		sortBy = "bankAccountId"
	default:
	}
	return sortBy
}

func (billingService *BillingService) GenerateBillingNumberDonation(billingType string, schoolGrade string, schoolYear string) (string, error) {
	lastSequence, err := billingService.billingRepository.GetLastSequenceNumberBilling()
	if err != nil {
		return "", err
	}

	newSequence := lastSequence + 1
	currentTime := time.Now()
	currentYear := currentTime.Format("06")
	monthDay := currentTime.Format("0102")

	years := strings.Split(schoolYear, "/")
	firstYear := years[0][2:]
	secondYear := years[1][2:]
	resultSchoolYear := firstYear + secondYear

	// newCode := fmt.Sprintf("INV/%s/%s/%s/%s/%s/%05d", billingType, schoolGrade, resultSchoolYear, currentYear, monthDay, newSequence)
	newCode := fmt.Sprintf("INV/"+billingType+"/"+resultSchoolYear+"/"+currentYear+"/"+monthDay+"/%05d", newSequence)

	return newCode, nil
}

func (billingService *BillingService) CreateDonation(billingName string, schoolGradeId int, bankAccountId int, userId uint) (response.BillingResponse, error) {
	// Dapatkan tahun ajaran terbaru

	user, err := billingService.userRepository.GetUserByID(userId)
	if err != nil {
		return response.BillingResponse{}, fmt.Errorf("Failed to get user: %w", err)
	}

	year := strconv.Itoa(time.Now().Year())
	schoolYears, _, _, err := billingService.schoolYearRepository.GetAllSchoolYear(1, 0, year, "", "", user.UserSchool.SchoolID)
	if err != nil || len(schoolYears) == 0 {
		return response.BillingResponse{}, fmt.Errorf("Failed to get latest school year")
	}
	latestSchoolYearID := schoolYears[0].ID
	schoolYear := schoolYears[0].SchoolYearName

	schoolGrade := strconv.Itoa(schoolGradeId)

	// Generate billing number
	billingNumber, err := billingService.GenerateBillingNumberDonation("Donation", schoolGrade, schoolYear)
	if err != nil {
		return response.BillingResponse{}, fmt.Errorf("Failed to generate billing number: %s", err)
	}

	// Buat billing baru
	newBilling := models.Billing{
		BillingName:   billingName,
		BillingNumber: billingNumber,
		SchoolGradeID: uint(schoolGradeId),
		BankAccountId: bankAccountId,
		IsDonation:    true,
		SchoolYearId:  uint(latestSchoolYearID),
		BillingType:   "non-regular",
	}

	newBilling.CreatedBy = int(userId)

	// Simpan billing baru ke database
	savedBilling, err := billingService.billingRepository.CreateBillingDonation(&newBilling)
	if err != nil {
		return response.BillingResponse{}, err
	}

	// Buat response
	billingResponse := response.BillingResponse{
		ID:              int(savedBilling.ID),
		BillingName:     savedBilling.BillingName,
		BillingType:     "Donation",
		BillingNumber:   savedBilling.BillingNumber,
		BankAccountName: "Bank Name - " + strconv.Itoa(savedBilling.BankAccountId),
		CreatedBy:       "admin",
		CreatedAt:       savedBilling.CreatedAt,
	}

	return billingResponse, nil
}

func (billingService *BillingService) GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) (*response.BillingByStudentResponse, error) {
	// Get billings with filters
	billings, err := billingService.billingRepository.GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID)
	if err != nil {
		return nil, err
	}

	res := &response.BillingByStudentResponse{
		BillingExist: false,
		Data:         make([]response.BillingStudentData, 0),
	}

	hasExistingBilling := false

	// Create a map to group billing details by billing ID
	billingMap := make(map[uint]response.BillingStudentData)

	// Process each billing
	for _, billing := range billings {
		// Check if we already have this billing in our map
		billingData, exists := billingMap[billing.ID]
		if !exists {
			billingData = response.BillingStudentData{
				BillingID:       billing.ID,
				BillingName:     billing.BillingName,
				BillingStudents: make([]response.BillingStudentDetailData, 0),
			}
		}

		detailData := response.BillingStudentDetailData{
			BillingDetailID:   billing.BillingDetailID,
			BillingDetailName: billing.DetailBillingName,
			Amount:            billing.Amount,
			IsExist:           billing.IsExist,
			Disabled:          billing.Disabled,
		}

		// Update hasExistingBilling if any billing exists
		if billing.IsExist {
			hasExistingBilling = true
		}

		billingData.BillingStudents = append(billingData.BillingStudents, detailData)
		billingMap[billing.ID] = billingData
	}

	// Convert map to slice for response
	for _, billingData := range billingMap {
		res.Data = append(res.Data, billingData)
	}

	res.BillingExist = hasExistingBilling

	return res, nil
}
