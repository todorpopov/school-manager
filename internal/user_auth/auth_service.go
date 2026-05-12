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
	RegisterAdminUser(ctx context.Context, registerAdminReq *RegisterAdminRequest) (*AuthResponse, *exceptions.AppError)
	LogUserIn(ctx context.Context, loginRequest *LoginRequest) (*AuthResponse, *exceptions.AppError)
	SetSessionRole(ctx context.Context, sessionId string, selectRoleReq *SelectRoleRequest) (*AuthResponse, *exceptions.AppError)
	IsRequestAuthorized(ctx context.Context, request *AuthRequest) (bool, *exceptions.AppError)
}

type AuthService struct {
	bcryptSvc       internal.IBCryptService
	userSvc         users.IUserService
	sessionSvc      sessions.ISessionService
	systemAuthToken string
}

func NewAuthService(
	bcryptSvc internal.IBCryptService,
	userSvc users.IUserService,
	sessionSvc sessions.ISessionService,
	systemAuthToken string,
) *AuthService {
	return &AuthService{
		bcryptSvc:       bcryptSvc,
		userSvc:         userSvc,
		sessionSvc:      sessionSvc,
		systemAuthToken: systemAuthToken,
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
		Roles:     registerRequest.Roles,
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
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SessionId: session.SessionId,
		Roles:     user.Roles,
	}
	return resp, nil
}

func (as *AuthService) RegisterAdminUser(ctx context.Context, registerAdminReq *RegisterAdminRequest) (*AuthResponse, *exceptions.AppError) {
	if err := ValidateRegisterAdminRequest(registerAdminReq); err != nil {
		return nil, err
	}

	if registerAdminReq.SystemAuthToken != as.systemAuthToken {
		return nil, exceptions.NewAppError("UNAUTHORIZED", "Unauthorized", nil)
	}

	createUser := &users.CreateUser{
		FirstName: registerAdminReq.FirstName,
		LastName:  registerAdminReq.LastName,
		Email:     registerAdminReq.Email,
		Password:  registerAdminReq.Password,
		Roles:     []string{"ADMIN"},
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
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SessionId: session.SessionId,
		Roles:     user.Roles,
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
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SessionId: session.SessionId,
		Roles:     user.Roles,
	}
	return resp, nil
}

func (as *AuthService) SetSessionRole(ctx context.Context, sessionId string, selectRoleReq *SelectRoleRequest) (*AuthResponse, *exceptions.AppError) {
	if err := ValidateSelectRoleRequest(selectRoleReq); err != nil {
		return nil, err
	}

	session, err := as.sessionSvc.GetActiveSessionById(ctx, nil, sessionId)
	if err != nil {
		return nil, err
	}

	user, err := as.userSvc.GetUserById(ctx, nil, session.UserId)
	if err != nil {
		return nil, err
	}

	hasRole := false
	for _, role := range user.Roles {
		if role == selectRoleReq.Role {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return nil, exceptions.NewAppError("UNAUTHORIZED", "User does not have the requested role", nil)
	}

	_, err = as.sessionSvc.SetSessionRole(ctx, nil, sessionId, selectRoleReq.Role)
	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SessionId: sessionId,
		Roles:     user.Roles,
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

	if session.ActiveRole != nil && *session.ActiveRole != "" {
		for _, role := range request.RequiredRoles {
			if role == *session.ActiveRole {
				return true, nil
			}
		}
		return false, exceptions.NewAppError("UNAUTHORIZED", "Unauthorized", nil)
	}

	rolesMap, err := as.userSvc.GetUsersRoles(ctx, nil, []int32{session.UserId})
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
