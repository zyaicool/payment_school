package services

import (
	"fmt"
	"strings"
	"time"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	repositories "schoolPayment/repositories"
	utilities "schoolPayment/utilities"
)

type BillingStudentServiceInterface interface {
	GetAllBillingStudent(page int, limit int, search string, studentID int, roleID int, schoolGradeId int, paymentType string, schoolID int, userID int, sortBy string, sortOrder string, schoolClassID int) (response.ListBillingPerStudent, error)
	GetDetailBillingIDService(studentId int, user models.User) (*response.BillingStudentByStudentIDResponse, error)
	GetDetailBillingStudentByID(billingStudentId int) (*response.BillingStudentDetailResponse, error)
	UpdateBillingStudentService(billingStudentId int, updateRequest request.UpdateBillingStudentRequest) (response.BillingStudentDetailResponse, error)
	DeleteBillingStudentService(billingStudentId, deletedBy int) error
	CreateBillingStudent(request request.CreateBillingStudentRequest, userID int) ([]models.BillingStudent, error)
}

type BillingStudentService struct {
	billingStudentRepository repositories.BillingStudentRepository
	userRepository           repositories.UserRepository
	schoolYearRepository     repositories.SchoolYearRepository
	schoolClassRepository    repositories.SchoolClassRepositoryInterface
	billingRepository        repositories.BillingRepositoryInterface
	schoolGradeRepository    repositories.SchoolGradeRepositoryInterface
	studentRepository		 repositories.StudentRepositoryInteface
}

func NewBillingStudentService(
	billingStudentRepository repositories.BillingStudentRepository, 
	userRepository repositories.UserRepository, 
	schoolYearRepository repositories.SchoolYearRepository, 
	schoolClassRepository repositories.SchoolClassRepositoryInterface, 
	billingRepository repositories.BillingRepositoryInterface, 
	schoolGradeRepository repositories.SchoolGradeRepositoryInterface,
	studentRepository repositories.StudentRepositoryInteface,
	) BillingStudentServiceInterface {
	return &BillingStudentService{
		billingStudentRepository: billingStudentRepository, 
		userRepository: userRepository, schoolYearRepository: 
		schoolYearRepository, 
		schoolClassRepository: schoolClassRepository, 
		billingRepository: billingRepository,
		schoolGradeRepository: schoolGradeRepository,
		studentRepository: studentRepository,
	}
}

func (billingStudentService *BillingStudentService) GetAllBillingStudent(page int, limit int, search string, studentID int, roleID int, schoolGradeId int, paymentType string, schoolID int, userID int, sortBy string, sortOrder string, schoolClassID int) (response.ListBillingPerStudent, error) {
	var listBilling response.ListBillingPerStudent

	listBilling.Page = page
	listBilling.Limit = limit
	listBilling.TotalData = 0
	listBilling.TotalPage = 0

	user, err := billingStudentService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return listBilling, err
	}
	if roleID == 2 && studentID == 0 {
		return listBilling, fmt.Errorf("please select student")
	} else if roleID == 2 && studentID != 0 {
		validateStudentId := false
		for _, userStudent := range user.UserStudents {
			if userStudent.StudentID == uint(studentID) {
				validateStudentId = true
			}
		}
		if !validateStudentId {
			return listBilling, fmt.Errorf("the student ID is not associated with the user")
		}
	}

	if roleID != 2 {
		if user.UserSchool != nil {
			schoolID = int(user.UserSchool.SchoolID)
		}
	}

	if sortBy != "" {
		sortBy = utilities.ChangeStringSortByBillingStudent(sortBy)
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	listBillingPerStudent, totalPage, totalData, err := billingStudentService.billingStudentRepository.GetAllBillingStudent(page, limit, search, studentID, roleID, schoolGradeId, paymentType, schoolID, sortBy, sortOrder, schoolClassID)
	if err != nil {
		listBilling.Data = []response.DataListBillingPerStudent{}
		return listBilling, err
	}

	if len(listBillingPerStudent) > 0 {
		listBilling.Data = listBillingPerStudent
	} else {
		listBilling.Data = []response.DataListBillingPerStudent{}
	}

	listBilling.TotalData = totalData
	listBilling.TotalPage = totalPage

	if len(listBillingPerStudent) > 0 && roleID == 2 {
		detailBillingStudentResponses, err := repositories.GetFirstBilling(studentID)
		if err == nil {
			schoolYear, err := repositories.GetLatestSchoolYear()
			if err != nil {
				return listBilling, err
			}

			years := strings.Split(schoolYear.SchoolYearName, "/")
			for i := range detailBillingStudentResponses {
				detailBillingStudentResponses[i].Semester = years[1]
			}
		}
	}

	return listBilling, nil
}

