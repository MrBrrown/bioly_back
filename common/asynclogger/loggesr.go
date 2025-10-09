package asynclogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LL_FATALERROR = iota
	LL_ERROR
	LL_WARNING
	LL_INFO
)

var logLevelNames = []string{
	"[LL_FATALERROR]", // 0
	"[LL_ERROR]",      // 1
	"[LL_WARNING]",    // 2
	"[LL_INFO]",       // 3
}

var (
	logChan = make(chan string, 100)
	wg      sync.WaitGroup

	maxLevelNameLen int
)

type LoggerInfo struct {
	FilePath   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	IsCompress bool   `json:"compressed"`
}

func New(info LoggerInfo) io.WriteCloser {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s%s", info.FilePath, timestamp, "access.log")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("failed to create log file: " + err.Error())
		}
		_ = f.Close()
	}
	_ = os.Chmod(fileName, 0644)
	ljLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    info.MaxSize,
		MaxBackups: info.MaxBackups,
		MaxAge:     info.MaxAge,
		Compress:   info.IsCompress,
	}

	for _, lvl := range logLevelNames {
		if len(lvl) > maxLevelNameLen {
			maxLevelNameLen = len(lvl)
		}
	}

	maxLevelNameLen = maxLevelNameLen + 1

	return struct {
		io.Writer
		io.Closer
	}{
		Writer: io.MultiWriter(os.Stdout, ljLogger),
		Closer: ljLogger,
	}
}

func StartAsyncLogWriter(out io.Writer) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range logChan {
			fmt.Fprintln(out, msg)
		}
	}()
}

func ShutdownLogger() {
	close(logChan)
	wg.Wait()
}

func flush() {
	close(logChan)
	wg.Wait()
}

func getLogLevelName(level int) string {
	if level >= 0 && level < len(logLevelNames) {
		return logLevelNames[level]
	}
	return "UNKNOWN"
}

func logf(level int, format string, args ...any) {
	timestamp := time.Now().UTC().Format("2006/01/02 15:04:05 UTC+0")
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("%s %*s %s", timestamp, maxLevelNameLen, getLogLevelName(level), msg)

	select {
	case logChan <- line:
	default:
		log.Printf("%*s %s", maxLevelNameLen, getLogLevelName(level), msg)
	}
}

func Log(level int, format string, args ...any) {
	logf(level, format, args...)
}

// Log with exit from app
func Fatal(format string, args ...any) {
	logf(LL_FATALERROR, format, args...)
	flush()
	os.Exit(1)
}

func FatalError(format string, args ...any) {
	logf(LL_FATALERROR, format, args...)
}

func Error(format string, args ...any) {
	logf(LL_ERROR, format, args...)
}

func Warning(format string, args ...any) {
	logf(LL_WARNING, format, args...)
}

func Info(format string, args ...any) {
	logf(LL_INFO, format, args...)
}
