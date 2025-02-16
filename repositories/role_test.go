package repositories

import (
	"errors"
	"schoolPayment/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock struct untuk roleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetAllRole(page int, limit int, search string, roleID int) ([]models.Role, error) {
	args := m.Called(page, limit, search, roleID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) GetRolesByNames(roleNames []string) (map[string]uint, error) {
	args := m.Called(roleNames)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]uint), args.Error(1)
	}
	return nil, args.Error(1)
}

// Test suite struct
type RoleRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockRoleRepository
}

func (suite *RoleRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(MockRoleRepository)
}

// Test Success untuk GetAllRole
func (suite *RoleRepositoryTestSuite) TestGetAllRole_Success() {
	dummyRoles := []models.Role{
		{Name: "Admin"},
		{Name: "User"},
	}

	// Setup mock untuk GetAllRole
	suite.mockRepo.On("GetAllRole", 1, 10, "", 0).Return(dummyRoles, nil)

	// Memanggil metode yang diuji
	result, err := suite.mockRepo.GetAllRole(1, 10, "", 0)

	// Assert hasil
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyRoles, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error untuk GetAllRole
func (suite *RoleRepositoryTestSuite) TestGetAllRole_Error() {
	// Setup mock untuk GetAllRole dengan error
	suite.mockRepo.On("GetAllRole", 1, 10, "", 0).Return(nil, errors.New("database error"))

	// Memanggil metode yang diuji
	result, err := suite.mockRepo.GetAllRole(1, 10, "", 0)

	// Assert error dan pastikan hasil nil
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result) // Hasil harus nil saat terjadi error
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Success untuk GetRolesByNames
func (suite *RoleRepositoryTestSuite) TestGetRolesByNames_Success() {
	roleNames := []string{"Admin", "User"}
	dummyRoleMap := map[string]uint{
		"Admin": 1,
		"User":  2,
	}

	// Setup mock untuk GetRolesByNames
	suite.mockRepo.On("GetRolesByNames", roleNames).Return(dummyRoleMap, nil)

	// Memanggil metode yang diuji
	result, err := suite.mockRepo.GetRolesByNames(roleNames)

	// Assert hasil
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyRoleMap, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error untuk GetRolesByNames
func (suite *RoleRepositoryTestSuite) TestGetRolesByNames_Error() {
	roleNames := []string{"Admin", "User"}

	// Setup mock untuk GetRolesByNames dengan error
	suite.mockRepo.On("GetRolesByNames", roleNames).Return(nil, errors.New("database error"))

	// Memanggil metode yang diuji
	result, err := suite.mockRepo.GetRolesByNames(roleNames)

	// Assert error dan pastikan hasil nil
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result) // Hasil harus nil saat terjadi error
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}
