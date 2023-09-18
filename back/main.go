package main

import (
	"languago/internal/server"
)

func main() {
	s := server.NewService()
	s.Run()
}
