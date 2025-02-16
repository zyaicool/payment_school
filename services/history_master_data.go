package services

import (
	"fmt"
	"time"

	"schoolPayment/models"
	"schoolPayment/repositories"
)

func CreateDataHistoryMasterData(fileName string) error {
	masterData := models.HistoryMasterData{
		GenerateDate: time.Now(),
		FileName:     fileName,
	}
	_, err := repositories.CreateDataHistoryMasterData(&masterData)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func FindHistoryData(fileName string) (models.HistoryMasterData, error) {
	return repositories.FindHistoryData(fileName)
}
