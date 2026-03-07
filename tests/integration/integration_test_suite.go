package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

	pgReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	dbContainer, err := testcontainers.GenericContainer(suite.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: pgReq,
		Started:          true,
	})
	suite.Require().NoError(err)
	suite.DbContainer = dbContainer

	dbHost, err := dbContainer.Host(suite.Ctx)
	suite.Require().NoError(err)
	dbPort, err := dbContainer.MappedPort(suite.Ctx, "5432")
	suite.Require().NoError(err)

	dbUrl := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", dbHost, dbPort.Port())

	suite.Env = configs.NewTestConfig(dbUrl)
	suite.Logger = logger

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
