package services

import (
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBillingStudentRepository struct {
	mock.Mock
}

func (m *MockBillingStudentRepository) GetAllBillingStudent(page int, limit int, search string, studentId int, roleID int, schoolGradeID int, paymentType string, schoolID int, sortBy string, sortOrder string, schoolClassID int) ([]response.DataListBillingPerStudent, int, int64, error) {
	args := m.Called(page, limit, search, studentId, roleID, schoolGradeID, paymentType, schoolID, sortBy, sortOrder, schoolClassID)
	return args.Get(0).([]response.DataListBillingPerStudent), args.Int(1), args.Get(2).(int64), args.Error(3)
}

func (m *MockBillingStudentRepository) GetDetailBillingStudentByID(billingStudentId int) (models.BillingStudent, error) {
	args := m.Called(billingStudentId)
	return args.Get(0).(models.BillingStudent), args.Error(1)
}

func (m *MockBillingStudentRepository) UpdateBillingStudent(billingStudent *models.BillingStudent) error {
	args := m.Called(billingStudent)
	return args.Error(0)
}

func (m *MockBillingStudentRepository) DeleteBillingStudent(billingStudentId int, deletedBy int) error {
	args := m.Called(billingStudentId, deletedBy)
	return args.Error(0)
}

func (m *MockBillingStudentRepository) GetTotalAmountByBillingStudentIds(billingStudentIds []string) (int, error) {
	args := m.Called(billingStudentIds)
	return args.Int(0), args.Error(1)
}

func (m *MockBillingStudentRepository) GetListBillingId(billingStudentIds []string) ([]int, error) {
	args := m.Called(billingStudentIds)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockBillingStudentRepository) CheckBillingStudentExists(studentID, billingID, billingDetailID uint) (bool, error) {
	args := m.Called(studentID, billingID, billingDetailID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBillingStudentRepository) GetBillingDetailsByIDs(ids []uint) ([]*models.BillingDetail, error) {
	args := m.Called(ids)
	return args.Get(0).([]*models.BillingDetail), args.Error(1)
}

func (m *MockBillingStudentRepository) CheckBulkBillingStudentExists(studentID uint, detailIDs []uint) ([]models.BillingStudent, error) {
	args := m.Called(studentID, detailIDs)
	return args.Get(0).([]models.BillingStudent), args.Error(1)
}

func (m *MockBillingStudentRepository) BulkCreateBillingStudents(billingStudents []models.BillingStudent) ([]models.BillingStudent, error) {
	args := m.Called(billingStudents)
	return args.Get(0).([]models.BillingStudent), args.Error(1)
}

func (m *MockBillingStudentRepository) GetInstallmentDetails(studentId int, schoolId uint) ([]response.BillingStudentByStudentIDDetailResponse, []response.DonationBillingResponse, error) {
	args := m.Called(studentId, schoolId)
	return args.Get(0).([]response.BillingStudentByStudentIDDetailResponse), args.Get(1).([]response.DonationBillingResponse), args.Error(2)
}

func TestGetDetailBillingIDService_Success(t *testing.T) {
	// Mock repositories
	mockStudentRepo := new(MockStudentRepository)
	mockSchoolYearRepo := new(MockSchoolYearRepository)
	mockSchoolGradeRepo := new(MockSchoolGradeRepository)
	mockSchoolClassRepo := new(MockSchoolClassRepository)
	mockBillingStudentRepo := new(MockBillingStudentRepository)

	// Mock dependencies
	mockUser := models.User{UserSchool: &models.UserSchool{School: &models.School{ID: 1}}}
	studentID := 1

	mockStudent := models.Student{
		Nis:          "12345",
		FullName:     "John Doe",
		SchoolYearID: 1,
		SchoolGradeID: 2,
		SchoolClassID: 3,
	}

	mockSchoolYear := models.SchoolYear{
		SchoolYearName: "2023/2024",
	}

	mockSchoolGrade := models.SchoolGrade{
		SchoolGradeName: "Grade 10",
	}

	mockSchoolClass := models.SchoolClass{
		SchoolClassName: "Class A",
	}

	mockInstallmentBillings := []response.BillingStudentByStudentIDDetailResponse{
		{BillingStudentID: 1, Amount: 500000},
	}

	mockDonations := []response.DonationBillingResponse{
		{BillingID: 1, BillingName: "Billing Name"},
	}

	// Mock expectations
	mockStudentRepo.On("GetStudentByID", uint(studentID), mockUser).Return(mockStudent, nil).Once()
	mockSchoolYearRepo.On("GetSchoolYearByID", mockStudent.SchoolYearID).Return(&models.SchoolYear{SchoolYearName: "2023/2024"}, nil).Once()
	mockSchoolGradeRepo.On("GetSchoolGradeByID", mockStudent.SchoolGradeID).Return(&models.SchoolGrade{SchoolGradeName: "Grade 10"}, nil).Once()
	mockSchoolClassRepo.On("GetSchoolClassByID", mockStudent.SchoolClassID).Return(mockSchoolClass, nil).Once()
	mockBillingStudentRepo.On("GetInstallmentDetails", 0, uint(1)).
    Return([]response.BillingStudentByStudentIDDetailResponse{{BillingStudentID: 1, Amount: 500000}}, 
		   []response.DonationBillingResponse{{BillingID: 1, BillingName: "Billing Name"}}, nil).Once()

	// Initialize service
	billingStudentService := BillingStudentService{
		studentRepository:       mockStudentRepo,
		schoolYearRepository:    mockSchoolYearRepo,
		schoolGradeRepository:   mockSchoolGradeRepo,
		schoolClassRepository:   mockSchoolClassRepo,
		billingStudentRepository: mockBillingStudentRepo,
	}

	// Call the method
	result, err := billingStudentService.GetDetailBillingIDService(1, mockUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockStudent.ID, result.ID)
	assert.Equal(t, mockStudent.Nis, result.Nis)
	assert.Equal(t, mockStudent.FullName, result.StudentName)
	assert.Equal(t, mockSchoolGrade.SchoolGradeName, result.SchoolGradeName)
	assert.Equal(t, mockSchoolClass.SchoolClassName, result.SchoolClassName)
	assert.Equal(t, mockSchoolYear.SchoolYearName, result.LatestSchoolYear)
	assert.Equal(t, mockInstallmentBillings, result.ListBilling)
	assert.Equal(t, mockDonations, result.ListDonation)

	// Verify mock expectations
	mockStudentRepo.AssertExpectations(t)
	mockSchoolYearRepo.AssertExpectations(t)
	mockSchoolGradeRepo.AssertExpectations(t)
	mockSchoolClassRepo.AssertExpectations(t)
	mockBillingStudentRepo.AssertExpectations(t)
}
