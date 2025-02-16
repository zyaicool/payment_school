package repositories

import (
	database "schoolPayment/configs"
	"schoolPayment/models"
)

func CreateBillingDetail(billingDetail *models.BillingDetail, studentNis string, isMaster bool) (*models.BillingDetail, error) {
	var existingBillingDetail models.BillingDetail

	if isMaster {
		result := database.DB.Create(&billingDetail)
		return billingDetail, result.Error
	} else {
		err := database.DB.Where("billing_id = ? AND detail_billing_name LIKE ?", billingDetail.BillingID, "%"+studentNis+"%").First(&existingBillingDetail).Error
		if err == nil {
			result := database.DB.Model(&existingBillingDetail).Updates(billingDetail)
			return &existingBillingDetail, result.Error
		} else if err.Error() == "record not found" {
			result := database.DB.Create(&billingDetail)
			return billingDetail, result.Error
		} else {
			return nil, err
		}
	}

}

func GetBillingDetailByID(id uint) (*models.BillingDetail, error) {
	var billingDetail models.BillingDetail
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&billingDetail)
	return &billingDetail, result.Error
}
