package main

import (
	"time"

	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppDeps struct {
	UserRepo users.IUserRepository
	UserSvc  users.IUserService
}

func NewAppDeps(db *persistence.Database, logger *zap.Logger) *AppDeps {
	usrRepo := users.NewUserRepository(db, logger)
	usrSvc := users.NewUserService(usrRepo, logger)

	return &AppDeps{
		UserRepo: usrRepo,
		UserSvc:  usrSvc,
	}
}

func NewLogger() (*zap.Logger, error) {
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = timeEncoder
	return cfg.Build()
}
