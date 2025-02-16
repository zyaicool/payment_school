package repositories

import (
	"errors"
	"schoolPayment/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockSchoolClassRepository struct {
	mock.Mock
}

func (m *MockSchoolClassRepository) CreateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	args := m.Called(schoolClass)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassRepository) GetSchoolClassByID(id uint) (models.SchoolClass, error) {
	args := m.Called(id)
	return args.Get(0).(models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassRepository) UpdateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	args := m.Called(schoolClass)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

type SchoolClassRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockSchoolClassRepository
}

func (suite *SchoolClassRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(MockSchoolClassRepository)
}
func (suite *SchoolClassRepositoryTestSuite) TestCreateSchoolClass_Success() {
	dummySchoolClass := &models.SchoolClass{
		SchoolClassName: "Class 10",
		SchoolGradeID:   1,
	}

	// Setting up expected behavior
	suite.mockRepo.On("CreateSchoolClass", dummySchoolClass).Return(dummySchoolClass, nil)

	// Calling the service
	result, err := suite.mockRepo.CreateSchoolClass(dummySchoolClass)

	// Assert the expectations
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Class 10", result.SchoolClassName)
	assert.Equal(suite.T(), 1, result.SchoolClassName)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolClassRepositoryTestSuite) TestCreateSchoolClass_error() {
	dummySchoolClass := &models.SchoolClass{
		SchoolClassName: "Class 10",
		SchoolGradeID:   1,
	}

	// Setting up expected behavior
	suite.mockRepo.On("CreateSchoolClass", dummySchoolClass).Return((*models.SchoolGrade)(nil), errors.New("create error"))

	// Calling the service
	result, err := suite.mockRepo.CreateSchoolClass(dummySchoolClass)

	// Assert the expectations
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolClassRepositoryTestSuite) TestGetSchoolClassByIDRepository() {
	id := 12

	dummySchoolClassID := models.SchoolClass{
		Master:          models.Master{ID: 12},
		SchoolGradeID:   3,
		PrefixClassID:   1,
		SchoolMajorID:   1,
		Suffix:          "1",
		SchoolClassName: "XI IPA 4",
	}

	suite.mockRepo.On("GetSchoolClassByID", dummySchoolClassID).Return(dummySchoolClassID, nil)

	result, err := suite.mockRepo.GetSchoolClassByID(uint(id))

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "XI IPA 4", result.SchoolClassName)
	// assert.Equal(suite.T(), 1, res)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolClassRepositoryTestSuite) TestUpdateSchoolClass_error() {
	dummy := models.SchoolClass{
		SchoolGradeID:   1,
		PrefixClassID:   1,
		SchoolMajorID:   1,
		Suffix:          "testing 0001",
		SchoolClassName: "Test001",
	}

	suite.mockRepo.On("UpdateSchoolClass", dummy).Return((*models.SchoolClass)(nil), errors.New("update error"))

	result, err := suite.mockRepo.UpdateSchoolClass(&dummy)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SchoolClassRepositoryTestSuite) TestUpdateSchoolClassRepository() {
	dummy := models.SchoolClass{
		SchoolGradeID:   1,
		PrefixClassID:   1,
		SchoolMajorID:   1,
		Suffix:          "tes",
		SchoolClassName: "tes1",
	}

	suite.mockRepo.On("UpdateSchoolClass", dummy).Return(dummy, nil)

	result, err := suite.mockRepo.UpdateSchoolClass(&dummy)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "tes1", result.SchoolClassName)
	suite.mockRepo.AssertExpectations(suite.T())
}
