package services

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempJSONFile(data []District) (string, error) {
	// Create a temporary file with os.CreateTemp
	tempFile, err := os.CreateTemp("", "districts_*.json")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Encode data to JSON and write to temp file
	if err := json.NewEncoder(tempFile).Encode(data); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func TestGetDataCities_Success(t *testing.T) {
	// Mock data
	mockDistricts := []District{
		{ID: 1, IDProvince: 1, District: "District1"},
		{ID: 2, IDProvince: 1, District: "District2"},
		{ID: 3, IDProvince: 2, District: "District3"},
	}

	// Create temp file with mock data
	filename, err := createTempJSONFile(mockDistricts)
	assert.NoError(t, err)
	defer os.Remove(filename) // Clean up the temp file after the test

	// Create service and call function
	service := NewCityService()
	districts, err := service.GetDataCities(filename, 1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, districts, 2)
	assert.Equal(t, "District1", districts[0].District)
	assert.Equal(t, "District2", districts[1].District)
}

func TestGetDataCities_NoDistrictsFound(t *testing.T) {
	// Mock data with a different province ID
	mockDistricts := []District{
		{ID: 1, IDProvince: 2, District: "District1"},
		{ID: 2, IDProvince: 3, District: "District2"},
	}

	// Create temp file with mock data
	filename, err := createTempJSONFile(mockDistricts)
	assert.NoError(t, err)
	defer os.Remove(filename)

	// Create service and call function
	service := NewCityService()
	districts, err := service.GetDataCities(filename, 1)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "No districts found for this province")
	assert.Nil(t, districts)
}

func TestGetDataCities_FileNotFound(t *testing.T) {
	// Create service and call function with a non-existent file
	service := NewCityService()
	districts, err := service.GetDataCities("non_existent_file.json", 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, districts)
}
