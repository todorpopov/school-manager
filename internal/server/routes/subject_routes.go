package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterSubjectRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, subjectSvc subjects.ISubjectService, authSvc user_auth.IAuthService) {
	logger.Info("Registering subject routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/subject",
		middleware.Chain(
			handlers.CreateSubjectHandler(writer, subjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/subject/{subject_id}",
		middleware.Chain(
			handlers.GetSubjectByIdHandler(writer, subjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/subjects",
		middleware.Chain(
			handlers.GetSubjectsHandler(writer, subjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/subject/{subject_id}",
		middleware.Chain(
			handlers.DeleteSubjectHandler(writer, subjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

