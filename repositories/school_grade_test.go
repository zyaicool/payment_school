package repositories

import (
	"errors"
	"testing"

	"schoolPayment/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock repository for SchoolGradeRepository
type MockSchoolGradeRepository struct {
	mock.Mock
}

func (m *MockSchoolGradeRepository) CreateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error) {
	args := m.Called(schoolGrade)
	return args.Get(0).(*models.SchoolGrade), args.Error(1)
}

func (m *MockSchoolGradeRepository) UpdateSchoolGrade(schoolGrade *models.SchoolGrade) (*models.SchoolGrade, error) {
	args := m.Called(schoolGrade)
	return args.Get(0).(*models.SchoolGrade), args.Error(1)
}

func (m *MockSchoolGradeRepository) GetLastSequenceNumberSchoolGrades() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// Test suite structure
type SchoolGradeRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockSchoolGradeRepository
}

func (suite *SchoolGradeRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(MockSchoolGradeRepository)
}

// Test for creating school grade
func (suite *SchoolGradeRepositoryTestSuite) TestCreateSchoolGrade_Success() {
	dummySchoolGrade := &models.SchoolGrade{
		SchoolGradeName: "Grade 10",
	}

	// Setting up expected behavior
	suite.mockRepo.On("CreateSchoolGrade", dummySchoolGrade).Return(dummySchoolGrade, nil)

	// Calling the service
	result, err := suite.mockRepo.CreateSchoolGrade(dummySchoolGrade)

	// Assert the expectations
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Grade 10", result.SchoolGradeName)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolGradeRepositoryTestSuite) TestCreateSchoolGrade_Error() {
	dummySchoolGrade := &models.SchoolGrade{
		SchoolGradeName: "Grade 10",
	}

	// Setting up expected error
	suite.mockRepo.On("CreateSchoolGrade", dummySchoolGrade).Return((*models.SchoolGrade)(nil), errors.New("create error"))

	// Calling the service
	result, err := suite.mockRepo.CreateSchoolGrade(dummySchoolGrade)

	// Assert the expectations
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test for updating school grade
func (suite *SchoolGradeRepositoryTestSuite) TestUpdateSchoolGrade_Success() {
	dummySchoolGrade := &models.SchoolGrade{
		SchoolGradeName: "Grade 11",
	}

	// Setting up expected behavior
	suite.mockRepo.On("UpdateSchoolGrade", dummySchoolGrade).Return(dummySchoolGrade, nil)

	// Calling the service
	result, err := suite.mockRepo.UpdateSchoolGrade(dummySchoolGrade)

	// Assert the expectations
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Grade 11", result.SchoolGradeName)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolGradeRepositoryTestSuite) TestGetLastSequenceNumberSchoolGrades_Success() {
	// Setting up expected behavior
	suite.mockRepo.On("GetLastSequenceNumberSchoolGrades").Return(10, nil)

	// Calling the service
	sequence, err := suite.mockRepo.GetLastSequenceNumberSchoolGrades()

	// Assert the expectations
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, sequence)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolGradeRepositoryTestSuite) TestGetLastSequenceNumberSchoolGrades_Error() {
	// Setting up expected error
	suite.mockRepo.On("GetLastSequenceNumberSchoolGrades").Return(0, errors.New("database error"))

	// Calling the service
	sequence, err := suite.mockRepo.GetLastSequenceNumberSchoolGrades()

	// Assert the expectations
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, sequence)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Running the test suite
func TestSchoolGradeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SchoolGradeRepositoryTestSuite))
}
