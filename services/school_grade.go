package services

import (
	"fmt"
	"time"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
)

type SchoolGradeService struct {
	schoolGradeRepository repositories.SchoolGradeRepositoryInterface
}

func NewSchoolGradeService(schoolGradeRepository repositories.SchoolGradeRepositoryInterface) SchoolGradeService {
	return SchoolGradeService{schoolGradeRepository: schoolGradeRepository}
}

func GetAllSchoolGrade(page int, limit int, search string) (response.SchoolGradeListResponse, error) {
	var mapSchoolGrade response.SchoolGradeListResponse
	mapSchoolGrade.Limit = limit
	mapSchoolGrade.Page = page

	listSchoolGrade, err := repositories.GetAllSchoolGrade(page, limit, search)
	if err != nil {
		mapSchoolGrade.Data = []models.SchoolGrade{}
		return mapSchoolGrade, nil
	}

	if len(listSchoolGrade) > 0 {
		mapSchoolGrade.Data = listSchoolGrade
	} else {
		mapSchoolGrade.Data = []models.SchoolGrade{}
	}

	return mapSchoolGrade, nil
}

func (schoolGradeService *SchoolGradeService) GetSchoolGradeByID(id uint) (models.SchoolGrade, error) {
	return schoolGradeService.schoolGradeRepository.GetSchoolGradeByID(id)
}

func (schoolGradeService *SchoolGradeService) CreateSchoolGrade(schoolGradeRequest *request.SchoolGradeCreateUpdateRequest, userID int) (*models.SchoolGrade, error) {
	schoolGradeCode, err := schoolGradeService.GenerateSchoolGradeCode()
	if err != nil {
		return nil, err
	}
	schoolGrade := models.SchoolGrade{
		SchoolGradeCode: schoolGradeCode,
		SchoolGradeName: schoolGradeRequest.SchoolGradeName,
	}
	schoolGrade.Master.CreatedBy = userID
	schoolGrade.Master.UpdatedBy = userID

	dataSchoolGrade, err := schoolGradeService.schoolGradeRepository.CreateSchoolGrade(&schoolGrade)
	if err != nil {
		return nil, err
	}

	return dataSchoolGrade, nil
}

func (schoolGradeService *SchoolGradeService) UpdateSchoolGrade(id uint, schoolGradeRequest *request.SchoolGradeCreateUpdateRequest, userID int) (*models.SchoolGrade, error) {
	getSchoolGrade, err := schoolGradeService.schoolGradeRepository.GetSchoolGradeByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getSchoolGrade.SchoolGradeName = schoolGradeRequest.SchoolGradeName
	getSchoolGrade.Master.UpdatedBy = userID

	dataSchoolGrade, err := schoolGradeService.schoolGradeRepository.UpdateSchoolGrade(&getSchoolGrade)
	if err != nil {
		return nil, err
	}

	return dataSchoolGrade, nil
}

func (schoolGradeService *SchoolGradeService) DeleteSchoolGrade(id uint, userID int) (*models.SchoolGrade, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getSchoolGrade, err := schoolGradeService.schoolGradeRepository.GetSchoolGradeByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getSchoolGrade.Master.DeletedAt = currentTimePointer
	getSchoolGrade.Master.DeletedBy = &userID
	dataSchoolGrade, err := schoolGradeService.schoolGradeRepository.UpdateSchoolGrade(&getSchoolGrade)
	if err != nil {
		return nil, err
	}

	return dataSchoolGrade, nil
}

func (schoolGradeService *SchoolGradeService) GenerateSchoolGradeCode() (string, error) {
	lastNumber, err := schoolGradeService.schoolGradeRepository.GetLastSequenceNumberSchoolGrades()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1
	newCode := fmt.Sprintf("SG%03d", newSequence)

	return newCode, nil
}
