package users

import (
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/tests/integration"
)

type UserServiceSuite struct {
	integration.TestSuite
	usersSvc users.IUserService
}

func (suite *UserServiceSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()
	bcryptSvc := internal.NewBCryptService()
	usersRepo := users.NewUserRepository(suite.Db, suite.Logger)
	suite.usersSvc = users.NewUserService(bcryptSvc, usersRepo)
}

func (suite *UserServiceSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *UserServiceSuite) TearDownSuite() {
	suite.TestSuite.TearDownSuite()
}
