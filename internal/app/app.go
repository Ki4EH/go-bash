package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ki4EH/go-bash/internal/config"
	"github.com/Ki4EH/go-bash/internal/logger"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

type App struct {
	server *http.Server
	Db     *sql.DB
	Router *mux.Router
}

// Run SetupRoutes sets up the routes for the application
func Run(cfg *config.Config) (*App, error) {
	logger.Info("starting http server...")

	srv := new(App)

	srv.server = &http.Server{
		Addr:         cfg.Address,
		Handler:      srv.SetupRoutes(),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := ConnectionToDB(srv, cfg.Database); err != nil {
		return nil, fmt.Errorf("error on connect to db %w", err)
	}

	if err := srv.server.ListenAndServe(); err != nil {
		return nil, fmt.Errorf("error on starting server %w", err)
	}

	logger.Info("server was shutdown")
	return srv, nil
}

// ConnectionToDB connects to the database and sets the connection to the App
func ConnectionToDB(srv *App, dbStruct config.Database) error {
	logger.Info("connecting to db...")
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbStruct.Host, dbStruct.Port, dbStruct.UserName, dbStruct.Password, dbStruct.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	srv.Db = db

	return nil
}

// Stop stops the server
func (a *App) stop(ctx context.Context) error {
	logger.Info("shutdown server...")
	err := a.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}
	logger.Info("server was shutdown")
	return nil
}

// GracefulStop stops the server gracefully
func (a *App) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 30 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			logger.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	err := a.stop(shutdownCtx)
	if err != nil {
		logger.Error("graceful shutdown timed out... forcing exit")
		os.Exit(1)
	}
	defer serverStopCtx()
	defer shutdownStopCtx()
}
