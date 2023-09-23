package main

import (
	"languago/internal/pkg/config"
	"languago/internal/server"
)

func main() {
	cfg := config.InitialConfiguration()

	svs := make(map[string]server.Service)
	svs["flashcard_service"] = server.NewService(cfg)
	node, err := server.NewNode(&server.NewNodeParams{
		Services: svs,
		Logger:   cfg.GetLoggerConfig().GetLogger(),
	})
	if err != nil {
		panic("ERROR!!!")
	}

	node.Run()
}
