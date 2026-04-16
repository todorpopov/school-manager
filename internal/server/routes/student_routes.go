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

func RegisterStudentRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, studentSvc students.IStudentService, authSvc user_auth.IAuthService) {
	logger.Info("Registering student routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/student",
		middleware.Chain(
			handlers.CreateStudentHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/student/{student_id}",
		middleware.Chain(
			handlers.GetStudentByIdHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/student/user/{user_id}",
		middleware.Chain(
			handlers.GetStudentByUserIdHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/students",
		middleware.Chain(
			handlers.GetStudentsHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("PUT /api/student/{student_id}",
		middleware.Chain(
			handlers.UpdateStudentHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/student/{student_id}",
		middleware.Chain(
			handlers.DeleteStudentHandler(writer, studentSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

