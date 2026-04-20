package server

import (
	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/absences"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/domain_model/directors"
	"github.com/todorpopov/school-manager/internal/domain_model/grades"
	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/domain_model/sessions"
	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type Dependencies struct {
	UserRepo       users.IUserRepository
	UserSvc        users.IUserService
	SessionRepo    sessions.ISessionRepository
	SessionSvc     sessions.ISessionService
	RoleRepo       roles.IRoleRepository
	RoleSvc        roles.IRoleService
	DirectorRepo   directors.IDirectorRepository
	DirectorSvc    directors.IDirectorService
	TeacherRepo    teachers.ITeacherRepository
	TeacherSvc     teachers.ITeacherService
	ParentRepo     parents.IParentRepository
	ParentSvc      parents.IParentService
	ClassRepo      classes.IClassRepository
	ClassSvc       classes.IClassService
	StudentRepo    students.IStudentRepository
	StudentSvc     students.IStudentService
	TermRepo       terms.ITermRepository
	TermSvc        terms.ITermService
	SubjectRepo    subjects.ISubjectRepository
	SubjectSvc     subjects.ISubjectService
	CurriculumRepo curricula.ICurriculumRepository
	CurriculumSvc  curricula.ICurriculumService
	GradeRepo      grades.IGradeRepository
	GradeSvc       grades.IGradeService
	AbsenceRepo    absences.IAbsenceRepository
	AbsenceSvc     absences.IAbsenceService
	AuthSvc        user_auth.IAuthService
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

	studentRepo := students.NewStudentRepository(db, logger)
	studentSvc := students.NewStudentService(studentRepo, usrSvc, txFactory)

	termRepo := terms.NewTermRepository(db, logger)
	termSvc := terms.NewTermService(termRepo)

	subjectRepo := subjects.NewSubjectRepository(db, logger)
	subjectSvc := subjects.NewSubjectService(subjectRepo)

	curriculumRepo := curricula.NewCurriculumRepository(db, logger)
	curriculumSvc := curricula.NewCurriculumService(curriculumRepo)

	gradeRepo := grades.NewGradeRepository(db, logger)
	gradeSvc := grades.NewGradeService(gradeRepo)

	absenceRepo := absences.NewAbsenceRepository(db, logger)
	absenceSvc := absences.NewAbsenceService(absenceRepo)

	authSvc := user_auth.NewAuthService(bcryptSvc, usrSvc, sessionSvc)
	return &Dependencies{
		UserRepo:       usrRepo,
		UserSvc:        usrSvc,
		SessionRepo:    sessionRepo,
		SessionSvc:     sessionSvc,
		RoleRepo:       roleRepo,
		RoleSvc:        roleSvc,
		DirectorRepo:   directorRepo,
		DirectorSvc:    directorSvc,
		TeacherRepo:    teacherRepo,
		TeacherSvc:     teacherSvc,
		ParentRepo:     parentRepo,
		ParentSvc:      parentSvc,
		ClassRepo:      classRepo,
		ClassSvc:       classSvc,
		StudentRepo:    studentRepo,
		StudentSvc:     studentSvc,
		TermRepo:       termRepo,
		TermSvc:        termSvc,
		SubjectRepo:    subjectRepo,
		SubjectSvc:     subjectSvc,
		CurriculumRepo: curriculumRepo,
		CurriculumSvc:  curriculumSvc,
		GradeRepo:      gradeRepo,
		GradeSvc:       gradeSvc,
		AbsenceRepo:    absenceRepo,
		AbsenceSvc:     absenceSvc,
		AuthSvc:        authSvc,
	}
}
