package main

import (
	"context"
	"go.uber.org/zap"
	http2 "mancala/handler/http"
	internalDB "mancala/internal/db"
	"mancala/internal/server"
	"mancala/lobby"
	"mancala/mancala"
	"net/http"
	"time"
)

func main() {

	l, _ := zap.NewDevelopment()
	defer l.Sync()
	zap.ReplaceGlobals(l)
	logger := l.Sugar()

	logger.Info("Starting the app")

	db, err := internalDB.New()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err.Error())
	}

	ctx := context.Background()

	// Migrate the schema
	_ = db.AutoMigrate(&lobby.Lobby{}, &mancala.MancalaDB{})

	lobbyService := lobby.NewService(db)
	gameService := mancala.NewService(db)

	srv := server.New(http2.NewHandler(logger, lobbyService, gameService))
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
		logger.With("error", err).Info("Server stopped")
	}
	<-done

}
