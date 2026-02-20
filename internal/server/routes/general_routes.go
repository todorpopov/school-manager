package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/server/handlers"
	"go.uber.org/zap"
)

func RegisterGeneralRoutes(s *http.ServeMux, logger *zap.Logger) {
	logger.Info("Registering general routes")
	s.Handle("GET /", handlers.PingHandler())
}
