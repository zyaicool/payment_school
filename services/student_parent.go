package services

import (
	"fmt"

	"schoolPayment/constants"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"
)

// Definisikan interface
type StudentParentService interface {
	GetAllParent(page int, limit int, search string, userID int) (response.StudentParentResponse, error)
	GetParentByID(id uint, userID int) (*models.StudentParent, error)
	CreateParent(parent *models.StudentParent, userID int) (*models.StudentParent, error)
	UpdateParent(id uint, parent *models.StudentParent, userID int) (*models.StudentParent, error)
	GetParentByUserLogin(id uint) (*models.StudentParent, error)
	CreateBatchParent(parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error)
	UpdateBatchParent(id uint, parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error)
}

// Definisikan struct dengan nama yang berbeda (huruf kecil untuk private struct)
type studentParentService struct {
	studentParentRepository repositories.StudentParentRepository
	userRepository          repositories.UserRepository
}

// Fungsi constructor untuk membuat instance dari studentParentService
func NewStudentParentService(studentParentRepository repositories.StudentParentRepository, userRepository repositories.UserRepository) StudentParentService {
	return &studentParentService{
		studentParentRepository: studentParentRepository,
		userRepository:          userRepository,
	}
}

// Implementasikan method-method dari interface StudentParentService
func (s *studentParentService) GetAllParent(page int, limit int, search string, userID int) (response.StudentParentResponse, error) {
	var mapParent response.StudentParentResponse
	mapParent.Limit = limit
	mapParent.Page = page

	user, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		mapParent.Data = []models.StudentParent{}
		return mapParent, nil
	}

	listParent, err := s.studentParentRepository.GetAllParent(page, limit, search, user)
	if err != nil {
		mapParent.Data = []models.StudentParent{}
		return mapParent, nil
	}

	if len(listParent) > 0 {
		mapParent.Data = listParent
	} else {
		mapParent.Data = []models.StudentParent{}
	}

	return mapParent, nil
}

func (s *studentParentService) GetParentByID(id uint, userID int) (*models.StudentParent, error) {
	user, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	getParent, err := s.studentParentRepository.GetParentByID(id, user)
	if err != nil {
		return nil, err
	}

	return &getParent, nil
}

func (s *studentParentService) CreateParent(parent *models.StudentParent, userID int) (*models.StudentParent, error) {
	// check format email
	if parent.ParentMail != "" {
		err := utilities.ValidateEmail(parent.ParentMail)
		if err != nil {
			return nil, err
		}
	}

	getUser, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	if getUser.RoleID != 2 {
		return nil, fmt.Errorf(constants.ErrorMessageLoginAsParent)
	}

	// set user id to created by and updated by
	parent.UserID = uint(userID)
	parent.Master.CreatedBy = userID
	parent.Master.UpdatedBy = userID

	dataParent, err := s.studentParentRepository.CreateParent(parent)
	if err != nil {
		return nil, err
	}

	return dataParent, nil
}

func (s *studentParentService) UpdateParent(id uint, parent *models.StudentParent, userID int) (*models.StudentParent, error) {
	user, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	getParent, err := s.studentParentRepository.GetParentByID(id, user)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	// check format email
	if parent.ParentMail != "" {
		err = utilities.ValidateEmail(parent.ParentMail)
		if err != nil {
			return nil, err
		}
	}

	getUser, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	if getUser.RoleID != 2 {
		return nil, fmt.Errorf(constants.ErrorMessageLoginAsParent)
	}

	getParent.ParentName = parent.ParentName
	getParent.ParentAddress = parent.ParentAddress
	getParent.ParentHandphone = parent.ParentHandphone
	getParent.ParentMail = parent.ParentMail
	getParent.ParentCitizenship = parent.ParentCitizenship
	getParent.ParentSalary = parent.ParentSalary
	getParent.ParentStatus = parent.ParentStatus
	getParent.Master.UpdatedBy = userID

	dataParent, err := s.studentParentRepository.UpdateParent(&getParent)
	if err != nil {
		return nil, err
	}

	return dataParent, nil
}

func (s *studentParentService) GetParentByUserLogin(id uint) (*models.StudentParent, error) {
	getUser, err := s.userRepository.GetUserByID(uint(id))
	if err != nil {
		return nil, err
	}

	if getUser.RoleID != 2 {
		return nil, fmt.Errorf(constants.ErrorMessageLoginAsParent)
	}

	return repositories.GetParentByUserLogin(id)
}

func (s *studentParentService) CreateBatchParent(parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error) {
	var listParent *[]*models.StudentParent = &[]*models.StudentParent{}

	for _, parent := range parents {
		// Check format email
		if parent.ParentMail != "" {
			err := utilities.ValidateEmail(parent.ParentMail)
			if err != nil {
				return nil, err
			}
		}

		getUser, err := s.userRepository.GetUserByID(uint(userID))
		if err != nil {
			return nil, err
		}

		if getUser.RoleID != 2 {
			return nil, fmt.Errorf(constants.ErrorMessageLoginAsParent)
		}

		// Set user ID to createdBy and updatedBy
		parent.UserID = uint(userID)
		parent.Master.CreatedBy = userID
		parent.Master.UpdatedBy = userID

		// Create parent data
		dataParent, err := s.studentParentRepository.CreateParent(parent)
		if err != nil {
			return nil, err
		}

		// Append dataParent to listParent
		*listParent = append(*listParent, dataParent)
	}

	return listParent, nil
}

func (s *studentParentService) UpdateBatchParent(id uint, parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error) {
	var listParent *[]*models.StudentParent = &[]*models.StudentParent{}
	user, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	for _, parent := range parents {
		getParent, err := s.studentParentRepository.GetParentByID(id, user)
		if err != nil {
			return nil, fmt.Errorf("Data not found.")
		}

		// check format email
		if parent.ParentMail != "" {
			err = utilities.ValidateEmail(parent.ParentMail)
			if err != nil {
				return nil, err
			}
		}

		getUser, err := s.userRepository.GetUserByID(uint(userID))
		if err != nil {
			return nil, err
		}

		if getUser.RoleID != 2 {
			return nil, fmt.Errorf(constants.ErrorMessageLoginAsParent)
		}

		getParent.ParentName = parent.ParentName
		getParent.ParentAddress = parent.ParentAddress
		getParent.ParentHandphone = parent.ParentHandphone
		getParent.ParentMail = parent.ParentMail
		getParent.ParentCitizenship = parent.ParentCitizenship
		getParent.ParentSalary = parent.ParentSalary
		getParent.ParentStatus = parent.ParentStatus
		getParent.Master.UpdatedBy = userID

		dataParent, err := s.studentParentRepository.UpdateParent(&getParent)
		if err != nil {
			return nil, err
		}

		*listParent = append(*listParent, dataParent)
	}
	return listParent, nil
}