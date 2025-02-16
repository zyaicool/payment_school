package services

import (
	"schoolPayment/models"

	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) GetAllStudentForBilling(user models.User, schoolGradeId int, schoolClassIds []int) ([]models.Student, error) {
	args := m.Called(user, schoolGradeId, schoolClassIds)
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepository) CreateStudent(student *models.Student) (*models.Student, error) {
	args := m.Called(student)
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) CreateStudentHistory(studentHistory *models.StudentHistory) (*models.StudentHistory, error) {
	args := m.Called(studentHistory)
	return args.Get(0).(*models.StudentHistory), args.Error(1)
}

func (m *MockStudentRepository) UpdateStudent(student *models.Student) (*models.Student, error) {
	args := m.Called(student)
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) CreateUserStudentRepository(userStudent *models.UserStudent) (*models.UserStudent, error) {
	args := m.Called(userStudent)
	return args.Get(0).(*models.UserStudent), args.Error(1)
}

func (m *MockStudentRepository) GetStudentByUserIdRepository(id uint) ([]models.Student, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepository) BulkCreateStudents(students []models.Student) error {
	args := m.Called(students)
	return args.Error(0)
}

func (m *MockStudentRepository) BulkUpdateStudents(students []models.Student) error {
	args := m.Called(students)
	return args.Error(0)
}

func (m *MockStudentRepository) GetStudentByNis(nis string, user models.User) (models.Student, error) {
	args := m.Called(nis, user)
	return args.Get(0).(models.Student), args.Error(1)
}

func (m *MockStudentRepository) GetStudentsByNIS(nisNumbers []string, user models.User) ([]models.Student, error) {
	args := m.Called(nisNumbers, user)
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepository) BulkCheckUserStudentExists(pairs []struct {
	UserID    uint
	StudentID uint
}) (map[string]bool, error) {
	args := m.Called(pairs)
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockStudentRepository) BulkCreateUserStudents(pairs []models.UserStudent) error {
	args := m.Called(pairs)
	return args.Error(0)
}

func (m *MockStudentRepository) BulkCreateStudentHistory(histories []models.StudentHistory) error {
	args := m.Called(histories)
	return args.Error(0)
}

func (m *MockStudentRepository) CreateImageStudentRepository(student *models.Student, id int) (*models.Student, error) {
	args := m.Called(student, id)
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) GetStudentByID(id uint, user models.User) (models.Student, error) {
	args := m.Called(id, user)
	return args.Get(0).(models.Student), args.Error(1)
}