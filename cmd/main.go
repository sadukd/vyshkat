package main

import (
	"fmt"
	"os"

	"vyshkat/internal/console"
	"vyshkat/internal/inspection"
	"vyshkat/internal/service"
)

func main() {
	inspector := inspection.NewInspector()
	inventoryService := service.NewService(inspector)
	numberGenerator := service.NewNumberGenerator()
	app := console.New(os.Stdin, os.Stdout, inventoryService, numberGenerator)

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