func (billingStudentService *BillingStudentService) GetDetailBillingIDService(studentId int, user models.User) (*response.BillingStudentByStudentIDResponse, error) {
	schoolGradeName := "-"
	schoolClassName := "-"

	student, err := billingStudentService.studentRepository.GetStudentByID(uint(studentId), user)
	if err != nil {
		return nil, err
	}

	years, _ := billingStudentService.schoolYearRepository.GetSchoolYearByID(student.SchoolYearID)
	latestSchoolYear := years.SchoolYearName

	getSchoolGrade, err := billingStudentService.schoolGradeRepository.GetSchoolGradeByID(student.SchoolGradeID)
	if err == nil {
		schoolGradeName = getSchoolGrade.SchoolGradeName
	}

	// Get School Class
	getSchoolClass, err := billingStudentService.schoolClassRepository.GetSchoolClassByID(student.SchoolClassID)
	if err == nil {
		schoolClassName = getSchoolClass.SchoolClassName
	}

	response := response.BillingStudentByStudentIDResponse{
		ID:               student.ID,
		Nis:              student.Nis,
		StudentName:      student.FullName,
		SchoolClassName:  schoolClassName,
		SchoolGradeName:  schoolGradeName,
		LatestSchoolYear: latestSchoolYear,
		ListBilling:      []response.BillingStudentByStudentIDDetailResponse{},
		ListDonation:     []response.DonationBillingResponse{},
	}

	installmentBillings, donations, err := billingStudentService.billingStudentRepository.GetInstallmentDetails(int(student.ID), user.UserSchool.School.ID)
	if err != nil {
		return &response, err
	}
	response.ListBilling = installmentBillings
	response.ListDonation = donations

	return &response, nil
}

func (billingStudentService *BillingStudentService) GetDetailBillingStudentByID(billingStudentId int) (*response.BillingStudentDetailResponse, error) {
	billingStudentDetail, err := billingStudentService.billingStudentRepository.GetDetailBillingStudentByID(billingStudentId)
	if err != nil {
		return nil, err
	}

	rsp := response.BillingStudentDetailResponse{
		BillingStudentId:  int(billingStudentDetail.ID),
		DetailBillingName: billingStudentDetail.DetailBillingName,
		DueDate:           billingStudentDetail.DueDate.Format(constants.DateFormatYYYYMMDD),
		Amount:            int(billingStudentDetail.Amount),
	}
	return &rsp, nil
}

func (service *BillingStudentService) UpdateBillingStudentService(billingStudentId int, updateRequest request.UpdateBillingStudentRequest) (response.BillingStudentDetailResponse, error) {
	// Ambil data BillingStudent berdasarkan ID
	billingStudent, err := service.billingStudentRepository.GetDetailBillingStudentByID(billingStudentId)
	if err != nil {
		return response.BillingStudentDetailResponse{}, err
	}

	// Update fields sesuai request
	billingStudent.DetailBillingName = updateRequest.DetailBillingName
	billingStudent.Amount = updateRequest.Amount
	// Parse due date dari string ke time.Time
	dueDate, err := time.Parse(constants.DateFormatYYYYMMDD, updateRequest.DueDate)
	if err != nil {
		return response.BillingStudentDetailResponse{}, fmt.Errorf("format tanggal tidak valid: %v", err)
	}
	billingStudent.DueDate = &dueDate

	// Simpan perubahan ke database
	if err := service.billingStudentRepository.UpdateBillingStudent(&billingStudent); err != nil {
		return response.BillingStudentDetailResponse{}, fmt.Errorf("gagal memperbarui data: %v", err)
	}

	// Mengemas data dalam variabel responseData
	responseData := response.BillingStudentDetailResponse{
		BillingStudentId:  int(billingStudent.ID),
		DetailBillingName: billingStudent.DetailBillingName,
		DueDate:           dueDate.Format(constants.DateFormatYYYYMMDD), // Konversi dueDate ke string
		Amount:            int(billingStudent.Amount),
	}

	return responseData, nil
}

