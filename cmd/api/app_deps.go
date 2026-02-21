package main

import (
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
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
