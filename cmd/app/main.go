package main

import (
	"context"
	"fmt"
	"github.com/Ki4EH/go-bash/internal/app"
	"github.com/Ki4EH/go-bash/internal/config"
	"github.com/Ki4EH/go-bash/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Info("reading config...")
	conf, err := config.LoadFromEnv()
	if err != nil {
		fmt.Println(err)
		logger.Error("failed to read config, ", err)
	}

	srv, err := app.Run(conf)

	if err != nil {
		fmt.Println(err)
		logger.Error(err)
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	// Notify the channel when a signal is received
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	done := make(chan struct{})
	// Graceful shutdown of the server when a signal is received or the context is done
	go func() {
		srv.GracefulStop(serverCtx, sig, serverStopCtx)
		serverStopCtx()
		close(done)
	}()

	<-done
}
