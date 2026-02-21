package main

import (
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
	defer func(env *configs.Config, db *persistence.Database) {
		err = persistence.CloseDatabase(env, db, logger)
		if err != nil {
			panic(err)
		}
	}(env, db)

	httpServer := server.NewHttpServer(env, logger)
	httpServer.RegisterRoutes()
	err = httpServer.Start()
	if err != nil {
		panic(err)
	}
}
