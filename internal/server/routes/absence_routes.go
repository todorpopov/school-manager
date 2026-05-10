package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/absences"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterAbsenceRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, absenceSvc absences.IAbsenceService, authSvc user_auth.IAuthService) {
	logger.Info("Registering absence routes")

	logging := middleware.Logging(logger)
	requireAdminOrTeacher := middleware.RequireRoles(writer, authSvc, "ADMIN", "TEACHER")
	requireAdminTeacherParentOrStudent := middleware.RequireRoles(writer, authSvc, "ADMIN", "TEACHER", "PARENT", "STUDENT")

	s.Handle("POST /api/absences/bulk",
		middleware.Chain(
			handlers.BulkCreateAbsencesHandler(writer, absenceSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)

	s.Handle("POST /api/absence",
		middleware.Chain(
			handlers.CreateAbsenceHandler(writer, absenceSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)

	s.Handle("GET /api/absence/{absence_id}",
		middleware.Chain(
			handlers.GetAbsenceByIdHandler(writer, absenceSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("GET /api/absences",
		middleware.Chain(
			handlers.GetAbsencesHandler(writer, absenceSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("GET /api/student/{student_id}/absences",
		middleware.Chain(
			handlers.GetAbsencesByStudentIdHandler(writer, absenceSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("DELETE /api/absence/{absence_id}",
		middleware.Chain(
			handlers.DeleteAbsenceHandler(writer, absenceSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)

	s.Handle("PATCH /api/absence/{absence_id}/excuse",
		middleware.Chain(
			handlers.ExcuseAbsenceHandler(writer, absenceSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)
}

