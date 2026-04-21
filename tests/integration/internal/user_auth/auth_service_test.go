package user_auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"github.com/todorpopov/school-manager/persistence"
	"github.com/todorpopov/school-manager/tests/integration"
)

type AuthServiceSuite struct {
	integration.TestSuite
	authSvc   user_auth.IAuthService
	bcryptSvc internal.IBCryptService
	txFactory persistence.ITransactionFactory
	usersSvc  users.IUserService
	rolesRepo roles.IRoleRepository
}

func (suite *AuthServiceSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.bcryptSvc = internal.NewBCryptService()
	suite.txFactory = persistence.NewTransactionFactory(suite.Db)
	usersRepo := users.NewUserRepository(suite.Db, suite.Logger)
	suite.usersSvc = users.NewUserService(suite.bcryptSvc, usersRepo, suite.txFactory)
	sessionsRepo := sessions.NewSessionRepository(suite.Db, suite.Env.SessionExpiration, suite.Logger)
	sessionsSvc := sessions.NewSessionService(sessionsRepo)
	suite.authSvc = user_auth.NewAuthService(suite.bcryptSvc, suite.usersSvc, sessionsSvc)
	suite.rolesRepo = roles.NewRoleRepository(suite.Db, suite.Logger)
}

func (suite *AuthServiceSuite) SetupTest() {
	suite.TestSuite.SetupTest()

	_, err := suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "ADMIN"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "USER"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "EMPLOYEE"})
	suite.Require().Nil(err)
}

func (suite *AuthServiceSuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}

func (suite *AuthServiceSuite) TestRegisterUser() {
	testCases := []struct {
		name            string
		registerRequest *user_auth.RegisterRequest
		expectError     bool
		errorCode       string
	}{
		{
			name: "Successfully register user with USER role",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"USER"},
			},
			expectError: false,
		},
		{
			name: "Successfully register user with ADMIN role",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "Admin",
				LastName:  "User",
				Email:     "admin.user@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"ADMIN"},
			},
			expectError: false,
		},
		{
			name: "Successfully register user with multiple roles",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "Multi",
				LastName:  "Role",
				Email:     "multi.role@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"USER", "EMPLOYEE"},
			},
			expectError: false,
		},
		{
			name: "Fail to register user with missing first name",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "",
				LastName:  "Doe",
				Email:     "missing.firstname@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"USER"},
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to register user with missing last name",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "John",
				LastName:  "",
				Email:     "missing.lastname@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"USER"},
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to register user with invalid email",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "SecurePass123!",
				Roles:     []string{"USER"},
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to register user with weak password",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "weak.password@example.com",
				Password:  "weak",
				Roles:     []string{"USER"},
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to register user with no roles",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "No",
				LastName:  "Roles",
				Email:     "no.roles@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{},
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to register user with invalid role",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "Invalid",
				LastName:  "Role",
				Email:     "invalid.role@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"INVALID_ROLE"},
			},
			expectError: true,
			errorCode:   "ROLES_DO_NOT_EXIST",
		},
		{
			name: "Fail to register user with duplicate email",
			registerRequest: &user_auth.RegisterRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "duplicate@example.com",
				Password:  "SecurePass123!",
				Roles:     []string{"USER"},
			},
			expectError: true,
			errorCode:   "UNIQUE_VIOLATION",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.name == "Fail to register user with duplicate email" {
				firstRequest := &user_auth.RegisterRequest{
					FirstName: "First",
					LastName:  "User",
					Email:     "duplicate@example.com",
					Password:  "SecurePass123!",
					Roles:     []string{"USER"},
				}
				_, err := suite.authSvc.RegisterUser(suite.Ctx, firstRequest)
				suite.Require().Nil(err, "Expected first registration to succeed")
			}

			authResp, appErr := suite.authSvc.RegisterUser(suite.Ctx, tc.registerRequest)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Equal(tc.errorCode, appErr.Code)
				suite.Require().Nil(authResp, "Expected response to be nil when error occurs")
			} else {
				suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
				suite.Require().NotNil(authResp, "Expected auth response to be returned")
				suite.Require().NotEmpty(authResp.SessionId, "Expected session ID to be generated")

				user, err := suite.usersSvc.GetUserByEmail(suite.Ctx, nil, tc.registerRequest.Email)
				suite.Require().Nil(err)
				suite.Require().NotNil(user)
				suite.Require().Equal(tc.registerRequest.FirstName, user.FirstName)
				suite.Require().Equal(tc.registerRequest.LastName, user.LastName)
				suite.Require().Equal(tc.registerRequest.Email, user.Email)
				suite.Require().ElementsMatch(tc.registerRequest.Roles, user.Roles, "Expected roles to match")
			}
		})
	}
}

