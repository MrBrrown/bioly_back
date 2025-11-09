package main

import (
	"bioly/asynclogger"
	"bioly/auth/internal/config"
	"bioly/auth/internal/repositories"
	"bioly/auth/internal/transport"
	"bioly/auth/internal/usecase"
	"bioly/storage"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	cfgFile := os.Getenv("CONFIG_PATH")
	cfg := config.New(cfgFile)

	logDirName := os.Getenv("LOG_DIR")
	loggerInfo := asynclogger.LoggerInfo{FilePath: logDirName, MaxSize: 10, MaxBackups: 5, MaxAge: 30, IsCompress: true}
	logger := asynclogger.New(loggerInfo)
	log.SetOutput(logger)
	asynclogger.StartAsyncLogWriter(logger)
	defer func() {
		asynclogger.ShutdownLogger()
		logger.Close()
	}()

	db, err := storage.New(&cfg.DBInfo)
	if err != nil {
		asynclogger.Fatal("Can't connect to auth DB: %v", err)
	}

	userRepo := repositories.NewUsers(db)
	refreshRepo := repositories.NewRefreshTokens(db)

	uc := usecase.NewAuth(userRepo, refreshRepo, &cfg.JWT)

	handler := transport.NewHandler(uc)
	router := transport.NewRouter(handler)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	asynclogger.Info("Starting auth service on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		asynclogger.Fatal("Server stopped: %v", err)
	}
}
