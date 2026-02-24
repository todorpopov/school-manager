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
	s.Handle("GET /api/user/{user_id}", handlers.GetUserByIdHandler(writer, usrSvc, logger))
	s.Handle("GET /api/user/email/{email}", handlers.GetUserByEmailHandler(writer, usrSvc, logger))
	s.Handle("GET /api/users", handlers.GetUsersHandler(writer, usrSvc, logger))
	s.Handle("PUT /api/user/{user_id}", handlers.UpdateUserHandler(writer, usrSvc, logger))
	s.Handle("PUT /api/user/password/{user_id}", handlers.UpdateUserPasswordHandler(writer, usrSvc, logger))
	s.Handle("DELETE /api/user/{user_id}", handlers.DeleteUserHandler(writer, usrSvc, logger))
}
