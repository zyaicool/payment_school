package request

type UserCreateRequest struct {
	SchoolID uint   `json:"schoolId"`
	RoleID   uint   `json:"roleId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	RoleID   int    `json:"roleId"`
	Status   string `json:"status"`
}

type UpdateUserImageRequest struct {
	Image string `json:"image"`
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
