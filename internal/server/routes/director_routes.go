package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/directors"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterDirectorRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, directorSvc directors.IDirectorService, authSvc user_auth.IAuthService) {
	logger.Info("Registering director routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/director",
		middleware.Chain(
			handlers.CreateDirectorHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/director/{director_id}",
		middleware.Chain(
			handlers.GetDirectorByIdHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/director/user/{user_id}",
		middleware.Chain(
			handlers.GetDirectorByUserIdHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/directors",
		middleware.Chain(
			handlers.GetDirectorsHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("PUT /api/director/{director_id}",
		middleware.Chain(
			handlers.UpdateDirectorHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/director/{director_id}",
		middleware.Chain(
			handlers.DeleteDirectorHandler(writer, directorSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}
