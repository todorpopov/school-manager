package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func RegisterUserRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, usrSvc users.IUserService) {
	logger.Info("Registering user routes")
	s.Handle("POST /api/user", handlers.CreateUserHandler(writer, usrSvc, logger))
}
