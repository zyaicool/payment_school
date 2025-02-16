package services_test

import (
	"schoolPayment/models"
	"schoolPayment/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRoleRepository is a mock of the RoleRepository interface
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
	return args.Get(0).(map[string]uint), args.Error(1)
}

func TestRoleService_GetAllRoles(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	service := services.NewRoleService(mockRepo)

	// Test case 1: Successfully retrieving roles
	mockRepo.On("GetAllRole", 1, 10, "", 0).Return([]models.Role{
		{Name: "Admin"},
		{Name: "User"},
	}, nil)

	resp, err := service.GetAllRoles(1, 10, "", 0)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "Admin", resp.Data[0].Name)
	assert.Equal(t, "User", resp.Data[1].Name)
	mockRepo.AssertCalled(t, "GetAllRole", 1, 10, "", 0)

	// Test case 2: No roles found
	mockRepo.On("GetAllRole", 1, 10, "search", 0).Return([]models.Role{}, nil)

	respEmpty, errEmpty := service.GetAllRoles(1, 10, "search", 0)

	assert.NoError(t, errEmpty)
	assert.Equal(t, 0, len(respEmpty.Data))
	mockRepo.AssertCalled(t, "GetAllRole", 1, 10, "search", 0)

	// Test case 3: Error while fetching roles
	mockRepo.On("GetAllRole", 1, 10, "error", 0).Return(nil, gorm.ErrInvalidTransaction)

	respErr, errErr := service.GetAllRoles(1, 10, "error", 0)

	assert.Error(t, errErr)
	assert.Nil(t, respErr.Data)  // Assert that Data is nil (not an empty slice)
	mockRepo.AssertCalled(t, "GetAllRole", 1, 10, "error", 0)
}

