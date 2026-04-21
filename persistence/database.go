package persistence

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/todorpopov/school-manager/configs"
	"go.uber.org/zap"
)

type PersistentStore interface {
	HealthCheck(ctx context.Context) error
	Close()
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
}

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) HealthCheck(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *Database) Close() {
	db.Pool.Close()
}

func NewDatabase(env *configs.Config) (*Database, error) {
	conf, err := pgxpool.ParseConfig(env.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	conf.MinConns = env.DBMinConnections
	conf.MaxConns = env.DBMaxConnections
	conf.MaxConnLifetime = env.DBMaxConnectionLifetime
	conf.MaxConnIdleTime = env.DBMaxConnectionIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

type MigrationDirection int

const (
	MigrationUp MigrationDirection = iota
	MigrationDown
)

func (md MigrationDirection) IsValid() bool {
	return md == MigrationUp || md == MigrationDown
}

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func Migrate(env *configs.Config, direction MigrationDirection) error {
	if !direction.IsValid() {
		return fmt.Errorf("invalid migration direction: %v", direction)
	}

	dsn := env.DBUrl
	goose.SetBaseFS(embeddedMigrations)

	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}
	sqlDB := stdlib.OpenDB(*conf)
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to close database: %v", err))
		}
	}(sqlDB)

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	switch direction {
	case MigrationDown:
		err = goose.Reset(sqlDB, "migrations")
	case MigrationUp:
		err = goose.Up(sqlDB, "migrations")
	}

	if err != nil {
		return err
	}

	return nil
}

func InitDatabase(env *configs.Config, logger *zap.Logger) (*Database, error) {
	err := Migrate(env, MigrationUp)
	if err != nil {
		logger.Error("Failed to apply database migrations", zap.Error(err))
		return nil, err
	}

	logger.Info("Database migrations applied")
	db, err := NewDatabase(env)
	if err != nil {
		logger.Error("Failed to establish database connection", zap.Error(err))
		return nil, err
	}
	logger.Info("Database connection established")

	err = db.HealthCheck(context.Background())
	if err != nil {
		logger.Error("Failed to perform database health check", zap.Error(err))
		return nil, err
	}
	logger.Info("Database health check passed")

	return db, nil
}

func ShutdownDatabase(env *configs.Config, db *Database, logger *zap.Logger) error {
	db.Close()
	err := Migrate(env, MigrationDown)
	if err != nil {
		logger.Error("Failed to rollback database migrations", zap.Error(err))
		return err
	}
	logger.Info("Database migrations rolled back")

	return nil
}
