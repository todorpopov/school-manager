package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterTeacherRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, teacherSvc teachers.ITeacherService, authSvc user_auth.IAuthService) {
	logger.Info("Registering teacher routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/teacher",
		middleware.Chain(
			handlers.CreateTeacherHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/teacher/{teacher_id}",
		middleware.Chain(
			handlers.GetTeacherByIdHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/teacher/user/{user_id}",
		middleware.Chain(
			handlers.GetTeacherByUserIdHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/teachers",
		middleware.Chain(
			handlers.GetTeachersHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("PUT /api/teacher/{teacher_id}",
		middleware.Chain(
			handlers.UpdateTeacherHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/teacher/{teacher_id}",
		middleware.Chain(
			handlers.DeleteTeacherHandler(writer, teacherSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

