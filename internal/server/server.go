package server

import (
	"context"
	"net/http"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/routes"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	RegisterRoutes()
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	env        *configs.Config
	mux        *http.ServeMux
	writer     *writer.HttpWriter
	logger     *zap.Logger
	server     *http.Server
	serverDeps *Dependencies
}

func NewHttpServer(env *configs.Config, logger *zap.Logger, deps *Dependencies) *HttpServer {
	errWriter := exceptions.NewErrorWriter()
	return &HttpServer{
		env,
		http.NewServeMux(),
		writer.NewHttpWriter(errWriter),
		logger,
		nil,
		deps,
	}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Session-Id")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *HttpServer) Start() error {
	s.registerRoutes()
	s.logger.Info("Starting HTTP server on port: " + s.env.ApiPort)
	srvAddr := ":" + s.env.ApiPort
	s.server = &http.Server{
		Addr:    srvAddr,
		Handler: enableCors(s.mux),
	}
	return s.server.ListenAndServe()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) registerRoutes() {
	routes.RegisterGeneralRoutes(s.mux, s.writer, s.logger)
	routes.RegisterUserRoutes(s.mux, s.writer, s.logger, s.serverDeps.UserSvc, s.serverDeps.AuthSvc)
	routes.RegisterAuthRoutes(s.mux, s.writer, s.logger, s.serverDeps.AuthSvc)
	routes.RegisterRoleRoutes(s.mux, s.writer, s.logger, s.serverDeps.RoleRepo, s.serverDeps.AuthSvc)
	routes.RegisterSchoolRoutes(s.mux, s.writer, s.logger, s.serverDeps.SchoolSvc, s.serverDeps.AuthSvc)
	routes.RegisterDirectorRoutes(s.mux, s.writer, s.logger, s.serverDeps.DirectorSvc, s.serverDeps.AuthSvc)
	routes.RegisterTeacherRoutes(s.mux, s.writer, s.logger, s.serverDeps.TeacherSvc, s.serverDeps.AuthSvc)
	routes.RegisterTeacherSubjectRoutes(s.mux, s.writer, s.logger, s.serverDeps.TeacherSubjectSvc, s.serverDeps.AuthSvc)
	routes.RegisterParentRoutes(s.mux, s.writer, s.logger, s.serverDeps.ParentSvc, s.serverDeps.AuthSvc)
	routes.RegisterClassRoutes(s.mux, s.writer, s.logger, s.serverDeps.ClassSvc, s.serverDeps.AuthSvc)
	routes.RegisterStudentRoutes(s.mux, s.writer, s.logger, s.serverDeps.StudentSvc, s.serverDeps.AuthSvc)
	routes.RegisterStudentParentRoutes(s.mux, s.writer, s.logger, s.serverDeps.StudentParentSvc, s.serverDeps.AuthSvc)
	routes.RegisterTermRoutes(s.mux, s.writer, s.logger, s.serverDeps.TermSvc, s.serverDeps.AuthSvc)
	routes.RegisterSubjectRoutes(s.mux, s.writer, s.logger, s.serverDeps.SubjectSvc, s.serverDeps.AuthSvc)
	routes.RegisterCurriculumRoutes(s.mux, s.writer, s.logger, s.serverDeps.CurriculumSvc, s.serverDeps.AuthSvc)
	routes.RegisterGradeRoutes(s.mux, s.writer, s.logger, s.serverDeps.GradeSvc, s.serverDeps.AuthSvc)
	routes.RegisterAbsenceRoutes(s.mux, s.writer, s.logger, s.serverDeps.AbsenceSvc, s.serverDeps.AuthSvc)
	routes.RegisterReportingRoutes(s.mux, s.writer, s.logger, s.serverDeps.DynamicRptSvc, s.serverDeps.AuthSvc)
}
