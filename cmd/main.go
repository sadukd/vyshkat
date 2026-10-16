package main

import (
	"fmt"
	"os"

	"github.com/samber/do/v2"

	"vyshkat/internal/app"
	"vyshkat/internal/console"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	injector := app.NewContainer(os.Stdin, os.Stdout)
	defer injector.Shutdown()

	application, err := do.Invoke[*console.Console](injector)
	if err != nil {
		return err
	}

	return application.Run()
}
