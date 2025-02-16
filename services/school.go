package services

import (
	"fmt"
	"time"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SchoolServiceInterface interface {
	GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) (response.SchoolListResponse, error)
	GetSchoolByID(id uint) (*response.DetailSchoolResponse, error)
	CreateSchool(schoolRequest *request.SchoolCreateUpdateRequest, userID int) (*models.School, error)
	UpdateSchool(id uint, schoolRequest *request.SchoolCreateUpdateRequest, userID int) (*models.School, error)
	DeleteSchool(id uint, userID int) (*models.School, error)
	GenerateSchoolCode() (string, error)
	CheckAccessToSchool(c *fiber.Ctx) error
	GetAllOnboardingSchools(search string) ([]response.SchoolOnBoardingResponse, error)
}

type SchoolService struct {
	schoolRepository repositories.SchoolRepository
	userRepository   repositories.UserRepository
}

func NewSchoolService(schoolRepository repositories.SchoolRepository, userRepository repositories.UserRepository) SchoolServiceInterface {
	return &SchoolService{schoolRepository: schoolRepository, userRepository: userRepository}
}

func (schoolService *SchoolService) GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) (response.SchoolListResponse, error) {
	var mapSchool response.SchoolListResponse
	totalPage := 0
	mapSchool.Limit = limit
	mapSchool.Page = page

	if sortBy != "" {
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	// Get the paginated list of schools and the total data count
	listSchool, totalData, err := schoolService.schoolRepository.GetAllSchoolList(page, limit, search, sortBy, sortOrder)
	if err != nil {
		mapSchool.Data = []response.SchoolResponse{}
		return mapSchool, err
	}

	// Map models.School to response.SchoolResponse
	var schoolResponses []response.SchoolResponse
	for _, school := range listSchool {
		var provinceName string

		if utilities.HasNumeric(school.SchoolProvince) {
			dataProvince, err := GetProvinceById(school.SchoolProvince)
			if err != nil {
				provinceName = "-"
			}
			provinceName = dataProvince[0].Name
		} else {
			provinceName = school.SchoolProvince
		}

		schoolResponses = append(schoolResponses, response.SchoolResponse{
			ID:             int(school.ID),
			Npsn:           school.Npsn,
			SchoolName:     school.SchoolName,
			SchoolProvince: provinceName,
			SchoolPhone:    school.SchoolPhone,
			CreatedBy:      school.CreatedByUsername,
		})
	}

	if limit > 0 {
		totalPage = (totalData + limit - 1) / limit
	}

	mapSchool.TotalData = totalData
	mapSchool.TotalPage = totalPage
	mapSchool.Data = schoolResponses

	return mapSchool, nil
}

func (schoolService *SchoolService) GetSchoolByID(id uint) (*response.DetailSchoolResponse, error) {
	school, err := schoolService.schoolRepository.GetSchoolByID(id)
	if err != nil {
		return nil, err
	}

	schoolGradeID := uint(0)
	schoolGradeName := ""
	if school.SchoolGrade != nil {
		schoolGradeID = school.SchoolGrade.ID
		schoolGradeName = school.SchoolGrade.SchoolGradeName
	}

	schoolLogo := utilities.ConvertPath(school.SchoolLogo)
	SchoolLetterhead := utilities.ConvertPath(school.SchoolLetterhead)
	response := response.DetailSchoolResponse{
		ID:               int(school.ID),
		Npsn:             school.Npsn,
		SchoolCode:       school.SchoolCode,
		SchoolName:       school.SchoolName,
		SchoolProvince:   school.SchoolProvince,
		SchoolCity:       school.SchoolCity,
		SchoolPhone:      school.SchoolPhone,
		SchoolAddress:    school.SchoolAddress,
		SchoolMail:       school.SchoolMail,
		SchoolFax:        school.SchoolFax,
		SchoolLogo:       schoolLogo,
		SchoolGradeID:    schoolGradeID,
		SchoolGradeName:  schoolGradeName,
		SchoolLetterhead: SchoolLetterhead,
	}
	return &response, nil
}

func (schoolService *SchoolService) CreateSchool(schoolRequest *request.SchoolCreateUpdateRequest, userID int) (*models.School, error) {
	schoolCode, err := schoolService.GenerateSchoolCode()
	if err != nil {
		return nil, err
	}

	err = utilities.ValidateEmail(schoolRequest.SchoolMail)
	if err != nil {
		return nil, err
	}

	_, err = schoolService.schoolRepository.CheckNpsn(uint(schoolRequest.Npsn))
	if err == nil {
		return nil, fmt.Errorf("NPSN sudah terdaftar")
	}

	school := models.School{
		Npsn:             schoolRequest.Npsn,
		SchoolGradeID:    uint(schoolRequest.SchoolGradeID),
		SchoolCode:       schoolCode,
		SchoolName:       schoolRequest.SchoolName,
		SchoolProvince:   schoolRequest.SchoolProvince,
		SchoolCity:       schoolRequest.SchoolCity,
		SchoolAddress:    schoolRequest.SchoolAddress,
		SchoolPhone:      schoolRequest.SchoolPhone,
		SchoolMail:       schoolRequest.SchoolMail,
		SchoolFax:        schoolRequest.SchoolFax,
		SchoolLogo:       schoolRequest.SchoolLogo,
		SchoolLetterhead: schoolRequest.SchoolLetterhead,
	}
	school.Master.CreatedBy = userID
	school.Master.UpdatedBy = userID

	dataSchool, err := schoolService.schoolRepository.CreateSchool(&school)
	if err != nil {
		return nil, err
	}

	return dataSchool, nil
}

func (schoolService *SchoolService) UpdateSchool(id uint, schoolRequest *request.SchoolCreateUpdateRequest, userID int) (*models.School, error) {
	// Retrieve the school record
	getSchool, err := schoolService.schoolRepository.GetSchoolByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	// Validate email if provided
	if schoolRequest.SchoolMail != "" {
		err = utilities.ValidateEmail(schoolRequest.SchoolMail)
		if err != nil {
			return nil, err
		}
	}

	// Update only if fields are non-empty or non-zero
	_, err = schoolService.schoolRepository.CheckNpsnExistsExcept(uint(schoolRequest.Npsn), int(id))

	fmt.Print(err)

	if err == nil {
		return nil, fmt.Errorf("NPSN sudah terdaftar")
	}
	getSchool.Npsn = schoolRequest.Npsn

	if schoolRequest.SchoolName != "" {
		getSchool.SchoolName = schoolRequest.SchoolName
	}
	if schoolRequest.SchoolProvince != "" {
		getSchool.SchoolProvince = schoolRequest.SchoolProvince
	}
	if schoolRequest.SchoolCity != "" {
		getSchool.SchoolCity = schoolRequest.SchoolCity
	}
	if schoolRequest.SchoolAddress != "" {
		getSchool.SchoolAddress = schoolRequest.SchoolAddress
	}
	if schoolRequest.SchoolPhone != "" {
		getSchool.SchoolPhone = schoolRequest.SchoolPhone
	}
	if schoolRequest.SchoolMail != "" {
		getSchool.SchoolMail = schoolRequest.SchoolMail
	}
	if schoolRequest.SchoolFax != "" {
		getSchool.SchoolFax = schoolRequest.SchoolFax
	}
	if schoolRequest.SchoolLogo != "" {
		getSchool.SchoolLogo = schoolRequest.SchoolLogo
	}
	if schoolRequest.SchoolGradeID != 0 {
		getSchool.SchoolGradeID = uint(schoolRequest.SchoolGradeID)
	}
	if schoolRequest.SchoolLetterhead != "" {
		getSchool.SchoolLetterhead = schoolRequest.SchoolLetterhead
	}

	// Set the user ID for tracking who updated the record
	getSchool.Master.UpdatedBy = userID

	// Update the school record in the repository
	dataSchool, err := schoolService.schoolRepository.UpdateSchool(&getSchool)
	if err != nil {
		return nil, err
	}

	return dataSchool, nil
}

func (schoolService *SchoolService) DeleteSchool(id uint, userID int) (*models.School, error) {
	currentTime := time.Now()
	currentTimePointer := &currentTime
	getSchool, err := schoolService.schoolRepository.GetSchoolByID(id)
	if err != nil {
		return nil, fmt.Errorf("Data not found.")
	}

	getSchool.Master.DeletedAt = currentTimePointer
	getSchool.Master.DeletedBy = &userID
	dataSchool, err := schoolService.schoolRepository.UpdateSchool(&getSchool)
	if err != nil {
		return nil, err
	}

	return dataSchool, nil
}

func (schoolService *SchoolService) GenerateSchoolCode() (string, error) {
	lastNumber, err := schoolService.schoolRepository.GetLastSequenceNumberSchool()
	if err != nil {
		return "", err
	}

	newSequence := lastNumber + 1
	newCode := fmt.Sprintf("SC%03d", newSequence)

	return newCode, nil
}

func (schoolService *SchoolService) CheckAccessToSchool(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID == 1 || roleID == 5 {
		return nil
	}

	return fmt.Errorf("User can't access this page")
}

func (schoolService *SchoolService) GetAllOnboardingSchools(search string) ([]response.SchoolOnBoardingResponse, error) {

	schools, err := schoolService.schoolRepository.GetAllOnboardingSchools(search)
	if err != nil {
		return nil, err
	}

	var responses []response.SchoolOnBoardingResponse
	for _, school := range schools {
		schoolLogo := utilities.ConvertPath(school.SchoolLogo)
		responses = append(responses, response.SchoolOnBoardingResponse{
			ID:         school.ID,
			SchoolName: school.SchoolName,
			Npsn:       school.Npsn,
			SchoolLogo: schoolLogo,
		})
	}

	return responses, nil
}
