package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterTermRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, termSvc terms.ITermService, authSvc user_auth.IAuthService) {
	logger.Info("Registering term routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/term",
		middleware.Chain(
			handlers.CreateTermHandler(writer, termSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/term/{term_id}",
		middleware.Chain(
			handlers.GetTermByIdHandler(writer, termSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/terms",
		middleware.Chain(
			handlers.GetTermsHandler(writer, termSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/term/{term_id}",
		middleware.Chain(
			handlers.DeleteTermHandler(writer, termSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

