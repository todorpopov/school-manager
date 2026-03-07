package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterRoleRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, roleSvc roles.IRoleService, authSvc user_auth.IAuthService) {
	logger.Info("Registering role routes")

	logging := middleware.Logging(logger)
	requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/role",
		middleware.Chain(
			handlers.CreateRoleHandler(writer, roleSvc, logger),
			logging,
			requireAdmin,
		),
	)

	s.Handle("GET /api/role/{role_id}",
		middleware.Chain(
			handlers.GetRoleByIdHandler(writer, roleSvc, logger),
			logging,
			requireAdmin,
		),
	)

	s.Handle("GET /api/roles",
		middleware.Chain(
			handlers.GetRolesHandler(writer, roleSvc, logger),
			logging,
			requireAdmin,
		),
	)

	s.Handle("DELETE /api/role/{role_id}",
		middleware.Chain(
			handlers.DeleteRoleHandler(writer, roleSvc, logger),
			logging,
			requireAdmin,
		),
	)
}
