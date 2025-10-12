package main

import (
	"log"
	"os"

	"bioly/common/asynclogger"
)

func main() {
	logDirName := os.Getenv("LOG_DIR")
	loggerInfo := asynclogger.LoggerInfo{FilePath: logDirName, MaxSize: 10, MaxBackups: 5, MaxAge: 30, IsCompress: true}
	logger := asynclogger.New(loggerInfo)
	log.SetOutput(logger)
	asynclogger.StartAsyncLogWriter(logger)
	defer func() {
		asynclogger.ShutdownLogger()
		logger.Close()
	}()
}
