package main

import (
	"context"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	cfg := config.InitialConfiguration()

	svs := make(map[string]server.Service)
	svs["flashcard_service"] = server.NewService(cfg)

	node, err := server.NewNode(&server.NewNodeParams{
		Services: svs,
		Logger:   cfg.GetLoggerConfig().GetLogger(),
		Closer:   closer.NewCloser(),
	})
	if err != nil {
		return err
	}

	node.Run()
	return nil
}
