package request

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	RoleID        uint   `json:"role_id"`
	SchoolID	  uint	 `json:"school_id"`
	FirebaseToken string `json:"firebase_token"`
	jwt.RegisteredClaims
}
