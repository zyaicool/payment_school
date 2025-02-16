package services

import (
	"fmt"
	"time"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"
)

type SchoolClassService struct {
	schoolClassRepository repositories.SchoolClassRepositoryInterface
	userRepository        repositories.UserRepository
}

func NewSchoolClassService(schoolClassRepository repositories.SchoolClassRepositoryInterface, userRepository repositories.UserRepository) SchoolClassService {
	return SchoolClassService{
		schoolClassRepository: schoolClassRepository,
		userRepository:        userRepository,
	}
}

func (schoolClassService *SchoolClassService) GetAllSchoolClass(page int, limit int, search string, sortBy string, sortOrder string, showDeletedData bool, user models.User) (response.SchoolClassListResponse, error) {
	var responseList response.SchoolClassListResponse
	var listSchoolClassNew []response.SchoolClassResponse
	responseList.Page = page
	responseList.Limit = limit
	responseList.TotalPage = 0
	responseList.TotalData = 0

	if sortBy != "" {
		sortBy = utilities.ChangeStringSortBySchoolClass(sortBy)
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	// Ambil data dari repository
	listSchoolClass, totalPage, totalData, err := repositories.GetAllSchoolClass(page, limit, search, sortBy, sortOrder, showDeletedData, user)
	if err != nil {
		responseList.Data = []response.SchoolClassResponse{}
		return responseList, err
	}

	// Total data dan perhitungan total halaman
	responseList.TotalData = totalData
	responseList.TotalPage = totalPage

	for _, schoolClass := range listSchoolClass {

		listSchoolClassNew = append(listSchoolClassNew, response.SchoolClassResponse{
			ID:              int(schoolClass.ID),
			Unit:            schoolClass.Unit,
			PrefixClass:     schoolClass.PrefixClass,
			SchoolMajor:     schoolClass.SchoolMajor,
			SchoolClassName: schoolClass.SchoolClassName,
			CreatedBy:       schoolClass.CreatedByUsername,
			Status:          schoolClass.Status,
			IsEdit:          schoolClass.IsEdit,
			Placeholder:     schoolClass.Placeholder,
		})
	}

	// Mapping data ke format response yang diinginkan
	responseList.Data = listSchoolClassNew

	return responseList, nil
}

func (schoolClassService *SchoolClassService) GetSchoolClassByID(id uint) (*response.SchoolClassDetailResponse, error) {
	schoolClass, err := schoolClassService.schoolClassRepository.GetSchoolClassByID(id)
	if err != nil {
		return nil, err
	}
	response := response.SchoolClassDetailResponse{
		ID:              int(schoolClass.ID),
		SchoolGradeID:   int(schoolClass.SchoolGradeID),
		PrefixClassID:   schoolClass.PrefixClassID,
		SchoolMajorID:   schoolClass.SchoolMajorID,
		Suffix:          schoolClass.Suffix,
		SchoolClassName: schoolClass.SchoolClassName,
	}
	return &response, nil
}

func (schoolClassService *SchoolClassService) CreateSchoolClass(schoolClassRequest *request.SchoolClassCreateUpdateRequest, userID int) (*models.SchoolClass, error) {
	schoolClassCode, err := schoolClassService.GenerateSchoolClassCode()
	if err != nil {
		return nil, err
	}

	schoolID := 0
	if schoolClassRequest.SchoolID == 0 {
		getSchoolId, err := repositories.GetUserSchoolByUserId(uint(userID))
		if err != nil {
			return nil, err
		}
		schoolID = int(getSchoolId.SchoolID)
	} else {
		schoolID = schoolClassRequest.SchoolID
	}

	schoolClass := models.SchoolClass{
		SchoolID:        uint(schoolID),
		SchoolGradeID:   uint(schoolClassRequest.SchoolGradeID),
		SchoolClassCode: schoolClassCode,
		SchoolClassName: schoolClassRequest.SchoolClassName,
		Suffix:          schoolClassRequest.Suffix,
		SchoolMajorID:   schoolClassRequest.SchoolMajorID,
		PrefixClassID:   schoolClassRequest.PrefixClassID,
	}
	schoolClass.Master.CreatedBy = userID
	schoolClass.Master.UpdatedBy = userID

	dataSchoolClass, err := schoolClassService.schoolClassRepository.CreateSchoolClass(&schoolClass)
	if err != nil {
		return nil, err
	}

	return dataSchoolClass, nil
}

func (schoolClassService *SchoolClassService) UpdateSchoolClass(id uint, schoolClassRequest *request.SchoolClassCreateUpdateRequest, userID int) (*models.SchoolClass, error) {
	getSchoolClass, err := schoolClassService.schoolClassRepository.GetSchoolClassByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getSchoolClass.SchoolClassName = schoolClassRequest.SchoolClassName
	getSchoolClass.PrefixClassID = schoolClassRequest.PrefixClassID
	getSchoolClass.SchoolID = uint(schoolClassRequest.SchoolID)
	getSchoolClass.SchoolGradeID = uint(schoolClassRequest.SchoolGradeID)
	getSchoolClass.SchoolMajorID = schoolClassRequest.SchoolMajorID
	getSchoolClass.Suffix = schoolClassRequest.Suffix
	getSchoolClass.Master.UpdatedBy = userID
	getSchoolClass.Master.CreatedAt = getSchoolClass.Master.CreatedAt
	getSchoolClass.Master.CreatedBy = getSchoolClass.Master.CreatedBy
	getSchoolClass.SchoolClassCode = getSchoolClass.SchoolClassCode

	dataSchoolClass, err := schoolClassService.schoolClassRepository.UpdateSchoolClass(&getSchoolClass)
	if err != nil {
		return nil, err
	}

	return dataSchoolClass, nil
}

func (schoolClassService *SchoolClassService) DeleteSchoolClass(id uint, userID int) (*models.SchoolClass, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getSchoolClass, err := schoolClassService.schoolClassRepository.GetSchoolClassByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	// getSchoolClass.Master.DeletedAt = currentTimePointer
	// cek del
	if getSchoolClass.Master.DeletedAt == nil {
		getSchoolClass.Master.DeletedAt = currentTimePointer
		getSchoolClass.Master.DeletedBy = &userID
	} else {
		getSchoolClass.Master.DeletedAt = nil
		getSchoolClass.Master.DeletedBy = nil
	}

	dataSchoolClass, err := schoolClassService.schoolClassRepository.UpdateSchoolClass(&getSchoolClass)
	if err != nil {
		return nil, err
	}

	return dataSchoolClass, nil
}

func (schoolClassService *SchoolClassService) GenerateSchoolClassCode() (string, error) {
	lastNumber, err := schoolClassService.schoolClassRepository.GetLastSequenceNumberSchoolClasss()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1
	newCode := fmt.Sprintf("SC%03d", newSequence)

	return newCode, nil
}
