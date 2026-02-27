package users

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/tests/integration"
)

type UserServiceSuite struct {
	integration.TestSuite
	bcryptSvc internal.IBCryptService
	usersSvc  users.IUserService
}

func (suite *UserServiceSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.bcryptSvc = internal.NewBCryptService()
	usersRepo := users.NewUserRepository(suite.Db, suite.Logger)
	suite.usersSvc = users.NewUserService(suite.bcryptSvc, usersRepo)
}

func (suite *UserServiceSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *UserServiceSuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (suite *UserServiceSuite) TestCreateUser() {
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
				_, createErr := suite.usersSvc.CreateUser(suite.Ctx, nil, tc.createUser)
				suite.Require().Nil(createErr, "Expected no error when creating initial duplicate user")
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			user, appErr := suite.usersSvc.CreateUser(suite.Ctx, tx, tc.createUser)

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
				suite.Require().Nil(user.Password, "Expected password to be nil in returned user object")

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				retrievedUser, getErr := suite.usersSvc.GetUserByEmailWithPass(suite.Ctx, nil, tc.createUser.Email)
				suite.Require().Nil(getErr, "Expected no error when retrieving user")
				suite.Require().Equal(user.UserId, retrievedUser.UserId)
				suite.Require().Equal(tc.createUser.FirstName, retrievedUser.FirstName)
				suite.Require().Equal(tc.createUser.LastName, retrievedUser.LastName)
				suite.Require().Equal(tc.createUser.Email, retrievedUser.Email)
				suite.Require().NotNil(retrievedUser.Password, "Expected password to be returned in GetUserByEmailWithPass")
				suite.Require().NotEqual(tc.createUser.Password, *retrievedUser.Password)
			}
		})
	}
}

