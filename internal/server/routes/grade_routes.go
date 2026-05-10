package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/grades"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterGradeRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, gradeSvc grades.IGradeService, authSvc user_auth.IAuthService) {
	logger.Info("Registering grade routes")

	logging := middleware.Logging(logger)
	requireAdminOrTeacher := middleware.RequireRoles(writer, authSvc, "ADMIN", "TEACHER")
	requireAdminTeacherParentOrStudent := middleware.RequireRoles(writer, authSvc, "ADMIN", "TEACHER", "PARENT", "STUDENT")

	s.Handle("POST /api/grades/bulk",
		middleware.Chain(
			handlers.BulkCreateGradesHandler(writer, gradeSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)

	s.Handle("POST /api/grade",
		middleware.Chain(
			handlers.CreateGradeHandler(writer, gradeSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)

	s.Handle("GET /api/grade/{grade_id}",
		middleware.Chain(
			handlers.GetGradeByIdHandler(writer, gradeSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("GET /api/grades",
		middleware.Chain(
			handlers.GetGradesHandler(writer, gradeSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("GET /api/student/{student_id}/grades",
		middleware.Chain(
			handlers.GetGradesByStudentIdHandler(writer, gradeSvc, logger),
			logging,
			requireAdminTeacherParentOrStudent,
		),
	)

	s.Handle("DELETE /api/grade/{grade_id}",
		middleware.Chain(
			handlers.DeleteGradeHandler(writer, gradeSvc, logger),
			logging,
			requireAdminOrTeacher,
		),
	)
}

