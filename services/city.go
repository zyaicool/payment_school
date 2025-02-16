package services

import (
	"encoding/json"
	"fmt"
	"os"
)

type CityService interface {
	GetDataCities(filename string, provinceID int) ([]DistrictNewMap, error)
}

type CityServiceImpl struct{}

func NewCityService() CityService {
	return &CityServiceImpl{}
}

type District struct {
	ID         int    `json:"id"`
	IDProvince int    `json:"id_province"`
	District   string `json:"district"`
}

type DistrictNewMap struct {
	ID         int    `json:"id"`
	IDProvince int    `json:"idProvince"`
	District   string `json:"district"`
}

func (cityService *CityServiceImpl) GetDataCities(filename string, provinceID int) ([]DistrictNewMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var districts []District
	if err := json.NewDecoder(file).Decode(&districts); err != nil {
		return nil, err
	}

	var filteredDistricts []DistrictNewMap
	for _, district := range districts {
		if district.IDProvince == provinceID {
			districtNew := DistrictNewMap{
				ID:         district.ID,
				IDProvince: district.IDProvince,
				District:   district.District,
			}
			filteredDistricts = append(filteredDistricts, districtNew)
		}
	}

	if len(filteredDistricts) == 0 {
		return nil, fmt.Errorf("No districts found for this province")
	}
	return filteredDistricts, nil
}
