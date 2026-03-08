package user_auth

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Roles     []string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	SessionId string `json:"sessionId"`
}

type AuthRequest struct {
	SessionId     string   `json:"sessionId"`
	RequiredRoles []string `json:"requiredRoles"`
}

func ValidateRegisterRequest(registerRequest *RegisterRequest) *exceptions.AppError {
	messages := map[string]string{}

	if msg := domain_model.ValidateString(&registerRequest.FirstName, 1, 255, true); msg != "" {
		messages["first_name"] = msg
	}

	if msg := domain_model.ValidateString(&registerRequest.LastName, 1, 255, true); msg != "" {
		messages["last_name"] = msg
	}

	if msg := domain_model.ValidateEmail(&registerRequest.Email, true); msg != "" {
		messages["email"] = msg
	}

	if msg := domain_model.ValidatePassword(&registerRequest.Password, true); msg != "" {
		messages["password"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during registration", messages)
	}
	return nil
}

func ValidateLoginRequest(loginRequest *LoginRequest) *exceptions.AppError {
	messages := map[string]string{}

	if msg := domain_model.ValidateEmail(&loginRequest.Email, true); msg != "" {
		messages["email"] = msg
	}

	if msg := domain_model.ValidatePassword(&loginRequest.Password, true); msg != "" {
		messages["password"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during login", messages)
	}
	return nil
}

func ValidateAuthRequest(authRequest *AuthRequest) *exceptions.AppError {
	messages := domain_model.ValidateRoleNames(authRequest.RequiredRoles, true)
	if len(messages) > 0 {
		return exceptions.NewValidationError("Invalid roles", messages)
	}
	return nil
}
