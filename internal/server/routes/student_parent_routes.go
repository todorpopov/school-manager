package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterStudentParentRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, studentParentSvc students.IStudentParentService, authSvc user_auth.IAuthService) {
	logger.Info("Registering student-parent relationship routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/student-parent/student/{student_id}/parent/{parent_id}",
		middleware.Chain(
			handlers.LinkParentToStudentHandler(writer, studentParentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/student-parent/student/{student_id}/parent/{parent_id}",
		middleware.Chain(
			handlers.UnlinkParentFromStudentHandler(writer, studentParentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/student-parent/student/{student_id}/parents",
		middleware.Chain(
			handlers.GetParentsForStudentHandler(writer, studentParentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/student-parent/parent/{parent_id}/students",
		middleware.Chain(
			handlers.GetStudentsForParentHandler(writer, studentParentSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}
