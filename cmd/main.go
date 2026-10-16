package main

import (
	"fmt"
	"os"

	"github.com/samber/do/v2"

	"vyshkat/internal/app"
	"vyshkat/internal/console"
)

func main() {
	injector := app.NewContainer(os.Stdin, os.Stdout)
	application, err := do.Invoke[*console.Console](injector)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
