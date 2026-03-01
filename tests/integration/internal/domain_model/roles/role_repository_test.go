package roles

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/tests/integration"
)

type RoleRepositorySuite struct {
	integration.TestSuite
	rolesRepo roles.IRoleRepository
}

func (suite *RoleRepositorySuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	suite.rolesRepo = roles.NewRoleRepository(suite.Db, suite.Logger)
}

func (suite *RoleRepositorySuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *RoleRepositorySuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}

func (suite *RoleRepositorySuite) TestCreateRole() {
	testCases := []struct {
		name           string
		createRole     *roles.CreateRole
		useTransaction bool
		expectError    bool
	}{
		{
			name: "Successfully create role without transaction",
			createRole: &roles.CreateRole{
				RoleName: "Admin",
			},
			useTransaction: false,
			expectError:    false,
		},
		{
			name: "Successfully create role with transaction",
			createRole: &roles.CreateRole{
				RoleName: "Teacher",
			},
			useTransaction: true,
			expectError:    false,
		},
		{
			name: "Fail to create role with duplicate name without transaction",
			createRole: &roles.CreateRole{
				RoleName: "Duplicate",
			},
			useTransaction: false,
			expectError:    true,
		},
		{
			name: "Fail to create role with duplicate name with transaction",
			createRole: &roles.CreateRole{
				RoleName: "Duplicate2",
			},
			useTransaction: true,
			expectError:    true,
		},
		{
			name: "Fail to create role with empty name",
			createRole: &roles.CreateRole{
				RoleName: "",
			},
			useTransaction: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error

			if tc.expectError && tc.createRole.RoleName != "" {
				_, createErr := suite.rolesRepo.CreateRole(suite.Ctx, nil, tc.createRole)
				suite.Require().Nil(createErr, "Expected no error when creating initial duplicate role")
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			role, appErr := suite.rolesRepo.CreateRole(suite.Ctx, tx, tc.createRole)

			if tc.expectError {
				suite.Require().NotNil(appErr, "Expected an error but got none")
				suite.Require().Nil(role, "Expected role to be nil when error occurs")
			} else {
				suite.Require().Nil(appErr, "Expected no error but got: %v", appErr)
				suite.Require().NotNil(role, "Expected role to be returned")
				suite.Require().NotZero(role.RoleId, "Expected role ID to be generated")
				suite.Require().Equal(tc.createRole.RoleName, role.RoleName)

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				retrievedRole, getErr := suite.rolesRepo.GetRoleById(suite.Ctx, nil, role.RoleId)
				suite.Require().Nil(getErr, "Expected no error when retrieving role")
				suite.Require().Equal(role.RoleId, retrievedRole.RoleId)
				suite.Require().Equal(tc.createRole.RoleName, retrievedRole.RoleName)
			}
		})
	}
}

