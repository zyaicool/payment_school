package response

import (
	"time"
)

type UserListResponse struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	TotalPage int                    `json:"totalPage"`
	TotalData int64                  `json:"totalData"`
	Data      []ListDataUserForIndex `json:"data"`
}

type DataUserForAuth struct {
	ID         int                  `json:"id"`
	Username   string               `json:"username"`
	Email      string               `json:"email"`
	Image      *string              `json:"image"`
	Role       DetailRole           `json:"role"`
	SchoolData DetailSchoolResponse `json:"schoolData"`
}

type ListDataUserForIndex struct {
	ID                int       `json:"id"`
	RoleID            int       `json:"roleId"`
	RoleName          string    `json:"roleName"`
	SchoolID          int       `json:"schoolId"`
	SchooolName       string    `json:"schoolName"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	CreatedDate       time.Time `json:"createdDate"`
	Status            string    `json:"status"`
	VerificationEmail bool      `json:"verificationEmail"`
}

type ResponseErrorUploadUser struct {
	Username string `json:"username"`
	Email    string `json:"Email"`
	Reason   string `json:"reason"`
}

type DetailUserResponse struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	RoleID            int    `json:"roleId"`
	RoleName          string `json:"roleName"`
	SchoolID          int    `json:"schoolId"`
	SchoolName        string `json:"schoolName"`
	VerificationEmail bool   `json:"verificationEmail"`
	Status            string `json:"status"`
	Image             string `json:"image"`
}

type GetUserByEmailAndSchoolIdDto struct {
	UserId   int `json:"userId"`
	SchoolId int `json:"schoolId"`
}
