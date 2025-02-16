// Code generated by MockGen. DO NOT EDIT.
// Source: D:\mas ikul\work\PT AIGEN\school biller\school-payment-backend\repositories\roleRepository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	models "schoolPayment/models"

	gomock "github.com/golang/mock/gomock"
)

// MockRoleRepository is a mock of RoleRepository interface.
type MockRoleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRoleRepositoryMockRecorder
}

// MockRoleRepositoryMockRecorder is the mock recorder for MockRoleRepository.
type MockRoleRepositoryMockRecorder struct {
	mock *MockRoleRepository
}

// NewMockRoleRepository creates a new mock instance.
func NewMockRoleRepository(ctrl *gomock.Controller) *MockRoleRepository {
	mock := &MockRoleRepository{ctrl: ctrl}
	mock.recorder = &MockRoleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoleRepository) EXPECT() *MockRoleRepositoryMockRecorder {
	return m.recorder
}

// GetAllRole mocks base method.
func (m *MockRoleRepository) GetAllRole(page, limit int, search string, roleID int) ([]models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRole", page, limit, search, roleID)
	ret0, _ := ret[0].([]models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRole indicates an expected call of GetAllRole.
func (mr *MockRoleRepositoryMockRecorder) GetAllRole(page, limit, search, roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRole", reflect.TypeOf((*MockRoleRepository)(nil).GetAllRole), page, limit, search, roleID)
}
