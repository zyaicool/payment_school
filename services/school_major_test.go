package services_test

import (
	"schoolPayment/dtos/request"
	"schoolPayment/models"
	"schoolPayment/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolMajorRepository struct {
	mock.Mock
}

func (m *MockSchoolMajorRepository) GetAllSchoolMajorRepository(search string, user models.User) ([]models.SchoolMajor, error) {
	args := m.Called(search, user)
	return args.Get(0).([]models.SchoolMajor), args.Error(1)
}

func (m *MockSchoolMajorRepository) CreateSchoolMajorRepository(major *models.SchoolMajor) (*models.SchoolMajor, error) {
	args := m.Called(major)
	return args.Get(0).(*models.SchoolMajor), args.Error(1)
}

func (m *MockSchoolMajorRepository) CheckSchoolMajorExists(majorName string, schoolID uint) (bool, error) {
	args := m.Called(majorName, schoolID)
	return args.Bool(0), args.Error(1)
}

func TestCreateSchoolMajorService(t *testing.T) {
	mockRepo := new(MockSchoolMajorRepository)
	service := services.NewSchoolMajorService(mockRepo, nil)

	// Test case where the major already exists
	mockRepo.On("CheckSchoolMajorExists", "Computer Science", uint(1)).Return(true, nil)

	majorRequest := &request.SchoolMajorCreate{
		SchoolMajorName: "Computer Science",
		SchoolID:        1,
	}

	// Call CreateSchoolMajorService
	_, err := service.CreateSchoolMajorService(majorRequest, 1)

	// Assert that error was returned
	assert.Error(t, err)
	assert.Equal(t, "majorName already exists for this school", err.Error())

	// Test case where the major does not exist
	mockRepo.On("CheckSchoolMajorExists", "Mathematics", uint(1)).Return(false, nil)
	mockRepo.On("CreateSchoolMajorRepository", mock.Anything).Return(&models.SchoolMajor{
		SchoolMajorName: "Mathematics",
		SchoolID:        1,
	}, nil)

	majorRequest.SchoolMajorName = "Mathematics"

	createdMajor, err := service.CreateSchoolMajorService(majorRequest, 1)

	// Assert successful creation
	assert.NoError(t, err)
	assert.Equal(t, "Mathematics", createdMajor.SchoolMajorName)
}
