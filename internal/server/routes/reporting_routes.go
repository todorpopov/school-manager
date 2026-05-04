package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/reporting"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterReportingRoutes(
	s *http.ServeMux,
	writer *writer.HttpWriter,
	logger *zap.Logger,
	rptSvc reporting.IDynamicReportingService,
	authSvc user_auth.IAuthService,
) {
	logger.Info("Registering reporting routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/reporting",
		middleware.Chain(
			handlers.ReportingQueryHandler(writer, rptSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}
