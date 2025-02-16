package services

import (
	"fmt"

	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"
)

func GetAllGuardian(page int, limit int, search string) (response.StudentGuardianResponse, error) {
	var mapGuardian response.StudentGuardianResponse
	mapGuardian.Limit = limit
	mapGuardian.Page = page

	listGuardian, err := repositories.GetAllGuardian(page, limit, search)
	if err != nil {
		mapGuardian.Data = nil
		return mapGuardian, nil
	}

	mapGuardian.Data = listGuardian

	return mapGuardian, nil
}

func GetGuardianByID(id uint) (models.StudentGuardian, error) {
	return repositories.GetGuardianByID(id)
}

func CreateGuardian(guardian *models.StudentGuardian, userID int) (*models.StudentGuardian, error) {
	// check format email
	if guardian.GuardianMail != "" {
		err := utilities.ValidateEmail(guardian.GuardianMail)
		if err != nil {
			return nil, err
		}
	}

	// getUser, err := repositories.GetUserByID(uint(userID))
	// if err != nil {
	// 	return nil, err
	// }

	// if getUser.RoleID != 2 {
	// 	return nil, fmt.Errorf("Silahkan login dengan akun orang tua atau wali siswa.")
	// }

	// set user id to created by and updated by
	guardian.UserID = uint(userID)
	guardian.Master.CreatedBy = userID
	guardian.Master.UpdatedBy = userID

	dataGuardian, err := repositories.CreateGuardian(guardian)
	if err != nil {
		return nil, err
	}

	return dataGuardian, nil
}

func UpdateGuardian(id uint, guardian *models.StudentGuardian, userID int) (*models.StudentGuardian, error) {
	getGuardian, err := repositories.GetGuardianByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	// check format email
	if guardian.GuardianMail != "" {
		err = utilities.ValidateEmail(guardian.GuardianMail)
		if err != nil {
			return nil, err
		}
	}

	// getUser, err := repositories.GetUserByID(uint(userID))
	// if err != nil {
	// 	return nil, err
	// }

	// if getUser.RoleID != 2 {
	// 	return nil, fmt.Errorf("Silahkan login dengan akun orang tua atau wali siswa.")
	// }

	getGuardian.GuardianName = guardian.GuardianName
	getGuardian.GuardianAddress = guardian.GuardianAddress
	getGuardian.GuardianHandphone = guardian.GuardianHandphone
	getGuardian.GuardianMail = guardian.GuardianMail
	getGuardian.GuardianCitizenship = guardian.GuardianCitizenship
	getGuardian.GuardianSalary = guardian.GuardianSalary
	getGuardian.Master.UpdatedBy = userID

	dataGuardian, err := repositories.UpdateGuardian(&getGuardian)
	if err != nil {
		return nil, err
	}

	return dataGuardian, nil
}

func GetGuardianByUserLogin(id uint) (*models.StudentGuardian, error) {
	// getUser, err := repositories.GetUserByID(uint(id))
	// if err != nil {
	// 	return nil, err
	// }

	// if getUser.RoleID != 2 {
	// 	return nil, fmt.Errorf("Silahkan login dengan akun orang tua atau wali siswa.")
	// }

	return repositories.GetGuardianByUserLogin(id)
}
