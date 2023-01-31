package main

import (
	"context"
	http2 "mancala/handler/http"
	internalDB "mancala/internal/db"
	internalLog "mancala/internal/log"
	"mancala/internal/server"
	"net/http"
	"os/user"
	"time"
)

func main() {

	logger := internalLog.New()

	logger.Info("Starting the app")

	db, err := internalDB.New()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Migrate the schema
	_ = db.AutoMigrate(&user.User{})

	srv := server.New(http2.NewHandler(logger))
	logger.With("addr", srv.Addr).Info("Starting the server")

	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		<-ctx.Done()

		logger.Info("Stopping the server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.With("error", err).Fatal("Could not gracefully shutdown the server.")
		}
		logger.Info("Server stopped.")
		close(done)
	}(done)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("Server stopped")
	}
	<-done

}
