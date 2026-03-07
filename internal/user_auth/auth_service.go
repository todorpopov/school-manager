package user_auth

import (
	"context"
	"slices"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IAuthService interface {
	RegisterUser(ctx context.Context, registerRequest *RegisterRequest) (*AuthResponse, *exceptions.AppError)
	LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*AuthResponse, *exceptions.AppError)
	IsRequestAuthorized(ctx context.Context, request *AuthRequest) (bool, *exceptions.AppError)
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

func (as *AuthService) RegisterUser(ctx context.Context, registerRequest *RegisterRequest) (*AuthResponse, *exceptions.AppError) {
	if err := ValidateRegisterRequest(registerRequest); err != nil {
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
	if err := ValidateLoginRequest(loginRequest); err != nil {
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

func (as *AuthService) IsRequestAuthorized(ctx context.Context, request *AuthRequest) (bool, *exceptions.AppError) {
	if err := ValidateAuthRequest(request); err != nil {
		return false, err
	}

	session, err := as.sessionSvc.GetActiveSessionById(ctx, nil, request.SessionId)
	if err != nil {
		return false, err
	}

	rolesMap, err := as.userSvc.GetUsersRoles(ctx, []int32{session.UserId})
	if err != nil {
		return false, err
	}

	userRoles := rolesMap[session.UserId]

	for _, role := range request.RequiredRoles {
		if slices.Contains(userRoles, role) {
			return true, nil
		}
	}
	return false, exceptions.NewAppError("UNAUTHORIZED", "Unauthorized", nil)
}
