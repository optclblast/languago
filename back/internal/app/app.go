package app

import (
	"context"
	"fmt"
	"languago/infrastructure/config"
	"languago/internal/server"
	errors2 "languago/pkg/errors"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout time.Duration = 5 * time.Second
)

type languagoApp struct {
	config *config.Config
}

func StartApp() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if config, ok := config.InitialConfiguration().(*config.Config); ok {
		app := languagoApp{config: config}
		return app.main(ctx)
	}

	return fmt.Errorf("error start application: can't get config")
}

func (a *languagoApp) main(ctx context.Context) error {
	logger := a.config.GetLoggerConfig().GetLogger()

	node := server.NewNode(&server.NewNodeParams{
		Logger:          logger,
		Config:          a.config,
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

	return nil
}
