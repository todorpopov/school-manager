package main

import (
	"time"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppDeps struct {
	UserRepo    users.IUserRepository
	UserSvc     users.IUserService
	SessionRepo sessions.ISessionRepository
	SessionSvc  sessions.ISessionService
	RoleRepo    roles.IRoleRepository
	RoleSvc     roles.IRoleService
	AuthSvc     user_auth.IAuthService
}

func NewAppDeps(env *configs.Config, db *persistence.Database, logger *zap.Logger) *AppDeps {
	bcryptSvc := internal.NewBCryptService()

	usrRepo := users.NewUserRepository(db, logger)
	usrSvc := users.NewUserService(bcryptSvc, usrRepo)

	sessionRepo := sessions.NewSessionRepository(db, env.SessionExpiration, logger)
	sessionSvc := sessions.NewSessionService(sessionRepo)

	roleRepo := roles.NewRoleRepository(db, logger)
	roleSvc := roles.NewRoleService(roleRepo)

	authSvc := user_auth.NewAuthService(bcryptSvc, usrSvc, sessionSvc)
	return &AppDeps{
		UserRepo:    usrRepo,
		UserSvc:     usrSvc,
		SessionRepo: sessionRepo,
		SessionSvc:  sessionSvc,
		RoleRepo:    roleRepo,
		RoleSvc:     roleSvc,
		AuthSvc:     authSvc,
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
