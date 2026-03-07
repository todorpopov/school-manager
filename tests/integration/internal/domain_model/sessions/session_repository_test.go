package sessions

import (
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/tests/integration"
)

type SessionRepositorySuite struct {
	integration.TestSuite
	sessionRepo sessions.ISessionRepository
	userRepo    users.IUserRepository
	rolesRepo   roles.IRoleRepository
}

func (suite *SessionRepositorySuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.sessionRepo = sessions.NewSessionRepository(suite.Db, suite.Env.SessionExpiration, suite.Logger)
	suite.userRepo = users.NewUserRepository(suite.Db, suite.Logger)
	suite.rolesRepo = roles.NewRoleRepository(suite.Db, suite.Logger)
}

func (suite *SessionRepositorySuite) SetupTest() {
	suite.TestSuite.SetupTest()

	_, err := suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "ADMIN"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "TEACHER"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "STUDENT"})
	suite.Require().Nil(err)
}

func (suite *SessionRepositorySuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}

func (suite *SessionRepositorySuite) TestCreateOrRenewSession() {
	testCases := []struct {
		name           string
		useTransaction bool
		seedSession    bool
		useInvalidUser bool
		expectError    bool
	}{
		{
			name:           "Create session without transaction",
			useTransaction: false,
			seedSession:    false,
			useInvalidUser: false,
			expectError:    false,
		},
		{
			name:           "Create session with transaction",
			useTransaction: true,
			seedSession:    false,
			useInvalidUser: false,
			expectError:    false,
		},
		{
			name:           "Renew session without transaction",
			useTransaction: false,
			seedSession:    true,
			useInvalidUser: false,
			expectError:    false,
		},
		{
			name:           "Renew session with transaction",
			useTransaction: true,
			seedSession:    true,
			useInvalidUser: false,
			expectError:    false,
		},
		{
			name:           "Fail to create session for missing user",
			useTransaction: false,
			seedSession:    false,
			useInvalidUser: true,
			expectError:    true,
		},
	}

	for idx, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var userId int32
			var seededSession *sessions.Session

			if tc.useInvalidUser {
				userId = 99999
			} else {
				createUser := &users.CreateUser{
					FirstName: "Session",
					LastName:  "User",
					Email:     fmt.Sprintf("session.%d@example.com", idx),
					Password:  "hashedPassword123",
					Roles:     []string{"STUDENT"},
				}
				savedUser, createErr := suite.userRepo.CreateUser(suite.Ctx, nil, createUser)
				suite.Require().Nil(createErr, "Expected no error when creating user")
				suite.Require().NotNil(savedUser, "Expected user to be returned")
				userId = savedUser.UserId
			}

			if tc.seedSession {
				seededSession, err = suite.sessionRepo.CreateOrRenewSession(suite.Ctx, nil, userId)
				suite.Require().Nil(err, "Expected no error when seeding session")
				suite.Require().NotNil(seededSession, "Expected seeded session")
				time.Sleep(10 * time.Millisecond)
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			startedAt := time.Now()
			session, appErr := suite.sessionRepo.CreateOrRenewSession(suite.Ctx, tx, userId)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Nil(session, "Expected session to be nil when error occurs")
				return
			}

			suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
			suite.Require().NotNil(session, "Expected session to be returned")
			suite.Require().NotEmpty(session.SessionId, "Expected session ID to be generated")
			suite.Require().Equal(userId, session.UserId)

			tolerance := 2 * time.Second
			minExpiry := startedAt.Add(suite.Env.SessionExpiration - tolerance)
			maxExpiry := startedAt.Add(suite.Env.SessionExpiration + tolerance)
			suite.Require().False(session.ExpiresAt.Before(minExpiry), "Expected expiry to be after min bound")
			suite.Require().False(session.ExpiresAt.After(maxExpiry), "Expected expiry to be before max bound")

			if tc.seedSession {
				suite.Require().Equal(seededSession.SessionId, session.SessionId)
				suite.Require().False(session.ExpiresAt.Before(seededSession.ExpiresAt), "Expected expiry to be renewed")
			}

			if tc.useTransaction {
				commitErr := tx.Commit(suite.Ctx)
				suite.Require().NoError(commitErr)
			}
		})
	}
}

func TestSessionRepositorySuite(t *testing.T) {
	suite.Run(t, new(SessionRepositorySuite))
}
