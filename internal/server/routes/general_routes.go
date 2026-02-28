package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func RegisterGeneralRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger) {
	logger.Info("Registering general routes")

	s.Handle("GET /",
		middleware.Chain(
			handlers.PingHandler(writer),
			middleware.Logging(logger),
		),
	)
}