func (suite *AuthServiceSuite) TestLogUserIn() {
	testCases := []struct {
		name         string
		loginRequest *user_auth.LoginRequest
		setupUser    bool
		expectError  bool
		errorCode    string
	}{
		{
			name: "Successfully log in user with correct credentials",
			loginRequest: &user_auth.LoginRequest{
				Email:    "login.success@example.com",
				Password: "CorrectPass123!",
			},
			setupUser:   true,
			expectError: false,
		},
		{
			name: "Fail to log in user with incorrect password",
			loginRequest: &user_auth.LoginRequest{
				Email:    "wrong.password@example.com",
				Password: "WrongPassword123!",
			},
			setupUser:   true,
			expectError: true,
			errorCode:   "INVALID_CREDENTIALS",
		},
		{
			name: "Fail to log in user that does not exist",
			loginRequest: &user_auth.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "AnyPass123!",
			},
			setupUser:   false,
			expectError: true,
			errorCode:   "NOT_FOUND",
		},
		{
			name: "Fail to log in with invalid email format",
			loginRequest: &user_auth.LoginRequest{
				Email:    "invalid-email",
				Password: "AnyPass123!",
			},
			setupUser:   false,
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to log in with empty password",
			loginRequest: &user_auth.LoginRequest{
				Email:    "empty.pass@example.com",
				Password: "",
			},
			setupUser:   false,
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.setupUser {
				registerReq := &user_auth.RegisterRequest{
					FirstName: "Test",
					LastName:  "User",
					Email:     tc.loginRequest.Email,
					Password:  "CorrectPass123!",
					Roles:     []string{"USER"},
				}
				_, err := suite.authSvc.RegisterUser(suite.Ctx, registerReq)
				suite.Require().Nil(err, "Expected user registration to succeed")
			}

			authResp, appErr := suite.authSvc.LogUserIn(suite.Ctx, tc.loginRequest)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Equal(tc.errorCode, appErr.Code)
				suite.Require().Nil(authResp, "Expected response to be nil when error occurs")
			} else {
				suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
				suite.Require().NotNil(authResp, "Expected auth response to be returned")
				suite.Require().NotEmpty(authResp.SessionId, "Expected session ID to be generated")
			}
		})
	}
}

func (suite *AuthServiceSuite) TestIsRequestAuthorized() {
	testCases := []struct {
		name        string
		authRequest *user_auth.AuthRequest
		setupUser   bool
		userRoles   []string
		expectAuth  bool
		expectError bool
		errorCode   string
	}{
		{
			name: "Successfully authorize user with ADMIN role",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{"ADMIN"},
			},
			setupUser:   true,
			userRoles:   []string{"ADMIN"},
			expectAuth:  true,
			expectError: false,
		},
		{
			name: "Successfully authorize user with USER role",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{"USER"},
			},
			setupUser:   true,
			userRoles:   []string{"USER"},
			expectAuth:  true,
			expectError: false,
		},
		{
			name: "Successfully authorize user with multiple roles - has one required",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{"ADMIN", "EMPLOYEE"},
			},
			setupUser:   true,
			userRoles:   []string{"EMPLOYEE"},
			expectAuth:  true,
			expectError: false,
		},
		{
			name: "Fail to authorize user without required role",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{"ADMIN"},
			},
			setupUser:   true,
			userRoles:   []string{"USER"},
			expectAuth:  false,
			expectError: true,
			errorCode:   "UNAUTHORIZED",
		},
		{
			name: "Fail to authorize with invalid session ID",
			authRequest: &user_auth.AuthRequest{
				SessionId:     "invalid-session-id",
				RequiredRoles: []string{"ADMIN"},
			},
			setupUser:   false,
			expectAuth:  false,
			expectError: true,
			errorCode:   "DATABASE_ERROR",
		},
		{
			name: "Fail to authorize with invalid role name",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{"INVALID_ROLE!"},
			},
			setupUser:   false,
			expectAuth:  false,
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Fail to authorize with empty required roles",
			authRequest: &user_auth.AuthRequest{
				RequiredRoles: []string{},
			},
			setupUser:   false,
			expectAuth:  false,
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var sessionId string

			if tc.setupUser {
				registerReq := &user_auth.RegisterRequest{
					FirstName: "Auth",
					LastName:  "Test",
					Email:     fmt.Sprintf("auth.test.%d@example.com", time.Now().UnixNano()),
					Password:  "SecurePass123!",
					Roles:     []string{"USER"},
				}
				authResp, err := suite.authSvc.RegisterUser(suite.Ctx, registerReq)
				suite.Require().Nil(err, "Expected user registration to succeed")
				sessionId = authResp.SessionId

				if len(tc.userRoles) > 0 && (len(tc.userRoles) != 1 || tc.userRoles[0] != "USER") {
					user, err := suite.usersSvc.GetUserByEmail(suite.Ctx, nil, registerReq.Email)
					suite.Require().Nil(err)

					updateUser := &users.UpdateUser{
						UserId:    user.UserId,
						FirstName: user.FirstName,
						LastName:  user.LastName,
						Email:     user.Email,
						Roles:     tc.userRoles,
					}
					_, err = suite.usersSvc.UpdateUser(suite.Ctx, nil, updateUser)
					suite.Require().Nil(err, "Expected role assignment to succeed")
				}

				tc.authRequest.SessionId = sessionId
			}

			isAuthorized, appErr := suite.authSvc.IsRequestAuthorized(suite.Ctx, tc.authRequest)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Equal(tc.errorCode, appErr.Code)
				suite.Require().False(isAuthorized, "Expected authorization to be false when error occurs")
			} else {
				suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
				suite.Require().Equal(tc.expectAuth, isAuthorized, "Authorization result mismatch")
			}
		})
	}
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
