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

type SchoolYearServiceInterface interface {
	GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string,  userId int) (response.SchoolYearListResponse, error)
	CreateSchoolYear(schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error)
	UpdateSchoolYear(id uint, schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error)
	DeleteSchoolYear(id uint, userID int) (*models.SchoolYear, error)
	GetSchoolYearByID(id uint) (*response.SchoolYearDetailResponse, error)
}

type SchoolYearService struct {
	schoolYearRepository repositories.SchoolYearRepository
	userRepository       repositories.UserRepository
}

func NewSchoolYearService(schoolYearRepository repositories.SchoolYearRepository, userRepository repositories.UserRepository) *SchoolYearService {
	return &SchoolYearService{schoolYearRepository: schoolYearRepository, userRepository: userRepository}
}

func (schoolYearService *SchoolYearService) GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, userId int) (response.SchoolYearListResponse, error) {
	var mapSchoolYear response.SchoolYearListResponse
	var listDetail []response.DetailSchoolYearResponse
	mapSchoolYear.Limit = limit
	mapSchoolYear.Page = page
	mapSchoolYear.TotalData = 0
	mapSchoolYear.TotalPage = 0

	user, err := schoolYearService.userRepository.GetUserByID(uint(userId))
    if err != nil {
		mapSchoolYear.Data = []response.DetailSchoolYearResponse{}
        return mapSchoolYear,nil
    }

	if sortBy != "" {
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	listSchoolYear, totalData, totalPage, err := schoolYearService.schoolYearRepository.GetAllSchoolYear(page, limit, search, sortBy, sortOrder, user.UserSchool.SchoolID)
	if err != nil {
		mapSchoolYear.Data = []response.DetailSchoolYearResponse{}
		return mapSchoolYear, nil
	}

	for _, schoolYear := range listSchoolYear {
		listDetail = append(listDetail, response.DetailSchoolYearResponse{
			ID:             schoolYear.ID,
			StartDate:      schoolYear.StartDate,
			EndDate:        schoolYear.EndDate,
			SchoolYearName: schoolYear.SchoolYearName,
			CreatedAt:      schoolYear.CreatedAt,
			CreatedBy:      schoolYear.CreateByUsername,
			UpdatedAt:      schoolYear.UpdatedAt,
		})
	}

	if len(listDetail) > 0 {
		mapSchoolYear.Data = listDetail
	} else {
		mapSchoolYear.Data = []response.DetailSchoolYearResponse{}
	}

	mapSchoolYear.TotalData = int(totalData)
	mapSchoolYear.TotalPage = totalPage

	return mapSchoolYear, nil
}

func (schoolYearService *SchoolYearService) GetSchoolYearByID(id uint) (*response.SchoolYearDetailResponse, error) {
	schoolYear, err := schoolYearService.schoolYearRepository.GetSchoolYearByID(id)
	if err != nil {
		return nil, err
	}
	response := response.SchoolYearDetailResponse{
		ID:             int(schoolYear.ID),
		SchoolYearName: schoolYear.SchoolYearName,
		StartDate:      schoolYear.StartDate,
		EndDate:        schoolYear.EndDate,
	}
	return &response, nil
}

func (schoolYearService *SchoolYearService) CreateSchoolYear(schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error) {
	var startDateFormated *time.Time
	var endDateFormated *time.Time

	schoolYearCode, err := schoolYearService.GenerateSchoolYearCode()
	if err != nil {
		return nil, err
	}

	startDateFormated, err = utilities.ChangeDate(schoolYearRequest.StartDate)
	if err != nil {
		return nil, err
	}

	if schoolYearRequest.EndDate != "" {
		// Validate that the end date is after the start date
		if err := utilities.ValidateDateRange(schoolYearRequest.StartDate, schoolYearRequest.EndDate); err != nil {
			return nil, err
		}

		endDateFormated, err = utilities.ChangeDate(schoolYearRequest.EndDate)
		if err != nil {
			return nil, err
		}
	}

	schoolYear := models.SchoolYear{
		SchoolYearCode: schoolYearCode,
		SchoolYearName: schoolYearRequest.SchoolYearName,
		SchoolId:       schoolYearRequest.SchoolId,
		StartDate:      startDateFormated,
		EndDate:        endDateFormated,
	}
	schoolYear.Master.CreatedBy = userID
	schoolYear.Master.UpdatedBy = userID

	dataSchoolYear, err := schoolYearService.schoolYearRepository.CreateSchoolYear(&schoolYear)
	if err != nil {
		return nil, err
	}

	return dataSchoolYear, nil
}

func (schoolYearService *SchoolYearService) UpdateSchoolYear(id uint, schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error) {
	getSchoolYear, err := schoolYearService.schoolYearRepository.GetSchoolYearByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	// Validate and format StartDate if provided
	if schoolYearRequest.StartDate != "" {
		startDateFormatted, err := utilities.ChangeDate(schoolYearRequest.StartDate)
		if err != nil {
			return nil, err
		}
		getSchoolYear.StartDate = startDateFormatted
	}

	// Validate and format EndDate if provided
	if schoolYearRequest.EndDate != "" {
		if schoolYearRequest.StartDate != "" {
			// Validate that EndDate is after StartDate
			if err := utilities.ValidateDateRange(schoolYearRequest.StartDate, schoolYearRequest.EndDate); err != nil {
				return nil, err
			}
		}

		endDateFormatted, err := utilities.ChangeDate(schoolYearRequest.EndDate)
		if err != nil {
			return nil, err
		}
		getSchoolYear.EndDate = endDateFormatted
	}

	// Conditionally update other fields
	if schoolYearRequest.SchoolYearName != "" {
		getSchoolYear.SchoolYearName = schoolYearRequest.SchoolYearName
	}
	if schoolYearRequest.SchoolId != 0 {
		getSchoolYear.SchoolId = schoolYearRequest.SchoolId
	}
	getSchoolYear.Master.UpdatedBy = userID

	dataSchoolYear, err := schoolYearService.schoolYearRepository.UpdateSchoolYear(&getSchoolYear)
	if err != nil {
		return nil, err
	}

	return dataSchoolYear, nil
}

func (schoolYearService *SchoolYearService) DeleteSchoolYear(id uint, userID int) (*models.SchoolYear, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getSchoolYear, err := schoolYearService.schoolYearRepository.GetSchoolYearByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getSchoolYear.Master.DeletedAt = currentTimePointer
	getSchoolYear.Master.DeletedBy = &userID
	dataSchoolYear, err := schoolYearService.schoolYearRepository.UpdateSchoolYear(&getSchoolYear)
	if err != nil {
		return nil, err
	}

	return dataSchoolYear, nil
}

func (schoolYearService *SchoolYearService) GenerateSchoolYearCode() (string, error) {
	lastNumber, err := schoolYearService.schoolYearRepository.GetLastSequenceNumberSchoolYears()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1
	newCode := fmt.Sprintf("SY%03d", newSequence)

	return newCode, nil
}
