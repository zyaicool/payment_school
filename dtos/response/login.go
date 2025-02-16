package response

type LoginResponse struct {
	Token string `json:"token"`
}

type InvalidateFirebaseTokenResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}