package app

import (
	"context"
	"fmt"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const (
	shutdownTimeout time.Duration = 10 * time.Second
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

	config := config.InitialConfiguration()
	app := languagoApp{config: config}

	return app.main(ctx)
}

func (a *languagoApp) main(ctx context.Context) error {
	log := logger.ProvideLogger(a.config)
	log.Info("starting")

	consul, _ := config.NewConsulManager("127.0.0.1:8500", "iii", log)
	consul.RegisterNode(&config.RegisterNodeParams{
		ID:      fmt.Sprintf("languago-server-%s", uuid.NewString()),
		Name:    "languago-cluster",
		Tags:    []string{"languago-server"},
		Address: "127.0.0.1",
		Port:    3301,
	})

	consul.DefaultConfiguration()

	time.Sleep(30 * time.Second)

	// node := server.NewNode(&server.NewNodeParams{
	// 	Version:         a.config.Node.Version,
	// 	Log:             log,
	// 	Config:          a.config,
	// 	ErrorsPresenter: errors2.NewErrorPresenter(log),
	// })

	// go func() {
	// 	node.Run()
	// }()

	// <-ctx.Done()
	// func() {
	// 	closeTimer := time.AfterFunc(shutdownTimeout, func() {
	// 		node.Log().Warn("error node shutdown timeout. stopping node with force")
	// 		node.Stop(ctx) // todo force
	// 	})
	// 	defer closeTimer.Stop()

	// 	node.Stop(ctx)
	// }()

	// node.Log().Infof("node shutting down | node_id: %v time: %v", node.ID().String(), time.Now())

	return nil
}
