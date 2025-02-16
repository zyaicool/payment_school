package models

import (
	"encoding/json"
	"os"
)

type BillingStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// LoadBanks reads the bank data from a JSON file
func LoadBillingStatus(filePath string) ([]BillingStatus, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var billingStatuses []BillingStatus
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&billingStatuses); err != nil {
		return nil, err
	}

	return billingStatuses, nil
}
