package repositories

import (
	"errors"
	"schoolPayment/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock struct
type MockSchoolMajorRepository struct {
	mock.Mock
}

func (m *MockSchoolMajorRepository) CreateSchoolMajorRepository(major *models.SchoolMajor) (*models.SchoolMajor, error) {
	args := m.Called(major)
	return args.Get(0).(*models.SchoolMajor), args.Error(1)
}

func (m *MockSchoolMajorRepository) GetAllSchoolMajorRepository(search string, user models.User) ([]models.SchoolMajor, error) {
	args := m.Called(search, user)
	return args.Get(0).([]models.SchoolMajor), args.Error(1)
}

func (m *MockSchoolMajorRepository) CheckSchoolMajorExists(majorName string, schoolID uint) (bool, error) {
	args := m.Called(majorName, schoolID)
	return args.Bool(0), args.Error(1)
}

// Test suite struct
type SchoolMajorRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockSchoolMajorRepository
}

func (suite *SchoolMajorRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(MockSchoolMajorRepository)
}

// Test Success for CreateSchoolMajorRepository
func (suite *SchoolMajorRepositoryTestSuite) TestCreateSchoolMajor_Success() {
	dummyMajor := &models.SchoolMajor{
		SchoolMajorName: "Computer Science",
	}

	// Setting up expected behavior
	suite.mockRepo.On("CreateSchoolMajorRepository", dummyMajor).Return(dummyMajor, nil)

	// Calling the repository function
	result, err := suite.mockRepo.CreateSchoolMajorRepository(dummyMajor)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Computer Science", result.SchoolMajorName)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for CreateSchoolMajorRepository
func (suite *SchoolMajorRepositoryTestSuite) TestCreateSchoolMajor_Error() {
	dummyMajor := &models.SchoolMajor{
		SchoolMajorName: "Computer Science",
	}

	suite.mockRepo.On("CreateSchoolMajorRepository", dummyMajor).Return((*models.SchoolMajor)(nil), errors.New("create error"))

	// Calling the repository function
	result, err := suite.mockRepo.CreateSchoolMajorRepository(dummyMajor)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Success for GetAllSchoolMajorRepository
func (suite *SchoolMajorRepositoryTestSuite) TestGetAllSchoolMajor_Success() {
	dummyUser := models.User{
		UserSchool: &models.UserSchool{
			SchoolID: 1,
		},
	}

	dummyMajors := []models.SchoolMajor{
		{SchoolMajorName: "Computer Science"},
		{SchoolMajorName: "Mathematics"},
	}

	// Setting up expected behavior
	suite.mockRepo.On("GetAllSchoolMajorRepository", "", dummyUser).Return(dummyMajors, nil)

	// Calling the repository function
	result, err := suite.mockRepo.GetAllSchoolMajorRepository("", dummyUser)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Computer Science", result[0].SchoolMajorName)
	assert.Equal(suite.T(), "Mathematics", result[1].SchoolMajorName)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for GetAllSchoolMajorRepositoryfunc (suite *SchoolMajorRepositoryTestSuite) TestGetAllSchoolMajor_Error() {
	func (suite *SchoolMajorRepositoryTestSuite) TestGetAllSchoolMajor_Error() {
		dummyUser := models.User{
			UserSchool: &models.UserSchool{
				SchoolID: 1,
			},
		}
	
		// Mock behavior with slice empty and error
		suite.mockRepo.On("GetAllSchoolMajorRepository", "", dummyUser).Return([]models.SchoolMajor{}, errors.New("fetch error"))
	
		// Calling the repository function
		result, err := suite.mockRepo.GetAllSchoolMajorRepository("", dummyUser)
	
		// Assertions
		assert.Error(suite.T(), err)
		assert.Empty(suite.T(), result) // Memastikan bahwa result adalah slice kosong
		suite.mockRepo.AssertExpectations(suite.T())
	}
	


// Test Success for CheckSchoolMajorExists
func (suite *SchoolMajorRepositoryTestSuite) TestCheckSchoolMajorExists_Success() {
	majorName := "Computer Science"
	schoolID := uint(1)

	// Setting up expected behavior
	suite.mockRepo.On("CheckSchoolMajorExists", majorName, schoolID).Return(true, nil)

	// Calling the repository function
	result, err := suite.mockRepo.CheckSchoolMajorExists(majorName, schoolID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for CheckSchoolMajorExists
func (suite *SchoolMajorRepositoryTestSuite) TestCheckSchoolMajorExists_Error() {
	majorName := "Computer Science"
	schoolID := uint(1)

	suite.mockRepo.On("CheckSchoolMajorExists", majorName, schoolID).Return(false, errors.New("database error"))

	// Calling the repository function
	result, err := suite.mockRepo.CheckSchoolMajorExists(majorName, schoolID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.False(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestSchoolMajorRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SchoolMajorRepositoryTestSuite))
}
