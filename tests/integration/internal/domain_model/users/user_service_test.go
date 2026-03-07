package users

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/tests/integration"
)

type UserServiceSuite struct {
	integration.TestSuite
	bcryptSvc internal.IBCryptService
	usersSvc  users.IUserService
	rolesRepo roles.IRoleRepository
}

func (suite *UserServiceSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.bcryptSvc = internal.NewBCryptService()
	usersRepo := users.NewUserRepository(suite.Db, suite.Logger)
	suite.usersSvc = users.NewUserService(suite.bcryptSvc, usersRepo)
	suite.rolesRepo = roles.NewRoleRepository(suite.Db, suite.Logger)
}

func (suite *UserServiceSuite) SetupTest() {
	suite.TestSuite.SetupTest()

	_, err := suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "ADMIN"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "TEACHER"})
	suite.Require().Nil(err)
	_, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, &roles.CreateRole{RoleName: "STUDENT"})
	suite.Require().Nil(err)
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
				Roles:     []string{"ADMIN", "TEACHER"},
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
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"TEACHER"},
			},
			useTransaction: true,
			expectError:    true,
		},
		{
			name: "Fail to create user with invalid role",
			createUser: &users.CreateUser{
				FirstName: "Invalid",
				LastName:  "Role",
				Email:     "invalid.role@example.com",
				Password:  "hashedPassword202",
				Roles:     []string{"NONEXISTENT_ROLE"},
			},
			useTransaction: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error

			if tc.expectError && tc.name != "Fail to create user with invalid role" {
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
				suite.Require().ElementsMatch(tc.createUser.Roles, user.Roles, "Expected roles to match")

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
				suite.Require().ElementsMatch(tc.createUser.Roles, retrievedUser.Roles, "Expected retrieved user roles to match")
			}
		})
	}
}

