package server

import (
	"context"
	"net/http"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/routes"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	RegisterRoutes()
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	env        *configs.Config
	mux        *http.ServeMux
	writer     *writer.HttpWriter
	logger     *zap.Logger
	server     *http.Server
	serverDeps *Dependencies
}

func NewHttpServer(env *configs.Config, logger *zap.Logger, deps *Dependencies) *HttpServer {
	errWriter := exceptions.NewErrorWriter()
	return &HttpServer{
		env,
		http.NewServeMux(),
		writer.NewHttpWriter(errWriter),
		logger,
		nil,
		deps,
	}
}

func (s *HttpServer) Start() error {
	s.registerRoutes()
	s.logger.Info("Starting HTTP server on port: " + s.env.ApiPort)
	srvAddr := ":" + s.env.ApiPort
	s.server = &http.Server{
		Addr:    srvAddr,
		Handler: s.mux,
	}
	return s.server.ListenAndServe()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) registerRoutes() {
	routes.RegisterGeneralRoutes(s.mux, s.writer, s.logger)
	routes.RegisterUserRoutes(s.mux, s.writer, s.logger, s.serverDeps.UserSvc, s.serverDeps.AuthSvc)
	routes.RegisterAuthRoutes(s.mux, s.writer, s.logger, s.serverDeps.AuthSvc)
	routes.RegisterRoleRoutes(s.mux, s.writer, s.logger, s.serverDeps.RoleRepo, s.serverDeps.AuthSvc)
	routes.RegisterDirectorRoutes(s.mux, s.writer, s.logger, s.serverDeps.DirectorSvc, s.serverDeps.AuthSvc)
}