func (suite *UserServiceSuite) TestUpdateUser() {
	testCases := []struct {
		name                    string
		initialUser             *users.CreateUser
		updateUser              *users.UpdateUser
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentUserId    bool
		useDuplicateEmail       bool
		expectError             bool
	}{
		{
			name: "Successfully update user without transaction",
			initialUser: &users.CreateUser{
				FirstName: "Charlie",
				LastName:  "Brown",
				Email:     "charlie.brown@example.com",
				Password:  "hashedPassword555",
			},
			updateUser: &users.UpdateUser{
				FirstName: "Charles",
				LastName:  "Brownson",
				Email:     "charles.brownson@example.com",
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			useDuplicateEmail:       false,
			expectError:             false,
		},
		{
			name: "Successfully update user with transaction",
			initialUser: &users.CreateUser{
				FirstName: "Diana",
				LastName:  "Prince",
				Email:     "diana.prince@example.com",
				Password:  "hashedPassword666",
			},
			updateUser: &users.UpdateUser{
				FirstName: "Wonder",
				LastName:  "Woman",
				Email:     "wonder.woman@example.com",
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			useDuplicateEmail:       false,
			expectError:             false,
		},
		{
			name: "Fail to update user without transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "Ghost",
				LastName:  "User",
				Email:     "ghost.user@example.com",
				Password:  "hashedPassword777",
			},
			updateUser: &users.UpdateUser{
				UserId:    99999,
				FirstName: "Updated",
				LastName:  "Ghost",
				Email:     "updated.ghost@example.com",
			},
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentUserId:    true,
			useDuplicateEmail:       false,
			expectError:             true,
		},
		{
			name: "Fail to update user with transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "Phantom",
				LastName:  "User",
				Email:     "phantom.user@example.com",
				Password:  "hashedPassword888",
			},
			updateUser: &users.UpdateUser{
				UserId:    99999,
				FirstName: "Updated",
				LastName:  "Phantom",
				Email:     "updated.phantom@example.com",
			},
			useTransaction:          true,
			shouldCreateInitialUser: false,
			useNonExistentUserId:    true,
			useDuplicateEmail:       false,
			expectError:             true,
		},
		{
			name: "Fail to update user without transaction - duplicate email",
			initialUser: &users.CreateUser{
				FirstName: "First",
				LastName:  "User",
				Email:     "first.user@example.com",
				Password:  "hashedPassword999",
			},
			updateUser: &users.UpdateUser{
				FirstName: "First",
				LastName:  "User",
				Email:     "duplicate.target@example.com",
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			useDuplicateEmail:       true,
			expectError:             true,
		},
		{
			name: "Fail to update user with transaction - duplicate email",
			initialUser: &users.CreateUser{
				FirstName: "Second",
				LastName:  "User",
				Email:     "second.user@example.com",
				Password:  "hashedPassword000",
			},
			updateUser: &users.UpdateUser{
				FirstName: "Second",
				LastName:  "User",
				Email:     "duplicate.target2@example.com",
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			useDuplicateEmail:       true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var createdUser *users.User
			var duplicateUser *users.User

			if tc.shouldCreateInitialUser {
				createdUser, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.initialUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(createdUser, "Expected initial user to not be nil")
				tc.updateUser.UserId = createdUser.UserId
			}

			if tc.useDuplicateEmail {
				duplicateUser, err = suite.usersSvc.CreateUser(suite.Ctx, nil, &users.CreateUser{
					FirstName: "Duplicate",
					LastName:  "Target",
					Email:     tc.updateUser.Email,
					Password:  "hashedPasswordDuplicate",
				})
				suite.Require().Nil(err, "Expected no error when creating duplicate target user")
				suite.Require().NotNil(duplicateUser, "Expected duplicate target user to not be nil")
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			updatedUser, updateErr := suite.usersSvc.UpdateUser(suite.Ctx, tx, tc.updateUser)

			if tc.expectError {
				suite.Require().NotNil(updateErr, "Expected an error but got none")
				suite.Require().Nil(updatedUser, "Expected updated user to be nil when error occurs")
			} else {
				suite.Require().Nil(updateErr, "Expected no error but got: %v", updateErr)
				suite.Require().NotNil(updatedUser, "Expected updated user to be returned")
				suite.Require().Equal(tc.updateUser.UserId, updatedUser.UserId)
				suite.Require().Equal(tc.updateUser.FirstName, updatedUser.FirstName)
				suite.Require().Equal(tc.updateUser.LastName, updatedUser.LastName)
				suite.Require().Equal(tc.updateUser.Email, updatedUser.Email)
				suite.Require().Nil(updatedUser.Password, "Expected password to be returned")

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				retrievedUser, getErr := suite.usersSvc.GetUserByEmailWithPass(suite.Ctx, nil, updatedUser.Email)
				suite.Require().Nil(getErr, "Expected no error when retrieving updated user")
				suite.Require().Equal(updatedUser.UserId, retrievedUser.UserId)
				suite.Require().Equal(tc.updateUser.FirstName, retrievedUser.FirstName)
				suite.Require().Equal(tc.updateUser.LastName, retrievedUser.LastName)
				suite.Require().Equal(tc.updateUser.Email, retrievedUser.Email)
				suite.Require().NotNil(retrievedUser.Password, "Expected password to be returned in GetUserByEmailWithPass")
				suite.Require().NotEqual(updatedUser.Password, *retrievedUser.Password)
			}
		})
	}
}

func (suite *UserServiceSuite) TestUpdateUserPassword() {
	testCases := []struct {
		name                    string
		initialUser             *users.CreateUser
		newPassword             string
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentUserId    bool
		expectError             bool
	}{
		{
			name: "Successfully update user password without transaction",
			initialUser: &users.CreateUser{
				FirstName: "Password",
				LastName:  "Test",
				Email:     "password.test@example.com",
				Password:  "oldHashedPassword123",
			},
			newPassword:             "newHashedPassword456",
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			expectError:             false,
		},
		{
			name: "Successfully update user password with transaction",
			initialUser: &users.CreateUser{
				FirstName: "Password",
				LastName:  "TestTwo",
				Email:     "password.test2@example.com",
				Password:  "oldHashedPassword789",
			},
			newPassword:             "newHashedPassword012",
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			expectError:             false,
		},
		{
			name: "Fail to update password without transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "User",
				Email:     "nonexistent.password@example.com",
				Password:  "somePassword",
			},
			newPassword:             "newPassword",
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentUserId:    true,
			expectError:             true,
		},
		{
			name: "Fail to update password with transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "UserTwo",
				Email:     "nonexistent.password2@example.com",
				Password:  "somePassword2",
			},
			newPassword:             "newPassword2",
			useTransaction:          true,
			shouldCreateInitialUser: false,
			useNonExistentUserId:    true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var createdUser *users.User
			var userIdToUpdate int32

			if tc.shouldCreateInitialUser {
				createdUser, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.initialUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(createdUser, "Expected initial user to not be nil")
				userIdToUpdate = createdUser.UserId
			} else if tc.useNonExistentUserId {
				userIdToUpdate = 99999
			}

			updateUserPass := &users.UpdateUserPassword{
				UserId:   userIdToUpdate,
				Password: tc.newPassword,
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			updateErr := suite.usersSvc.UpdateUserPassword(suite.Ctx, tx, updateUserPass)

			if tc.expectError {
				suite.Require().NotNil(updateErr, "Expected an error but got none")
			} else {
				suite.Require().Nil(updateErr, "Expected no error but got: %v", updateErr)

				if tc.shouldCreateInitialUser {
					if tc.useTransaction {
						commitErr := tx.Commit(suite.Ctx)
						suite.Require().NoError(commitErr)
					}

					retrievedUser, getErr := suite.usersSvc.GetUserByEmailWithPass(suite.Ctx, nil, tc.initialUser.Email)
					suite.Require().Nil(getErr, "Expected no error when retrieving user")
					suite.Require().NotNil(retrievedUser, "Expected user to be returned")
					suite.Require().NotNil(retrievedUser.Password, "Expected password to be returned")

					match := suite.bcryptSvc.PasswordsMatch(*retrievedUser.Password, tc.newPassword)
					suite.Require().True(match, "Password should match the new value")

					suite.Require().Equal(createdUser.UserId, retrievedUser.UserId)
					suite.Require().Equal(createdUser.FirstName, retrievedUser.FirstName)
					suite.Require().Equal(createdUser.LastName, retrievedUser.LastName)
					suite.Require().Equal(createdUser.Email, retrievedUser.Email)
				}
			}
		})
	}
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
