package app

import (
	"context"
	"fmt"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
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
	log := logger.ProvideLogger(a.config.LoggerCfg)

	node := server.NewNode(&server.NewNodeParams{
		Log:             logger.ProvideLogger(a.config.LoggerCfg),
		Logger:          a.config.GetLoggerConfig().GetLogger(),
		Config:          a.config,
		ErrorsPresenter: errors2.NewErrorPresenter(log),
	})

	go func() {
		node.Run()
	}()

	<-ctx.Done()
	func() {
		closeTimer := time.AfterFunc(shutdownTimeout, func() {
			node.Log().Warn().Msg("error node shutdown timeout. stopping node with force")
			node.Stop(ctx) // todo force
		})
		defer closeTimer.Stop()

		node.Stop(ctx)
	}()

	node.Log().Info().Msgf("node shutting down | node_id: %v time: %v", node.ID().String(), time.Now())

	return nil
}
