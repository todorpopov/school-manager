package main

import (
	"time"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal/server"
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

	httpServer := server.NewHttpServer(env, logger)
	httpServer.RegisterRoutes()
	err = httpServer.Start()
	if err != nil {
		panic(err)
	}
}
