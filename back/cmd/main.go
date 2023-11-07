package main

import (
	"context"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/server"
	"log"
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

	if err := Run(ctx); err != nil {
		log.Fatal(err)
	}

	server := 
}

func Run(ctx context.Context) error {
	cfg := config.InitialConfiguration()

	node := server.NewNode(&server.NewNodeParams{
		Logger: cfg.GetLoggerConfig().GetLogger(),
		Closer: closer.NewCloser(),
	})

	node.ServiceBuilder(cfg)

	go func() {
		node.Run()
	}()
	<-ctx.Done()

	node.Log().Info("node shutting down", logger.LogFields{
		"node_id": node.ID().String(),
		"time":    time.Now(),
	})

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	node.Stop(shutdownCtx)
	node.Log().Info("node stoped", nil)

	return nil
}
