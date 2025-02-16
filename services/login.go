package services

import (
	"fmt"
	"regexp"
	"time"

	"schoolPayment/constants"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	repositories "schoolPayment/repositories"
	utilities "schoolPayment/utilities"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepository       repositories.UserRepository
	auditTrailRepository repositories.AuditTrailRepository
}

func NewLoginService(userRepository repositories.UserRepository, auditTrailRepository repositories.AuditTrailRepository) LoginService {
	return LoginService{userRepository: userRepository, auditTrailRepository: auditTrailRepository}
}

// LoginService validates the user's credentials and generates a JWT token
func (loginService *LoginService) LoginService(email string, password string, firebaseToken string, schoolId uint) (string, error) {
	var user models.User
	var err error
	platform := "website"
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailPattern)
	if re.MatchString(email) {
		user, err = repositories.GetUserByEmail(email)
		if err != nil {
			return "", fmt.Errorf(constants.EmailPasswordSalahMessage)
		}
		if user.UserSchool != nil && schoolId != 0 && user.UserSchool.SchoolID != schoolId {
			return "", fmt.Errorf(constants.EmailPasswordSalahMessage)
		}
	} else {

		user, err = repositories.GetUserByUsername(email)
		if err != nil {
			return "", fmt.Errorf(constants.EmailPasswordSalahMessage)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf(constants.EmailPasswordSalahMessage)
	}

	// Step 3: Generate JWT token
	token, err := utilities.GenerateJWT(user, firebaseToken)
	if err != nil {
		return "", fmt.Errorf("Oops something wrong.")
	}

	if firebaseToken != "" {
		platform = "mobile"

		if user.RoleID != 2 {
			return "", fmt.Errorf("Oops something wrong.")
		}
	}

	// save to audit trail
	newAuditTrail := models.AuditTrail{
		UserID:     user.ID,
		Email:      user.Email,
		Role:       user.Role.Name,
		UserAction: "Login",
		ApiPath:    "/login",
		LogTime:    time.Now(),
		Platform:   platform,
		FirebaseID: firebaseToken,
		IsValid:    true,
	}

	newAuditTrail.CreatedBy = int(user.ID)
	newAuditTrail.UpdatedBy = int(user.ID)
	err = loginService.auditTrailRepository.CreateDataAuditTrail(&newAuditTrail)
	return token, nil
}

func (loginService *LoginService) GetUserFromAuth(userID int) (response.DataUserForAuth, error) {
	var detailDataUser response.DataUserForAuth
	var roleMatrix []response.DetailRoleMatrix
	var school response.DetailSchoolResponse

	user, err := loginService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return detailDataUser, fmt.Errorf("User Not Found.")
	}
	for _, rm := range user.Role.RoleMatrix {
		roleMatrix = append(roleMatrix, response.DetailRoleMatrix{
			PageName: rm.PageName,
			PageCode: rm.PageCode,
			IsCreate: rm.IsCreate,
			IsRead:   rm.IsRead,
			IsUpdate: rm.IsUpdate,
			IsDelete: rm.IsDelete,
		})
	}

	role := response.DetailRole{
		ID:         int(user.RoleID),
		RoleName:   user.Role.Name,
		RoleMatrix: roleMatrix,
	}

	schoolGradeID := uint(0)
	schoolGradeName := ""
	if user.UserSchool != nil {
		schoolGradeID = user.UserSchool.School.SchoolGrade.ID
		schoolGradeName = user.UserSchool.School.SchoolGrade.SchoolGradeName
	}

	if user.UserSchool != nil {
		school = response.DetailSchoolResponse{
			ID:              int(user.UserSchool.SchoolID),
			SchoolName:      user.UserSchool.School.SchoolName,
			SchoolCode:      user.UserSchool.School.SchoolCode,
			SchoolProvince:  user.UserSchool.School.SchoolProvince,
			SchoolCity:      user.UserSchool.School.SchoolCity,
			SchoolPhone:     user.UserSchool.School.SchoolPhone,
			SchoolAddress:   user.UserSchool.School.SchoolAddress,
			SchoolMail:      user.UserSchool.School.SchoolMail,
			SchoolFax:       user.UserSchool.School.SchoolFax,
			SchoolLogo:      user.UserSchool.School.SchoolLogo,
			SchoolGradeID:   schoolGradeID,
			SchoolGradeName: schoolGradeName,
		}
	}

	image := utilities.ConvertPathImage(user.Image)
	detailDataUser = response.DataUserForAuth{
		ID:         int(user.ID),
		Username:   user.Username,
		Email:      user.Email,
		Image:      image,
		Role:       role,
		SchoolData: school,
	}

	return detailDataUser, nil
}

func (loginService *LoginService) UpdateFirebaseTokenStatus(firebaseToken string, userId int) error {
	return loginService.auditTrailRepository.InvalidateFirebaseToken(firebaseToken, userId)
}
