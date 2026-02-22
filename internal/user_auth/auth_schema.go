package user_auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	SessionId string `json:"sessionId"`
}
