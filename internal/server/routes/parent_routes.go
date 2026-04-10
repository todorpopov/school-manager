package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterParentRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, parentSvc parents.IParentService, authSvc user_auth.IAuthService) {
	logger.Info("Registering parent routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/parent",
		middleware.Chain(
			handlers.CreateParentHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/parent/{parent_id}",
		middleware.Chain(
			handlers.GetParentByIdHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/parent/user/{user_id}",
		middleware.Chain(
			handlers.GetParentByUserIdHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/parents",
		middleware.Chain(
			handlers.GetParentsHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("PUT /api/parent/{parent_id}",
		middleware.Chain(
			handlers.UpdateParentHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/parent/{parent_id}",
		middleware.Chain(
			handlers.DeleteParentHandler(writer, parentSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

