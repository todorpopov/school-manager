package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterClassRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, classSvc classes.IClassService, authSvc user_auth.IAuthService) {
	logger.Info("Registering class routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/class",
		middleware.Chain(
			handlers.CreateClassHandler(writer, classSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/class/{class_id}",
		middleware.Chain(
			handlers.GetClassByIdHandler(writer, classSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/classes",
		middleware.Chain(
			handlers.GetClassesHandler(writer, classSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/class/{class_id}",
		middleware.Chain(
			handlers.DeleteClassHandler(writer, classSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}
