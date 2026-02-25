package users

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/tests/integration"
)

type UserRepositorySuite struct {
	integration.TestSuite
	usersRepo users.IUserRepository
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.usersRepo = users.NewUserRepository(suite.Db, suite.Logger)
}

func (suite *UserRepositorySuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *UserRepositorySuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}

func (suite *UserRepositorySuite) TestCreateUser() {
	testCases := []struct {
		name           string
		createUser     *users.CreateUser
		useTransaction bool
		expectError    bool
	}{
		{
			name: "Successfully create user without transaction",
			createUser: &users.CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedPassword123",
			},
			useTransaction: false,
			expectError:    false,
		},
		{
			name: "Successfully create user with transaction",
			createUser: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Password:  "hashedPassword456",
			},
			useTransaction: true,
			expectError:    false,
		},
		{
			name: "Fail to create user with duplicate email without transaction",
			createUser: &users.CreateUser{
				FirstName: "Duplicate",
				LastName:  "User",
				Email:     "duplicate@example.com",
				Password:  "hashedPassword789",
			},
			useTransaction: false,
			expectError:    true,
		},
		{
			name: "Fail to create user with duplicate email with transaction",
			createUser: &users.CreateUser{
				FirstName: "Another",
				LastName:  "Duplicate",
				Email:     "duplicate2@example.com",
				Password:  "hashedPassword101",
			},
			useTransaction: true,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error

			if tc.expectError {
				_, createErr := suite.usersRepo.CreateUser(suite.Ctx, nil, tc.createUser)
				suite.Require().Nil(createErr, "Expected no error when creating initial duplicate user")
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			user, appErr := suite.usersRepo.CreateUser(suite.Ctx, tx, tc.createUser)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Nil(user, "Expected user to be nil when error occurs")
			} else {
				suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
				suite.Require().NotNil(user, "Expected user to be returned")
				suite.Require().NotZero(user.UserId, "Expected user ID to be generated")
				suite.Require().Equal(tc.createUser.FirstName, user.FirstName)
				suite.Require().Equal(tc.createUser.LastName, user.LastName)
				suite.Require().Equal(tc.createUser.Email, user.Email)
				suite.Require().NotNil(user.Password, "Expected password to be returned")
				suite.Require().Equal(tc.createUser.Password, *user.Password)

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				retrievedUser, getErr := suite.usersRepo.GetUserByEmail(suite.Ctx, nil, tc.createUser.Email)
				suite.Require().Nil(getErr, "Expected no error when retrieving user")
				suite.Require().Equal(user.UserId, retrievedUser.UserId)
				suite.Require().Equal(tc.createUser.FirstName, retrievedUser.FirstName)
				suite.Require().Equal(tc.createUser.LastName, retrievedUser.LastName)
				suite.Require().Equal(tc.createUser.Email, retrievedUser.Email)
			}
		})
	}
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
