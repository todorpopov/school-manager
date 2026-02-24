package user_auth

import (
	"context"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IAuthService interface {
	validateRegisterRequest(registerRequest *RegisterRequest) *exceptions.AppError
	validateLoginRequest(loginRequest *LoginRequest) *exceptions.AppError

	RegisterUser(ctx context.Context, registerRequest *RegisterRequest) (*AuthResponse, *exceptions.AppError)
	LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*AuthResponse, *exceptions.AppError)
}

type AuthService struct {
	bcryptSvc  internal.IBCryptService
	userSvc    users.IUserService
	sessionSvc sessions.ISessionService
}

func NewAuthService(bcryptSvc internal.IBCryptService, userSvc users.IUserService, sessionSvc sessions.ISessionService) *AuthService {
	return &AuthService{
		bcryptSvc:  bcryptSvc,
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
	}
}

func (as *AuthService) validateRegisterRequest(registerRequest *RegisterRequest) *exceptions.AppError {
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

func (as *AuthService) validateLoginRequest(loginRequest *LoginRequest) *exceptions.AppError {
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

func (as *AuthService) RegisterUser(ctx context.Context, registerRequest *RegisterRequest) (*AuthResponse, *exceptions.AppError) {
	if err := as.validateRegisterRequest(registerRequest); err != nil {
		return nil, err
	}

	createUser := &users.CreateUser{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  registerRequest.Password,
	}
	user, err := as.userSvc.CreateUser(ctx, nil, createUser)
	if err != nil {
		return nil, err
	}

	session, err := as.sessionSvc.CreateOrRenewSession(ctx, nil, user.UserId)
	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{
		SessionId: session.SessionId,
	}
	return resp, nil
}

func (as *AuthService) LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*AuthResponse, *exceptions.AppError) {
	if err := as.validateLoginRequest(loginRequest); err != nil {
		return nil, err
	}

	user, err := as.userSvc.GetUserByEmailWithPass(ctx, nil, loginRequest.Email)
	if err != nil {
		return nil, err
	}

	if match := as.bcryptSvc.PasswordsMatch(*user.Password, loginRequest.Password); !match {
		return nil, exceptions.NewAppError("INVALID_CREDENTIALS", "Invalid credentials", nil)
	}

	session, err := as.sessionSvc.CreateOrRenewSession(ctx, nil, user.UserId)
	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{
		SessionId: session.SessionId,
	}
	return resp, nil
}
