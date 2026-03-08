package integration

import (
	"context"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	return zap.NewExample()
}

type TestSuite struct {
	suite.Suite
	Env         *configs.Config
	Logger      *zap.Logger
	DbContainer testcontainers.Container
	Db          *persistence.Database
	Ctx         context.Context
}

func (suite *TestSuite) SetupSuite() {
	suite.Ctx = context.Background()
	logger := newLogger()

	suite.Env = configs.NewTestConfig()
	suite.Logger = logger

	var err error
	suite.Db, err = persistence.InitDatabase(suite.Env, suite.Logger)
	suite.Require().NoError(err)
}

func (suite *TestSuite) TearDownSuite() {
	err := persistence.ShutdownDatabase(suite.Env, suite.Db, suite.Logger)
	suite.Require().NoError(err)
}

func (suite *TestSuite) CleanDatabase() {
	_, err := suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM sessions;")
	suite.Require().NoError(err)

	_, err = suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM user_roles;")
	suite.Require().NoError(err)

	_, err = suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM users;")
	suite.Require().NoError(err)

	_, err = suite.Db.Pool.Exec(suite.Ctx, "DELETE FROM roles;")
	suite.Require().NoError(err)
}

func (suite *TestSuite) SetupTest() {
	suite.CleanDatabase()
}
