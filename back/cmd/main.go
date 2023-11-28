package main

import (
	"fmt"
	"languago/internal/app"
	"os"
)

func main() {
	if err := app.StartApp(); err != nil {
		fmt.Fprintf(os.Stderr, "error at application runtime: %s", err.Error())
		os.Exit(1)
	}
}