func (suite *RoleRepositorySuite) TestGetRoleById() {
	testCases := []struct {
		name                    string
		createRole              *roles.CreateRole
		useTransaction          bool
		shouldCreateInitialRole bool
		useNonExistentId        bool
		useInvalidId            bool
		expectError             bool
	}{
		{
			name: "Successfully get role without transaction",
			createRole: &roles.CreateRole{
				RoleName: "Student",
			},
			useTransaction:          false,
			shouldCreateInitialRole: true,
			useNonExistentId:        false,
			useInvalidId:            false,
			expectError:             false,
		},
		{
			name: "Successfully get role with transaction",
			createRole: &roles.CreateRole{
				RoleName: "Parent",
			},
			useTransaction:          true,
			shouldCreateInitialRole: true,
			useNonExistentId:        false,
			useInvalidId:            false,
			expectError:             false,
		},
		{
			name: "Fail to get role without transaction - role does not exist",
			createRole: &roles.CreateRole{
				RoleName: "NonExistent",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        true,
			useInvalidId:            false,
			expectError:             true,
		},
		{
			name: "Fail to get role with transaction - role does not exist",
			createRole: &roles.CreateRole{
				RoleName: "AnotherNonExistent",
			},
			useTransaction:          true,
			shouldCreateInitialRole: false,
			useNonExistentId:        true,
			useInvalidId:            false,
			expectError:             true,
		},
		{
			name: "Fail to get role with invalid ID (negative)",
			createRole: &roles.CreateRole{
				RoleName: "Invalid",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        false,
			useInvalidId:            true,
			expectError:             true,
		},
		{
			name: "Fail to get role with invalid ID (zero)",
			createRole: &roles.CreateRole{
				RoleName: "Invalid2",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        false,
			useInvalidId:            true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var savedRole *roles.Role
			var roleIdToQuery int32

			if tc.shouldCreateInitialRole {
				savedRole, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, tc.createRole)
				suite.Require().Nil(err, "Expected no error when creating initial role")
				suite.Require().NotNil(savedRole, "Expected initial role to not be nil")
				roleIdToQuery = savedRole.RoleId
			} else if tc.useNonExistentId {
				roleIdToQuery = 99999
			} else if tc.useInvalidId {
				if tc.createRole.RoleName == "Invalid" {
					roleIdToQuery = -1
				} else {
					roleIdToQuery = 0
				}
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			foundRole, getErr := suite.rolesRepo.GetRoleById(suite.Ctx, tx, roleIdToQuery)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(foundRole, "Expected role to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)
				suite.Require().NotNil(foundRole, "Expected role to be returned")
				suite.Require().NotZero(foundRole.RoleId, "Expected role ID to be generated")
				suite.Require().Equal(savedRole.RoleId, foundRole.RoleId)
				suite.Require().Equal(savedRole.RoleName, foundRole.RoleName)

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
			}
		})
	}
}

func (suite *RoleRepositorySuite) TestGetRoles() {
	testCases := []struct {
		name              string
		rolesToCreate     []*roles.CreateRole
		useTransaction    bool
		expectEmptyResult bool
		expectError       bool
	}{
		{
			name: "Successfully get all roles without transaction",
			rolesToCreate: []*roles.CreateRole{
				{RoleName: "Admin"},
				{RoleName: "Teacher"},
				{RoleName: "Student"},
			},
			useTransaction:    false,
			expectEmptyResult: false,
			expectError:       false,
		},
		{
			name: "Successfully get all roles with transaction",
			rolesToCreate: []*roles.CreateRole{
				{RoleName: "Manager"},
				{RoleName: "Staff"},
			},
			useTransaction:    true,
			expectEmptyResult: false,
			expectError:       false,
		},
		{
			name:              "Successfully get empty list when no roles exist without transaction",
			rolesToCreate:     []*roles.CreateRole{},
			useTransaction:    false,
			expectEmptyResult: true,
			expectError:       false,
		},
		{
			name:              "Successfully get empty list when no roles exist with transaction",
			rolesToCreate:     []*roles.CreateRole{},
			useTransaction:    true,
			expectEmptyResult: true,
			expectError:       false,
		},
		{
			name: "Successfully get single role without transaction",
			rolesToCreate: []*roles.CreateRole{
				{RoleName: "SingleRole"},
			},
			useTransaction:    false,
			expectEmptyResult: false,
			expectError:       false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var createdRoles []*roles.Role

			for _, createRole := range tc.rolesToCreate {
				role, createErr := suite.rolesRepo.CreateRole(suite.Ctx, nil, createRole)
				suite.Require().Nil(createErr, "Expected no error when creating initial role")
				suite.Require().NotNil(role, "Expected role to not be nil")
				createdRoles = append(createdRoles, role)
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			foundRoles, getErr := suite.rolesRepo.GetRoles(suite.Ctx, tx)

			if tc.expectError {
				suite.Require().NotNil(getErr, "Expected an error but got none")
				suite.Require().Nil(foundRoles, "Expected roles to be nil when error occurs")
			} else {
				suite.Require().Nil(getErr, "Expected no error but got: %v", getErr)

				if tc.expectEmptyResult {
					suite.Require().Empty(foundRoles, "Expected roles list to be empty")
				} else {
					suite.Require().NotNil(foundRoles, "Expected roles to be returned")
					suite.Require().Len(foundRoles, len(tc.rolesToCreate), "Expected to get all created roles")

					for _, createdRole := range createdRoles {
						found := false
						for _, foundRole := range foundRoles {
							if foundRole.RoleId == createdRole.RoleId {
								found = true
								suite.Require().Equal(createdRole.RoleName, foundRole.RoleName)
								break
							}
						}
						suite.Require().True(found, "Expected to find created role in result")
					}
				}

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}
			}

			suite.CleanDatabase()
		})
	}
}

func (suite *RoleRepositorySuite) TestDeleteRole() {
	testCases := []struct {
		name                    string
		createRole              *roles.CreateRole
		useTransaction          bool
		shouldCreateInitialRole bool
		useNonExistentId        bool
		useInvalidId            bool
		expectError             bool
	}{
		{
			name: "Successfully delete role without transaction",
			createRole: &roles.CreateRole{
				RoleName: "ToDelete1",
			},
			useTransaction:          false,
			shouldCreateInitialRole: true,
			useNonExistentId:        false,
			useInvalidId:            false,
			expectError:             false,
		},
		{
			name: "Successfully delete role with transaction",
			createRole: &roles.CreateRole{
				RoleName: "ToDelete2",
			},
			useTransaction:          true,
			shouldCreateInitialRole: true,
			useNonExistentId:        false,
			useInvalidId:            false,
			expectError:             false,
		},
		{
			name: "Fail to delete role without transaction - role does not exist",
			createRole: &roles.CreateRole{
				RoleName: "NonExistent",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        true,
			useInvalidId:            false,
			expectError:             true,
		},
		{
			name: "Fail to delete role with transaction - role does not exist",
			createRole: &roles.CreateRole{
				RoleName: "AnotherNonExistent",
			},
			useTransaction:          true,
			shouldCreateInitialRole: false,
			useNonExistentId:        true,
			useInvalidId:            false,
			expectError:             true,
		},
		{
			name: "Fail to delete role with invalid ID (negative)",
			createRole: &roles.CreateRole{
				RoleName: "Invalid",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        false,
			useInvalidId:            true,
			expectError:             true,
		},
		{
			name: "Fail to delete role with invalid ID (zero)",
			createRole: &roles.CreateRole{
				RoleName: "Invalid2",
			},
			useTransaction:          false,
			shouldCreateInitialRole: false,
			useNonExistentId:        false,
			useInvalidId:            true,
			expectError:             true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var tx pgx.Tx
			var err error
			var savedRole *roles.Role
			var roleIdToDelete int32

			if tc.shouldCreateInitialRole {
				savedRole, err = suite.rolesRepo.CreateRole(suite.Ctx, nil, tc.createRole)
				suite.Require().Nil(err, "Expected no error when creating initial role")
				suite.Require().NotNil(savedRole, "Expected initial role to not be nil")
				roleIdToDelete = savedRole.RoleId
			} else if tc.useNonExistentId {
				roleIdToDelete = 99999
			} else if tc.useInvalidId {
				if tc.createRole.RoleName == "Invalid" {
					roleIdToDelete = -1
				} else {
					roleIdToDelete = 0
				}
			}

			if tc.useTransaction {
				tx, err = suite.Db.Pool.Begin(suite.Ctx)
				suite.Require().NoError(err)
				defer func() {
					_ = tx.Rollback(suite.Ctx)
				}()
			}

			deleteErr := suite.rolesRepo.DeleteRole(suite.Ctx, tx, roleIdToDelete)

			if tc.expectError {
				suite.Require().NotNil(deleteErr, "Expected an error but got none")
			} else {
				suite.Require().Nil(deleteErr, "Expected no error but got: %v", deleteErr)

				if tc.useTransaction {
					commitErr := tx.Commit(suite.Ctx)
					suite.Require().NoError(commitErr)
				}

				foundRole, getErr := suite.rolesRepo.GetRoleById(suite.Ctx, nil, roleIdToDelete)
				suite.Require().NotNil(getErr, "Expected an error when trying to get deleted role")
				suite.Require().Nil(foundRole, "Expected role to be nil after deletion")
			}
		})
	}
}

func TestRoleRepositorySuite(t *testing.T) {
	suite.Run(t, new(RoleRepositorySuite))
}
