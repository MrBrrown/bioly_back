package main

import (
	"log"
	"os"
	"time"

	"bioly/common/asynclogger"
	"bioly/common/yamlconf"

	"bioly/gateway/internal/config"
	"bioly/gateway/internal/server"
	"bioly/gateway/internal/transport"
)

func main() {
	cfgFile := os.Getenv("CONFIG_PATH")
	cfg := &config.Config{}
	err := yamlconf.Load(cfgFile, cfg)
	if err != nil {
		asynclogger.Fatal("failed to load config: %s", err.Error())
	}

	logDirName := os.Getenv("LOG_DIR")
	loggerInfo := asynclogger.LoggerInfo{FilePath: logDirName, MaxSize: 10, MaxBackups: 5, MaxAge: 30, IsCompress: true}
	logger := asynclogger.New(loggerInfo)
	log.SetOutput(logger)
	asynclogger.StartAsyncLogWriter(logger)
	defer func() {
		asynclogger.ShutdownLogger()
		logger.Close()
	}()

	handler, err := transport.NewHandler(cfg.Proxies)
	if err != nil {
		asynclogger.Fatal("failed to create handler: %s", err.Error())
	}

	corsCfg := server.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "application/json"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}
	requestTimeout := 10 * time.Second

	router := server.NewRouter(handler, requestTimeout, corsCfg)

	if err := server.Start(cfg.Server, router); err != nil {
		asynclogger.FatalError("server error: %s", err.Error())
	}
}
