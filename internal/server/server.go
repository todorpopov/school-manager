package server

import (
	"context"
	"net/http"

	"github.com/todorpopov/school-manager/configs"
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
	env    *configs.Config
	mux    *http.ServeMux
	writer *writer.HttpWriter
	logger *zap.Logger
	server *http.Server
}

func NewHttpServer(env *configs.Config, logger *zap.Logger) *HttpServer {
	return &HttpServer{
		env,
		http.NewServeMux(),
		writer.NewHttpWriter(),
		logger,
		nil,
	}
}

func (s *HttpServer) Start() error {
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

func (s *HttpServer) RegisterRoutes() {
	routes.RegisterGeneralRoutes(s.mux, s.writer, s.logger)
}
