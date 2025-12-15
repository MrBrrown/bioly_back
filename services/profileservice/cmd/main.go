package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bioly/common/asynclogger"
	"bioly/common/storage"
	"bioly/profileservice/internal/cache"
	"bioly/profileservice/internal/config"
	"bioly/profileservice/internal/repositories"
	"bioly/profileservice/internal/transport"
	"bioly/profileservice/internal/usecases"
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
		asynclogger.Fatal("Can't connect to profile DB: %v", err)
	}

	lruCache := cache.NewLruProfileCache(10000)

	profile := repositories.NewProfile(db)
	service := usecases.NewProfile(profile, lruCache)

	handler := transport.NewHandler(service)
	router := transport.NewRouter(handler)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	asynclogger.Info("Starting profile service on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		asynclogger.Fatal("Server stopped: %v", err)
	}
}
