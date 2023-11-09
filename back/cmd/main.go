package main

import (
	"context"
	"languago/internal/pkg/config"
	errors2 "languago/internal/pkg/errors"
	"languago/internal/server"

	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout time.Duration = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.InitialConfiguration()
	logger := cfg.GetLoggerConfig().GetLogger()

	node := server.NewNode(&server.NewNodeParams{
		Logger:          logger,
		Config:          cfg,
		ErrorsPresenter: errors2.NewErrorPresenter(logger),
	})

	go func() {
		node.Run()
	}()
	<-ctx.Done()

	node.Log().Info("node shutting down", map[string]any{
		"node_id": node.ID().String(),
		"time":    time.Now(),
	})
}
