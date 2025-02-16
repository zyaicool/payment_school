package models

import (
	"encoding/json"
	"os"
)

type Bank struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Code  string `json:"code"`
}

// LoadBanks reads the bank data from a JSON file
func LoadBanks(filePath string) ([]Bank, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var banks []Bank
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&banks); err != nil {
		return nil, err
	}

	return banks, nil
}
