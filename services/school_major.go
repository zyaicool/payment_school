package services

import (
	"errors"

	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
)

type SchoolMajorService interface {
	GetAllSchoolMajorService(search string, userID int) (response.SchoolMajorResponse, error)
	CreateSchoolMajorService(schoolMajorRequest *request.SchoolMajorCreate, userID int) (*models.SchoolMajor, error)
}

type schoolMajorService struct {
	schoolMajorRepository repositories.SchoolMajorRepository
	userRepository        repositories.UserRepository
}

func NewSchoolMajorService(schoolMajorRepository repositories.SchoolMajorRepository, userRepository repositories.UserRepository) SchoolMajorService {
	return &schoolMajorService{
		schoolMajorRepository: schoolMajorRepository,
		userRepository:        userRepository,
	}
}


func (service *schoolMajorService) GetAllSchoolMajorService(search string, userID int) (response.SchoolMajorResponse, error) {
	var resp response.SchoolMajorResponse
	var listDetailMajor []response.DetailSchoolMajorResponse

	user, err := service.userRepository.GetUserByID(uint(userID))
	if err != nil {
		resp.Data = []response.DetailSchoolMajorResponse{}
	}

	dataMajor, err := service.schoolMajorRepository.GetAllSchoolMajorRepository(search, user)
	if err != nil {
		resp.Data = []response.DetailSchoolMajorResponse{}
		return resp, nil
	}

	for _, major := range dataMajor {
		listDetailMajor = append(listDetailMajor, response.DetailSchoolMajorResponse{
			ID:              major.ID,
			SchoolMajorName: major.SchoolMajorName,
		})
	}

	if len(listDetailMajor) > 0 {
		resp.Data = listDetailMajor
	} else {
		resp.Data = []response.DetailSchoolMajorResponse{}
	}

	return resp, nil
}

func (service *schoolMajorService) CreateSchoolMajorService(schoolMajorRequest *request.SchoolMajorCreate, userID int) (*models.SchoolMajor, error) {
	exists, err := service.schoolMajorRepository.CheckSchoolMajorExists(schoolMajorRequest.SchoolMajorName, schoolMajorRequest.SchoolID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("majorName already exists for this school")
	}

	schoolMajor := models.SchoolMajor{
		SchoolMajorName: schoolMajorRequest.SchoolMajorName,
		SchoolID:        schoolMajorRequest.SchoolID,
	}

	schoolMajor.Master.CreatedBy = userID

	dataMajor, err := service.schoolMajorRepository.CreateSchoolMajorRepository(&schoolMajor)
	if err != nil {
		return nil, err
	}

	return dataMajor, nil
}