package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal/server"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() (*zap.Logger, error) {
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = timeEncoder
	return cfg.Build()
}

func main() {
	env := configs.ParseConfig()
	logger, err := getLogger()
	if err != nil {
		panic(err)
	}

	db, err := persistence.InitDatabase(env, logger)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = persistence.ShutdownDatabase(env, db, logger)
		if err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}()

	deps := NewAppDeps(db, logger)

	httpServer := server.NewHttpServer(env, logger)
	httpServer.RegisterRoutes(deps.UserSvc)

	go func() {
		err = httpServer.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown HTTP server", zap.Error(err))
	} else {
		logger.Info("HTTP server shut down")
	}
}
