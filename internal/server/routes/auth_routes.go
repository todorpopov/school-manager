package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterAuthRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, authSvc user_auth.IAuthService) {
	logger.Info("Registering auth routes")

	logging := middleware.Logging(logger)

	s.Handle("POST /api/auth/register",
		middleware.Chain(
			handlers.RegisterUserHandler(writer, authSvc, logger),
			logging,
		),
	)

	s.Handle("POST /api/auth/register-admin",
		middleware.Chain(
			handlers.RegisterAdminHandler(writer, authSvc, logger),
			logging,
		),
	)

	s.Handle("POST /api/auth/login",
		middleware.Chain(
			handlers.LogUserInHandler(writer, authSvc, logger),
			logging,
		),
	)
}