func (suite *UserServiceSuite) TestGetUserById() {
	testCases := []struct {
		name                    string
		createUser              *users.CreateUser
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentId        bool
		expectError             bool
	}{
		{
			name: "Successfully get user without transaction",
			createUser: &users.CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashedPassword123",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentId:        false,
			expectError:             false,
		},
		{
			name: "Successfully get user with transaction",
			createUser: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Password:  "hashedPassword456",
				Roles:     []string{"TEACHER"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentId:        false,
			expectError:             false,
		},
		{
			name: "Fail to get user without transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "User",
				Email:     "nonexistent@example.com",
				Password:  "hashedPassword789",
				Roles:     []string{"STUDENT"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentId:        true,
			expectError:             true,
		},
		{
			name: "Fail to get user with transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "Another",
				LastName:  "NonExistent",
				Email:     "another.nonexistent@example.com",
				Password:  "hashedPassword101",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: false,
			useNonExistentId:        true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var savedUsr *users.User
			var userIdToQuery int32

			if tc.shouldCreateInitialUser {
				savedUsr, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.createUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(savedUsr, "Expected initial user to not be nil")
				userIdToQuery = savedUsr.UserId
			} else if tc.useNonExistentId {
				userIdToQuery = 99999
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			foundUsr, getErr := suite.usersSvc.GetUserById(suite.Ctx, tx, userIdToQuery)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(foundUsr, "Expected user to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)
				suite.Require().NotNil(foundUsr, "Expected user to be returned")
				suite.Require().NotZero(foundUsr.UserId, "Expected user ID to be generated")
				suite.Require().Equal(savedUsr.UserId, foundUsr.UserId)
				suite.Require().Equal(savedUsr.FirstName, foundUsr.FirstName)
				suite.Require().Equal(savedUsr.LastName, foundUsr.LastName)
				suite.Require().Equal(savedUsr.Email, foundUsr.Email)
				suite.Require().Nil(foundUsr.Password, "Expected password to not be returned")
				suite.Require().Equal(savedUsr.Password, foundUsr.Password)
				suite.Require().ElementsMatch(savedUsr.Roles, foundUsr.Roles, "Expected roles to match")

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
			}
		})
	}
}

func (suite *UserServiceSuite) TestGetUserByEmail() {
	testCases := []struct {
		name                    string
		createUser              *users.CreateUser
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentEmail     bool
		expectError             bool
	}{
		{
			name: "Successfully get user by email without transaction",
			createUser: &users.CreateUser{
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Password:  "hashedPassword111",
				Roles:     []string{"TEACHER"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentEmail:     false,
			expectError:             false,
		},
		{
			name: "Successfully get user by email with transaction",
			createUser: &users.CreateUser{
				FirstName: "Bob",
				LastName:  "Williams",
				Email:     "bob.williams@example.com",
				Password:  "hashedPassword222",
				Roles:     []string{"STUDENT"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentEmail:     false,
			expectError:             false,
		},
		{
			name: "Fail to get user by email without transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "User",
				Email:     "nonexistent.email@example.com",
				Password:  "hashedPassword333",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentEmail:     true,
			expectError:             true,
		},
		{
			name: "Fail to get user by email with transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "Another",
				LastName:  "NonExistent",
				Email:     "another.nonexistent.email@example.com",
				Password:  "hashedPassword444",
				Roles:     []string{"TEACHER"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: false,
			useNonExistentEmail:     true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var savedUsr *users.User
			var emailToQuery string

			if tc.shouldCreateInitialUser {
				savedUsr, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.createUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(savedUsr, "Expected initial user to not be nil")
				emailToQuery = savedUsr.Email
			} else if tc.useNonExistentEmail {
				emailToQuery = "doesnotexist@example.com"
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			foundUsr, getErr := suite.usersSvc.GetUserByEmail(suite.Ctx, tx, emailToQuery)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(foundUsr, "Expected user to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)
				suite.Require().NotNil(foundUsr, "Expected user to be returned")
				suite.Require().NotZero(foundUsr.UserId, "Expected user ID to be generated")
				suite.Require().Equal(savedUsr.UserId, foundUsr.UserId)
				suite.Require().Equal(savedUsr.FirstName, foundUsr.FirstName)
				suite.Require().Equal(savedUsr.LastName, foundUsr.LastName)
				suite.Require().Equal(savedUsr.Email, foundUsr.Email)
				suite.Require().Nil(foundUsr.Password, "Expected password to not be returned")
				suite.Require().Equal(savedUsr.Password, foundUsr.Password)
				suite.Require().ElementsMatch(savedUsr.Roles, foundUsr.Roles, "Expected roles to match")

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
			}
		})
	}
}

func (suite *UserServiceSuite) TestGetUserByEmailWithPass() {
	testCases := []struct {
		name                    string
		createUser              *users.CreateUser
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentEmail     bool
		expectError             bool
	}{
		{
			name: "Successfully get user by email with pass without transaction",
			createUser: &users.CreateUser{
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Password:  "hashedPassword111",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentEmail:     false,
			expectError:             false,
		},
		{
			name: "Successfully get user by email with pass with transaction",
			createUser: &users.CreateUser{
				FirstName: "Bob",
				LastName:  "Williams",
				Email:     "bob.williams@example.com",
				Password:  "hashedPassword222",
				Roles:     []string{"TEACHER"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentEmail:     false,
			expectError:             false,
		},
		{
			name: "Fail to get user by email with pass without transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "User",
				Email:     "nonexistent.email@example.com",
				Password:  "hashedPassword333",
				Roles:     []string{"STUDENT"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentEmail:     true,
			expectError:             true,
		},
		{
			name: "Fail to get user by email with pass with transaction - user does not exist",
			createUser: &users.CreateUser{
				FirstName: "Another",
				LastName:  "NonExistent",
				Email:     "another.nonexistent.email@example.com",
				Password:  "hashedPassword444",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: false,
			useNonExistentEmail:     true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var savedUsr *users.User
			var emailToQuery string

			if tc.shouldCreateInitialUser {
				savedUsr, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.createUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(savedUsr, "Expected initial user to not be nil")
				emailToQuery = savedUsr.Email
			} else if tc.useNonExistentEmail {
				emailToQuery = "doesnotexist@example.com"
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			foundUsr, getErr := suite.usersSvc.GetUserByEmailWithPass(suite.Ctx, tx, emailToQuery)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(foundUsr, "Expected user to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)
				suite.Require().NotNil(foundUsr, "Expected user to be returned")
				suite.Require().NotZero(foundUsr.UserId, "Expected user ID to be generated")
				suite.Require().Equal(savedUsr.UserId, foundUsr.UserId)
				suite.Require().Equal(savedUsr.FirstName, foundUsr.FirstName)
				suite.Require().Equal(savedUsr.LastName, foundUsr.LastName)
				suite.Require().Equal(savedUsr.Email, foundUsr.Email)
				suite.Require().NotNil(foundUsr.Password, "Expected password to be returned")
				match := suite.bcryptSvc.PasswordsMatch(*foundUsr.Password, tc.createUser.Password)
				suite.Require().True(match, "Password should match the hash")
				suite.Require().ElementsMatch(savedUsr.Roles, foundUsr.Roles, "Expected roles to match")

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
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
				Roles:     []string{"STUDENT"},
			},
			updateUser: &users.UpdateUser{
				FirstName: "Charles",
				LastName:  "Brownson",
				Email:     "charles.brownson@example.com",
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"TEACHER"},
			},
			updateUser: &users.UpdateUser{
				FirstName: "Wonder",
				LastName:  "Woman",
				Email:     "wonder.woman@example.com",
				Roles:     []string{"TEACHER"},
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
				Roles:     []string{"ADMIN"},
			},
			updateUser: &users.UpdateUser{
				UserId:    99999,
				FirstName: "Updated",
				LastName:  "Ghost",
				Email:     "updated.ghost@example.com",
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"TEACHER"},
			},
			updateUser: &users.UpdateUser{
				UserId:    99999,
				FirstName: "Updated",
				LastName:  "Phantom",
				Email:     "updated.phantom@example.com",
				Roles:     []string{"ADMIN"},
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
				Roles:     []string{"STUDENT"},
			},
			updateUser: &users.UpdateUser{
				FirstName: "First",
				LastName:  "User",
				Email:     "duplicate.target@example.com",
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"TEACHER"},
			},
			updateUser: &users.UpdateUser{
				FirstName: "Second",
				LastName:  "User",
				Email:     "duplicate.target2@example.com",
				Roles:     []string{"TEACHER"},
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
					Roles:     []string{"STUDENT"},
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
				suite.Require().ElementsMatch(tc.updateUser.Roles, updatedUser.Roles, "Expected roles to match after update")

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
				suite.Require().ElementsMatch(tc.updateUser.Roles, retrievedUser.Roles, "Expected roles to match in retrieved user")
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
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"TEACHER"},
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
				Roles:     []string{"ADMIN"},
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
				Roles:     []string{"STUDENT"},
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

func (suite *UserServiceSuite) TestDeleteUser() {
	testCases := []struct {
		name                    string
		initialUser             *users.CreateUser
		useTransaction          bool
		shouldCreateInitialUser bool
		useNonExistentUserId    bool
		expectError             bool
	}{
		{
			name: "Successfully delete user without transaction",
			initialUser: &users.CreateUser{
				FirstName: "Delete",
				LastName:  "Test",
				Email:     "delete.test@example.com",
				Password:  "hashedPasswordDelete123",
				Roles:     []string{"STUDENT"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			expectError:             false,
		},
		{
			name: "Successfully delete user with transaction",
			initialUser: &users.CreateUser{
				FirstName: "Delete",
				LastName:  "TestTwo",
				Email:     "delete.test2@example.com",
				Password:  "hashedPasswordDelete456",
				Roles:     []string{"TEACHER"},
			},
			useTransaction:          true,
			shouldCreateInitialUser: true,
			useNonExistentUserId:    false,
			expectError:             false,
		},
		{
			name: "Fail to delete user without transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "Delete",
				Email:     "nonexistent.delete@example.com",
				Password:  "somePassword",
				Roles:     []string{"ADMIN"},
			},
			useTransaction:          false,
			shouldCreateInitialUser: false,
			useNonExistentUserId:    true,
			expectError:             true,
		},
		{
			name: "Fail to delete user with transaction - user does not exist",
			initialUser: &users.CreateUser{
				FirstName: "NonExistent",
				LastName:  "DeleteTwo",
				Email:     "nonexistent.delete2@example.com",
				Password:  "somePassword2",
				Roles:     []string{"STUDENT"},
			},
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
			var userIdToDelete int32

			if tc.shouldCreateInitialUser {
				createdUser, err = suite.usersSvc.CreateUser(suite.Ctx, nil, tc.initialUser)
				suite.Require().Nil(err, "Expected no error when creating initial user")
				suite.Require().NotNil(createdUser, "Expected initial user to not be nil")
				userIdToDelete = createdUser.UserId
			} else if tc.useNonExistentUserId {
				userIdToDelete = 99999
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			deleteErr := suite.usersSvc.DeleteUser(suite.Ctx, tx, userIdToDelete)

			if tc.expectError {
				suite.Require().NotNil(deleteErr, "Expected an error but got none")
			} else {
				suite.Require().Nil(deleteErr, "Expected no error but got: %v", deleteErr)

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				retrievedUser, getErr := suite.usersSvc.GetUserById(suite.Ctx, nil, userIdToDelete)
				suite.Require().NotNil(getErr, "Expected an error when retrieving deleted user")
				suite.Require().Nil(retrievedUser, "Expected user to be nil after deletion")
			}
		})
	}
}

func (suite *UserServiceSuite) TestGetUsers() {
	testCases := []struct {
		name           string
		usersToCreate  []*users.CreateUser
		useTransaction bool
		expectError    bool
	}{
		{
			name:           "Successfully get all users without transaction - empty database",
			usersToCreate:  []*users.CreateUser{},
			useTransaction: false,
			expectError:    false,
		},
		{
			name: "Successfully get all users without transaction - single user",
			usersToCreate: []*users.CreateUser{
				{
					FirstName: "Single",
					LastName:  "User",
					Email:     "single.user@example.com",
					Password:  "hashedPassword001",
					Roles:     []string{"STUDENT"},
				},
			},
			useTransaction: false,
			expectError:    false,
		},
		{
			name: "Successfully get all users without transaction - multiple users",
			usersToCreate: []*users.CreateUser{
				{
					FirstName: "First",
					LastName:  "User",
					Email:     "first.user@example.com",
					Password:  "hashedPassword002",
					Roles:     []string{"ADMIN", "TEACHER"},
				},
				{
					FirstName: "Second",
					LastName:  "User",
					Email:     "second.user@example.com",
					Password:  "hashedPassword003",
					Roles:     []string{"STUDENT"},
				},
				{
					FirstName: "Third",
					LastName:  "User",
					Email:     "third.user@example.com",
					Password:  "hashedPassword004",
					Roles:     []string{"TEACHER"},
				},
			},
			useTransaction: false,
			expectError:    false,
		},
		{
			name:           "Successfully get all users with transaction - empty database",
			usersToCreate:  []*users.CreateUser{},
			useTransaction: true,
			expectError:    false,
		},
		{
			name: "Successfully get all users with transaction - single user",
			usersToCreate: []*users.CreateUser{
				{
					FirstName: "Single",
					LastName:  "UserTx",
					Email:     "single.usertx@example.com",
					Password:  "hashedPassword005",
					Roles:     []string{"ADMIN"},
				},
			},
			useTransaction: true,
			expectError:    false,
		},
		{
			name: "Successfully get all users with transaction - multiple users",
			usersToCreate: []*users.CreateUser{
				{
					FirstName: "Alpha",
					LastName:  "User",
					Email:     "alpha.user@example.com",
					Password:  "hashedPassword006",
					Roles:     []string{"STUDENT", "TEACHER"},
				},
				{
					FirstName: "Beta",
					LastName:  "User",
					Email:     "beta.user@example.com",
					Password:  "hashedPassword007",
					Roles:     []string{"ADMIN"},
				},
				{
					FirstName: "Gamma",
					LastName:  "User",
					Email:     "gamma.user@example.com",
					Password:  "hashedPassword008",
					Roles:     []string{"TEACHER"},
				},
				{
					FirstName: "Delta",
					LastName:  "User",
					Email:     "delta.user@example.com",
					Password:  "hashedPassword009",
					Roles:     []string{"STUDENT"},
				},
			},
			useTransaction: true,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var createdUsers []*users.User

			for _, userToCreate := range tc.usersToCreate {
				createdUser, createErr := suite.usersSvc.CreateUser(suite.Ctx, nil, userToCreate)
				suite.Require().Nil(createErr, "Expected no error when creating user")
				suite.Require().NotNil(createdUser, "Expected user to be created")
				createdUsers = append(createdUsers, createdUser)
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			retrievedUsers, getErr := suite.usersSvc.GetUsers(suite.Ctx, tx)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(retrievedUsers, "Expected users to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)

				if len(createdUsers) == 0 {
					suite.Require().True(retrievedUsers == nil || len(retrievedUsers) == 0,
						"Expected users to be nil or empty slice for empty database")
				} else {
					suite.Require().NotNil(retrievedUsers, "Expected users slice to be returned")
					suite.Require().Equal(len(createdUsers), len(retrievedUsers), "Expected number of retrieved users to match created users")

					for _, createdUser := range createdUsers {
						found := false
						for _, retrievedUser := range retrievedUsers {
							if retrievedUser.UserId == createdUser.UserId {
								found = true
								suite.Require().Equal(createdUser.FirstName, retrievedUser.FirstName)
								suite.Require().Equal(createdUser.LastName, retrievedUser.LastName)
								suite.Require().Equal(createdUser.Email, retrievedUser.Email)
								suite.Require().Nil(retrievedUser.Password, "Expected password to not be returned")
								suite.Require().Equal(createdUser.Password, retrievedUser.Password)
								suite.Require().ElementsMatch(createdUser.Roles, retrievedUser.Roles, "Expected roles to match")
								break
							}
						}
						suite.Require().True(found, "Expected created user with ID %d to be in retrieved users", createdUser.UserId)
					}
				}

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
			}

			_, err = suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM user_roles;")
			suite.Require().NoError(err)
			_, err = suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM users;")
			suite.Require().NoError(err)
		})
	}
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
