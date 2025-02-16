package services

import (
	"schoolPayment/models"

	"github.com/stretchr/testify/mock"
)

type MockSchoolGradeRepository struct {
	mock.Mock
}

func (m *MockSchoolGradeRepository) GetSchoolGradeByID(id uint) (models.SchoolGrade, error) {
	args := m.Called(id)
	if schoolGrade, ok := args.Get(0).(*models.SchoolGrade); ok {
		return *schoolGrade, args.Error(1)
	}
	return models.SchoolGrade{}, args.Error(1)
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