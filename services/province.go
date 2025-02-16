package services

import (
	"encoding/json"
	"fmt"
	"os"
)

// type ProvinceService struct {
// }

// func NewProvinceService() ProvinceService {
// 	return ProvinceService{}
// }

// ProvinceService defines the interface for province-related data retrieval
type ProvinceService interface {
	GetDataProvinces(filename string) ([]Region, error)
}

// ProvinceServiceImpl is the implementation of the ProvinceService interface
type ProvinceServiceImpl struct{}

// NewProvinceService creates a new instance of ProvinceServiceImpl
func NewProvinceService() ProvinceService {
	return &ProvinceServiceImpl{}
}

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RegionNewMap struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (provinceService *ProvinceServiceImpl) GetDataProvinces(filename string) ([]Region, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var regions []Region
	if err := json.NewDecoder(file).Decode(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}

func GetProvinceById(provinceID string) ([]RegionNewMap, error) {
	file, err := os.Open("data/mst_province.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var regions []Region
	if err := json.NewDecoder(file).Decode(&regions); err != nil {
		return nil, err
	}

	var filteredRegions []RegionNewMap
	for _, region := range regions {
		if region.ID == provinceID {
			regionNew := RegionNewMap{
				ID:   region.ID,
				Name: region.Name,
			}
			filteredRegions = append(filteredRegions, regionNew)
		}
	}

	if len(filteredRegions) == 0 {
		return nil, fmt.Errorf("No province found")
	}
	return filteredRegions, nil
}
