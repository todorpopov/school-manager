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
	validateLoginRequest(loginRequest *LoginRequest) *exceptions.AppError

	LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, *exceptions.AppError)
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

func (as *AuthService) LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, *exceptions.AppError) {
	if err := as.validateLoginRequest(loginRequest); err != nil {
		return nil, err
	}

	user, err := as.userSvc.GetUserByEmail(ctx, nil, loginRequest.Email)
	if err != nil {
		return nil, exceptions.NewAppError("USER_NOT_FOUND", "No user found with the provided email", err)
	}

	if match := as.bcryptSvc.PasswordsMatch(*user.Password, loginRequest.Password); !match {
		return nil, exceptions.NewAppError("INVALID_CREDENTIALS", "Invalid credentials", nil)
	}

	session, err := as.sessionSvc.CreateOrRenewSession(ctx, nil, user.UserId)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		SessionId: session.SessionId,
	}
	return resp, nil
}
