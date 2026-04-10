package server

import (
	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/directors"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type Dependencies struct {
	UserRepo     users.IUserRepository
	UserSvc      users.IUserService
	SessionRepo  sessions.ISessionRepository
	SessionSvc   sessions.ISessionService
	RoleRepo     roles.IRoleRepository
	RoleSvc      roles.IRoleService
	DirectorRepo directors.IDirectorRepository
	DirectorSvc  directors.IDirectorService
	AuthSvc      user_auth.IAuthService
}

func NewDependencies(env *configs.Config, db *persistence.Database, logger *zap.Logger) *Dependencies {
	bcryptSvc := internal.NewBCryptService()
	txFactory := persistence.NewTransactionFactory(db)

	usrRepo := users.NewUserRepository(db, logger)
	usrSvc := users.NewUserService(bcryptSvc, usrRepo, txFactory)

	sessionRepo := sessions.NewSessionRepository(db, env.SessionExpiration, logger)
	sessionSvc := sessions.NewSessionService(sessionRepo)

	roleRepo := roles.NewRoleRepository(db, logger)
	roleSvc := roles.NewRoleService(roleRepo)

	directorRepo := directors.NewDirectorRepository(db, logger)
	directorSvc := directors.NewDirectorService(directorRepo, usrSvc, txFactory)

	authSvc := user_auth.NewAuthService(bcryptSvc, usrSvc, sessionSvc)
	return &Dependencies{
		UserRepo:     usrRepo,
		UserSvc:      usrSvc,
		SessionRepo:  sessionRepo,
		SessionSvc:   sessionSvc,
		RoleRepo:     roleRepo,
		RoleSvc:      roleSvc,
		DirectorRepo: directorRepo,
		DirectorSvc:  directorSvc,
		AuthSvc:      authSvc,
	}
}
