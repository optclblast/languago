package server

import "github.com/google/uuid"

type (
	Node interface {
		*Service
		Run() error
		Stop() error
		Healthcheck() error
	}

	node struct {
		Id uuid.UUID
		*Service
	}
)
