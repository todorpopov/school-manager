package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func RegisterUserRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, usrSvc users.IUserService) {
	logger.Info("Registering user routes")

	logging := middleware.Logging(logger)

	s.Handle("POST /api/user",
		middleware.Chain(
			handlers.CreateUserHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("GET /api/user/{user_id}",
		middleware.Chain(
			handlers.GetUserByIdHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("GET /api/user/email/{email}",
		middleware.Chain(
			handlers.GetUserByEmailHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("GET /api/users",
		middleware.Chain(
			handlers.GetUsersHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("PUT /api/user/{user_id}",
		middleware.Chain(
			handlers.UpdateUserHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("PUT /api/user/password/{user_id}",
		middleware.Chain(
			handlers.UpdateUserPasswordHandler(writer, usrSvc, logger),
			logging,
		),
	)

	s.Handle("DELETE /api/user/{user_id}",
		middleware.Chain(
			handlers.DeleteUserHandler(writer, usrSvc, logger),
			logging,
		),
	)
}
