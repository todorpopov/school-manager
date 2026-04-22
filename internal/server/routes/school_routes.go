package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterSchoolRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, schoolSvc schools.ISchoolService, authSvc user_auth.IAuthService) {
	logger.Info("Registering school routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/school",
		middleware.Chain(
			handlers.CreateSchoolHandler(writer, schoolSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/school/{school_id}",
		middleware.Chain(
			handlers.GetSchoolByIdHandler(writer, schoolSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/schools",
		middleware.Chain(
			handlers.GetSchoolsHandler(writer, schoolSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("PUT /api/school/{school_id}",
		middleware.Chain(
			handlers.UpdateSchoolHandler(writer, schoolSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/school/{school_id}",
		middleware.Chain(
			handlers.DeleteSchoolHandler(writer, schoolSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

