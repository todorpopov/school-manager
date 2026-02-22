package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterAuthRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, authSvc user_auth.IAuthService) {
	logger.Info("Registering auth routes")
	s.Handle("POST /api/auth/login", handlers.LogUserIn(writer, authSvc, logger))
}
