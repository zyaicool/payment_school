package repositories

import (
	database "schoolPayment/configs"
	models "schoolPayment/models"
)

func FindHistoryData(fileName string) (models.HistoryMasterData, error) {
	var masterData models.HistoryMasterData
	result := database.DB.Where("deleted_at IS NULL AND file_name = ?", fileName).First(&masterData)
	return masterData, result.Error
}

func CreateDataHistoryMasterData(masterData *models.HistoryMasterData) (*models.HistoryMasterData, error) {
	result := database.DB.Create(&masterData)
	return masterData, result.Error
}
