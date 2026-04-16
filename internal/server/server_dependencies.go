package server

import (
	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/directors"
	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
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
	TeacherRepo  teachers.ITeacherRepository
	TeacherSvc   teachers.ITeacherService
	ParentRepo   parents.IParentRepository
	ParentSvc    parents.IParentService
	ClassRepo    classes.IClassRepository
	ClassSvc     classes.IClassService
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

	teacherRepo := teachers.NewTeacherRepository(db, logger)
	teacherSvc := teachers.NewTeacherService(teacherRepo, usrSvc, txFactory)

	parentRepo := parents.NewParentRepository(db, logger)
	parentSvc := parents.NewParentService(parentRepo, usrSvc, txFactory)

	classRepo := classes.NewClassRepository(db, logger)
	classSvc := classes.NewClassService(classRepo)

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
		TeacherRepo:  teacherRepo,
		TeacherSvc:   teacherSvc,
		ParentRepo:   parentRepo,
		ParentSvc:    parentSvc,
		ClassRepo:    classRepo,
		ClassSvc:     classSvc,
		AuthSvc:      authSvc,
	}
}
