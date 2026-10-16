package app

import (
	"io"

	"github.com/samber/do/v2"

	"vyshkat/internal/console"
	"vyshkat/internal/inspection"
	"vyshkat/internal/service"
)

func NewContainer(input io.Reader, output io.Writer) *do.RootScope {
	injector := do.New()

	do.ProvideValue[io.Reader](injector, input)
	do.ProvideValue[io.Writer](injector, output)
	do.Provide[service.VehicleInspector](injector, func(do.Injector) (service.VehicleInspector, error) {
		return inspection.NewInspector(), nil
	})
	do.Provide[*service.Service](injector, func(i do.Injector) (*service.Service, error) {
		inspector, err := do.Invoke[service.VehicleInspector](i)
		if err != nil {
			return nil, err
		}
		return service.NewService(inspector), nil
	})
	do.Provide[console.InventoryService](injector, func(i do.Injector) (console.InventoryService, error) {
		return do.Invoke[*service.Service](i)
	})
	do.Provide[*service.NumberGenerator](injector, func(do.Injector) (*service.NumberGenerator, error) {
		return service.NewNumberGenerator(), nil
	})
	do.Provide[console.NumberProvider](injector, func(i do.Injector) (console.NumberProvider, error) {
		return do.Invoke[*service.NumberGenerator](i)
	})
	do.Provide[*console.Console](injector, func(i do.Injector) (*console.Console, error) {
		reader, err := do.Invoke[io.Reader](i)
		if err != nil {
			return nil, err
		}
		writer, err := do.Invoke[io.Writer](i)
		if err != nil {
			return nil, err
		}
		inventoryService, err := do.Invoke[console.InventoryService](i)
		if err != nil {
			return nil, err
		}
		numberProvider, err := do.Invoke[console.NumberProvider](i)
		if err != nil {
			return nil, err
		}
		return console.New(reader, writer, inventoryService, numberProvider), nil
	})

	return injector
}
