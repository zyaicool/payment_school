package mocks

import (
	"schoolPayment/models"

	"github.com/stretchr/testify/mock"
)

type RoleData struct {
	mock.Mock
}

func (_m *RoleData) GetAllData(page int, limit int, search string, roleID int) ([]models.Role, error) {
	ret := _m.Called(page, limit, search, roleID)

	var r0 []models.Role
	if rf, ok := ret.Get(0).(func(int, int, string, int) []models.Role); ok {
		r0 = rf(page, limit, search, roleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, string, int) error); ok {
		r1 = rf(page, limit, search, roleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