func (service *BillingStudentService) DeleteBillingStudentService(billingStudentId, deletedBy int) error {
	// Pastikan data BillingStudent dengan ID tersebut ada
	_, err := service.billingStudentRepository.GetDetailBillingStudentByID(billingStudentId)
	if err != nil {
		return fmt.Errorf("data billing student tidak ditemukan")
	}
	// Lakukan soft delete
	err = service.billingStudentRepository.DeleteBillingStudent(billingStudentId, deletedBy)
	if err != nil {
		return fmt.Errorf("gagal menghapus data billing student: %v", err)
	}

	return nil
}

func (service *BillingStudentService) CreateBillingStudent(request request.CreateBillingStudentRequest, userID int) ([]models.BillingStudent, error) {
	var createdBillingStudents []models.BillingStudent

	// Validate student exists - single query
	_, err := service.studentRepository.GetStudentByID(uint(request.StudentID), models.User{})
	if err != nil {
		return nil, fmt.Errorf("student not found")
	}

	// Collect all billing IDs and detail IDs for batch validation
	var billingIDs []uint
	var detailIDs []uint
	billingDetailMap := make(map[uint]uint) // map[detailID]billingID

	for _, detail := range request.DetailDataBilling {
		billingIDs = append(billingIDs, uint(detail.BillingID))
		detailIds := utilities.SplitBillingDetailIds(detail.BillingDetailIds)
		if len(detailIds) == 0 {
			return nil, fmt.Errorf("no billing detail IDs provided for billing ID %d", detail.BillingID)
		}
		for _, detailID := range detailIds {
			detailIDs = append(detailIDs, uint(detailID))
			billingDetailMap[uint(detailID)] = uint(detail.BillingID)
		}
	}

	// Batch validate billing details and get their data
	billingDetails, err := service.billingStudentRepository.GetBillingDetailsByIDs(detailIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching billing details: %v", err)
	}

	// Validate billing detail associations and create map for quick access
	detailsMap := make(map[uint]*models.BillingDetail)
	for _, detail := range billingDetails {
		if expectedBillingID, exists := billingDetailMap[detail.ID]; !exists || detail.BillingID != expectedBillingID {
			return nil, fmt.Errorf("billing detail ID %d does not belong to expected billing ID", detail.ID)
		}
		detailsMap[detail.ID] = detail
	}

	// Batch check for existing combinations
	existingCombinations, err := service.billingStudentRepository.CheckBulkBillingStudentExists(uint(request.StudentID), detailIDs)
	if err != nil {
		return nil, fmt.Errorf("error checking existing combinations: %v", err)
	}

	// Create a map of existing combinations for quick lookup
	existingMap := make(map[string]bool)
	for _, combo := range existingCombinations {
		key := fmt.Sprintf("%d-%d", combo.StudentID, combo.BillingDetailID)
		existingMap[key] = true
	}

	// Prepare batch insert
	var billingStudentsToCreate []models.BillingStudent
	uniqueCombinations := make(map[string]bool)

	for detailID, detail := range detailsMap {
		// Check if combination already exists
		existingKey := fmt.Sprintf("%d-%d", request.StudentID, detailID)
		if existingMap[existingKey] {
			continue
		}

		uniqueKey := fmt.Sprintf("%d-%d-%d", request.StudentID, detail.BillingID, detailID)
		if uniqueCombinations[uniqueKey] {
			continue
		}

		billingStudent := models.BillingStudent{
			Master: models.Master{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
			BillingID:         detail.BillingID,
			StudentID:         uint(request.StudentID),
			PaymentStatus:     "1",
			DueDate:           detail.DueDate,
			DetailBillingName: detail.DetailBillingName,
			Amount:            detail.Amount,
			BillingDetailID:   detailID,
		}

		billingStudentsToCreate = append(billingStudentsToCreate, billingStudent)
		uniqueCombinations[uniqueKey] = true
	}

	// Only proceed if there are new combinations to create
	if len(billingStudentsToCreate) == 0 {
		return []models.BillingStudent{}, nil
	}

	// Batch insert all billing students
	createdBillingStudents, err = service.billingStudentRepository.BulkCreateBillingStudents(billingStudentsToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing students: %v", err)
	}

	if len(createdBillingStudents) == 0 {
		return nil, fmt.Errorf("no billing students created")
	}

	return createdBillingStudents, nil
}
