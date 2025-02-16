package services

import (
	"fmt"
	"time"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
)

type BillingTypeService struct {
	billingRepo repositories.BillingTypeRepository
}

func NewBillingTypeService(billingRepo repositories.BillingTypeRepository) BillingTypeService {
	return BillingTypeService{billingRepo: billingRepo}
}

func (billingTypeService *BillingTypeService) GetAllBillingType(page int, limit int, search string) (response.BillingTypeListResponse, error) {
	var mapBillingType response.BillingTypeListResponse
	mapBillingType.Limit = limit
	mapBillingType.Page = page

	listBillingType, err := billingTypeService.billingRepo.GetAllBillingType(page, limit, search)
	if err != nil {
		mapBillingType.Data = []models.BillingType{}
		return mapBillingType, nil
	}
	if len(listBillingType) > 0 {
		mapBillingType.Data = listBillingType
	} else {
		mapBillingType.Data = []models.BillingType{}
	}

	return mapBillingType, nil
}

func (billingTypeService *BillingTypeService) GetBillingTypeByID(id uint) (models.BillingType, error) {
	return billingTypeService.billingRepo.GetBillingTypeByIDTest(id)
}

func (billingTypeService *BillingTypeService) CreateBillingType(billingTypeRequest *request.BillingTypeCreateUpdateRequest, userID int) (*models.BillingType, error) {
	// di atas sini ada proses get data school untuk kebutuham generate billing type code
	billingTypeCode, err := billingTypeService.GenerateBillingTypeCode()
	if err != nil {
		return nil, err
	}

	schoolID := 0
	if billingTypeRequest.SchoolID == 0 {
		getSchoolId, err := billingTypeService.billingRepo.GetUserSchoolByUserIdTest(uint(userID))
		if err != nil {
			return nil, err
		}
		schoolID = int(getSchoolId.SchoolID)
	} else {
		schoolID = billingTypeRequest.SchoolID
	}

	billingType := models.BillingType{
		SchoolID:          uint(schoolID),
		BillingTypeCode:   billingTypeCode,
		BillingTypeName:   billingTypeRequest.BillingTypeName,
		BillingTypePeriod: billingTypeRequest.BillingTypePeriod,
	}
	billingType.Master.CreatedBy = userID
	billingType.Master.UpdatedBy = userID

	dataBillingType, err := billingTypeService.billingRepo.CreateBillingType(&billingType)
	if err != nil {
		return nil, err
	}

	return dataBillingType, nil
}

func (billingTypeService *BillingTypeService) UpdateBillingType(id uint, billingTypeRequest *request.BillingTypeCreateUpdateRequest, userID int) (*models.BillingType, error) {
	getBillingType, err := billingTypeService.billingRepo.GetBillingTypeByIDTest(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getBillingType.BillingTypeName = billingTypeRequest.BillingTypeName
	getBillingType.BillingTypePeriod = billingTypeRequest.BillingTypePeriod
	getBillingType.Master.UpdatedBy = userID

	dataBillingType, err := billingTypeService.billingRepo.UpdateBillingType(&getBillingType)
	if err != nil {
		return nil, err
	}

	return dataBillingType, nil
}

func (billingTypeService *BillingTypeService) DeleteBillingType(id uint, userID int) (*models.BillingType, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getBillingType, err := billingTypeService.billingRepo.GetBillingTypeByIDTest(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getBillingType.Master.DeletedAt = currentTimePointer
	getBillingType.Master.DeletedBy = &userID
	dataBillingType, err := billingTypeService.billingRepo.UpdateBillingType(&getBillingType)
	if err != nil {
		return nil, err
	}

	return dataBillingType, nil
}

func (billingTypeService *BillingTypeService) GenerateBillingTypeCode() (string, error) {
	lastNumber, err := billingTypeService.billingRepo.GetLastSequenceNumberBillingType()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1
	newCode := fmt.Sprintf("BT%03d", newSequence)

	return newCode, nil
}
