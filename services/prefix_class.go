package services

import (
	"errors"

	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
)

type PrefixClassService struct {
	prefixClassRepository repositories.IPrefixClassRepository
	userRepository        repositories.UserRepository
}

func NewPrefixClassService() *PrefixClassService {
	return &PrefixClassService{
		prefixClassRepository: repositories.NewPrefixClassRepository(),
		userRepository:        repositories.NewUserRepository(),
	}
}

func NewPrefixClassServiceWithRepo(prefixRepo repositories.IPrefixClassRepository, userRepo repositories.UserRepository) *PrefixClassService {
	return &PrefixClassService{
		prefixClassRepository: prefixRepo,
		userRepository:        userRepo,
	}
}

func (service *PrefixClassService) GetAllPrefixClassService(search string, userID int) (response.PrefixCLassResponse, error) {
	var resp response.PrefixCLassResponse
	var listDetail []response.DetailPrefixResponse

	user, err := service.userRepository.GetUserByID(uint(userID))
	if err != nil {
		resp.Data = []response.DetailPrefixResponse{}
		return resp, nil
	}

	dataPrefix, err := service.prefixClassRepository.GetAllPrefixClassRepository(search, user)
	if err != nil {
		resp.Data = []response.DetailPrefixResponse{}
		return resp, nil
	}

	for _, prefix := range dataPrefix {
		listDetail = append(listDetail, response.DetailPrefixResponse{
			ID:         prefix.ID,
			PrefixName: prefix.PrefixName,
		})
	}

	if len(listDetail) > 0 {
		resp.Data = listDetail
	} else {
		resp.Data = []response.DetailPrefixResponse{}
	}

	return resp, nil
}

func (s *PrefixClassService) CreatePrefixClassService(prefixClass *request.PrefixClassCreate, userID int) (*models.PrefixClass, error) {
	exists, err := s.prefixClassRepository.CheckPrefixClassExists(prefixClass.PrefixName, prefixClass.SchoolID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("prefixName already exists for this school")
	}

	prefixClassModel := &models.PrefixClass{
		PrefixName: prefixClass.PrefixName,
		SchoolID:   prefixClass.SchoolID,
	}
	prefixClassModel.Master.CreatedBy = userID

	dataPrefix, err := s.prefixClassRepository.CreatePrefixClassRepository(prefixClassModel)
	if err != nil {
		return nil, err
	}

	return dataPrefix, nil
}
