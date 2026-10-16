# VyshKat

VyshKat is a console application for managing a small kick-sharing fleet. It accepts vehicles only after a technical inspection, keeps inventory numbers for vehicles and equipment, calculates daily energy consumption, and identifies vehicles suitable for beginners.

## Run and test

The project requires Go 1.22 or newer.

```bash
go run ./cmd
go test ./...
go test -race ./...
go vet ./...
```

To inspect combined coverage across the project:

```bash
go test ./... -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

The console supports these scenarios:

- add a bicycle or scooter, with acceptance determined by technical inspection;
- add a helmet or charging station to the balance sheet;
- remove vehicles or other inventory items by inventory number after confirmation;
- list all inventory and beginner-friendly vehicles;
- report total inventory count, vehicle count, and daily energy consumption.

Inventory numbers are generated sequentially. A rejected vehicle still consumes its generated number, so gaps are expected.

## Architecture

The application is split into four layers:

- `internal/domain` contains entities and small capability interfaces: `InventoryItem`, `Vehicle`, and `EnergyConsumer`;
- `internal/service` implements use cases and declares the `VehicleInspector` interface it consumes;
- `internal/inspection` contains the production inspection policy;
- `internal/console` handles input, output, object creation, and presentation.

`internal/app` is the composition root. It uses `github.com/samber/do/v2` to register the inspector, service, number generator, console, and I/O streams. The DI container is not imported by the domain, service, inspection, or console packages. Tests can inject an inspector mock directly into `service.NewService`.

The inventory service stores data in memory. Its query methods work through capabilities rather than concrete entity types: every `EnergyConsumer` contributes to the energy total, while every `Vehicle` participates in vehicle and beginner reports.

## SOLID decisions

- **Single Responsibility:** domain objects protect their invariants, the inspector decides technical acceptance, the service owns inventory use cases, and the console owns interaction.
- **Open/Closed:** adding a type that implements `Vehicle` or `EnergyConsumer` does not require changes to inspection or report algorithms. The service works through interfaces and type capabilities.
- **Liskov Substitution:** bicycles, scooters, and test doubles can all replace `Vehicle` without changing service behavior. The interfaces expose only operations whose contracts every implementation can satisfy.
- **Interface Segregation:** inventory identity, energy consumption, vehicle properties, console use cases, and number generation are separate narrow interfaces.
- **Dependency Inversion:** `Service` depends on `VehicleInspector`, not the concrete inspector. `Console` depends on its own `InventoryService` and `NumberProvider` interfaces. Production implementations are selected only in the composition root.

DIP is the most explicit principle: the inspection dependency is replaced by a mock in `TestAddVehicle`, and application dependencies are assembled by the DI container. OCP is also central to the capability-based reports.

OCP is intentionally limited at the console boundary. A new vehicle automatically works with service-level inspection and reports, but making it creatable from the console requires a menu option and an input form. The item formatter also needs a case if a type-specific label is desired; without one, the item is still listed with its name and inventory number. This explicit presentation mapping is preferred over adding UI concerns to domain interfaces.

## Extension scenarios

- A new transport type implements `domain.Vehicle`; if it consumes energy, it also implements `domain.EnergyConsumer`. Service logic and existing reports remain unchanged, while the console gains the corresponding creation flow.
- An employee does not implement `InventoryItem`, so it cannot accidentally be added to the balance sheet. Employee accounting would be a separate module.
- A spare-parts warehouse would be a separate aggregate and service. It can be added without changing the current fleet model; only combined application workflows would depend on both modules.

The technical-inspection decision is covered with an injected `inspectorMock`, including accepted and rejected paths. The beginner boundary is fixed by a service test in which ease `6` is included and ease `5` is excluded.

## AI usage

The original application was implemented by the author. AI was used to help create tests and, during the final review, to audit the requirements and architecture, suggest and implement focused fixes, add DI-container wiring, extend tests, and draft documentation. The author is responsible for reviewing, understanding, and defending the final design and code.
