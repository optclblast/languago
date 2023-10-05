package main

import (
	"context"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/server"
	"log"
)

func main() {
	// ctx, stop := signal.NotifyContext(
	// 	context.Background(),
	// 	syscall.SIGINT,
	// 	syscall.SIGTERM,
	// )
	// defer stop()
	ctx := context.Background()

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	cfg := config.InitialConfiguration()

	svs := map[string]server.Service{
		"flashcard_service": server.NewService(cfg),
	}

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
