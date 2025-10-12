package server

import (
	"bioly/common/asynclogger"
	"bioly/gateway/internal/config"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start(serInfo config.Server, handler http.Handler) error {
	srv := &http.Server{
		Addr:         serInfo.Address,
		Handler:      handler,
		ReadTimeout:  serInfo.ReadTimeout,
		WriteTimeout: serInfo.WriteTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		asynclogger.Info("starting server on %s", serInfo.Address)
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		asynclogger.Info("shutting down server... signal: %s", sig.String())
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), serInfo.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(ctx)
}
