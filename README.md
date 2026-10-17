# VyshKat

VyshKat is a Go console application developed as a coursework project for managing the inventory of a small kick-sharing service. The application stores vehicles and supporting equipment, assigns inventory numbers, requires a technical inspection before accepting vehicles, and provides operational inventory reports.

## Features

- Add bicycles and scooters to inventory after mandatory technical inspection.
- Add helmets and charging stations as non-vehicle inventory items.
- Generate sequential inventory numbers that are unique within the running application.
- Display all inventory with pagination of 10 records per page.
- Remove vehicles and non-vehicle items by inventory number after confirmation.
- Count all inventory records and vehicles separately.
- Sum the daily energy consumption of every object implementing `domain.EnergyConsumer`.
- List vehicles suitable for beginners.

Inventory numbers start at `1` and increase monotonically. A number is generated before a vehicle is inspected, so a rejected vehicle consumes its number and may leave a gap. Removing an object does not make its number available for reuse.

## Domain rules and validation

The constructors in `internal/domain` protect the entity invariants, while technical admission is handled separately by `internal/inspection`.

| Rule | Implemented constraint |
|---|---|
| Item name | Must contain at least one non-whitespace character. |
| Inventory number in domain constructors | Must be non-negative. The production generator issues positive numbers starting at `1`. |
| Ease of use | Integer from `1` to `10`, inclusive. |
| Technical condition | Integer from `0` to `10`, inclusive. |
| Vehicle inspection | A vehicle is accepted only when `Condition() >= 7`. |
| Beginner suitability | A stored vehicle is included when `EaseOfUse() >= 6`. |
| Energy consumption | Must be finite and non-negative; negative values, `NaN`, and infinities are rejected. |
| Inventory uniqueness | `service.Service` rejects duplicate inventory numbers. |

The inspection threshold and beginner threshold represent different rules. A vehicle must first pass inspection with condition at least `7`; among accepted vehicles, ease of use at least `6` makes it suitable for beginners.

The intended business unit for `EnergyConsumption()` and `TotalEnergyConsumption()` is **kWh/day**. The service calculates the total by summing the values of all inventory objects implementing `EnergyConsumer`, currently scooters and charging stations.

## Requirements, installation, and running

The module declares **Go 1.22.2** in `go.mod`. A compatible Go installation is required.

From a checked-out copy of the repository:

```bash
go mod download
go run ./cmd
```

To build and run a local executable:

```bash
go build -o vyshkat ./cmd
./vyshkat
```

The application keeps all data in memory. Stopping the process clears the inventory.

## Console interface

The root menu contains five operations:

1. `Add vehicle`
2. `Add item`
3. `Remove vehicle`
4. `Remove item`
5. `Reports`

`Add vehicle` supports `Bicycle` and `Scooter`. Both require a name, ease of use, and technical condition; a scooter additionally requires daily energy consumption. The object receives an inventory number and is then passed to `Service.AddVehicle`, which invokes the configured `VehicleInspector`.

`Add item` supports `Helmet` and `Charging station`. A charging station additionally requires daily energy consumption. Non-vehicle items are added through `Service.AddItem` and do not undergo vehicle inspection.

Removal is separated by category. The user enters an inventory number, the selected object is shown, and deletion requires `y`/`n` confirmation. `Remove vehicle` rejects equipment, while `Remove item` rejects vehicles.

The reports submenu provides:

- all inventory;
- vehicles suitable for beginners;
- total inventory count;
- vehicle count;
- total energy consumption.

Menus and inventory reports support pages of up to 10 entries. Keys `1`–`9` select the first nine entries, `0` selects the tenth entry, `<` and `>` change pages, and `q` returns to the previous menu or exits from the root menu.

### Example interaction

```text
≡≡≡≡≡≡ Inventory ≡≡≡≡≡≡
1. Add vehicle
2. Add item
3. Remove vehicle
4. Remove item
5. Reports
≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡
q. Exit
> 1

≡≡≡≡≡≡ Add vehicle ≡≡≡≡≡≡
1. Bicycle
2. Scooter
≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡
q. Back
> 1
Name: City bicycle
Ease of use (1-10): 6
Condition (0-10): 7
Added with inventory number 1.
```

## Architecture

The project separates domain rules, use cases, infrastructure, and user interaction.

| Package | Responsibility |
|---|---|
| `internal/domain` | Defines `InventoryItem`, `Vehicle`, and `EnergyConsumer`; implements `Bicycle`, `Scooter`, `Helmet`, and `ChargingStation`; validates entity construction. |
| `internal/inspection` | Implements the production technical-inspection policy in `Inspector.Inspect`, using `MinAcceptedCondition = 7`. |
| `internal/service` | Owns the in-memory inventory and application use cases: addition, lookup, removal, counts, energy aggregation, and beginner filtering. It also contains `NumberGenerator`. |
| `internal/console` | Handles terminal input and output, paginated menus, input validation, entity creation flows, confirmations, and presentation formatting. |
| `internal/app` | Acts as the composition root and assembles the dependency graph with `github.com/samber/do/v2`. |
| `cmd` | Contains the executable entry point, resolves the console application, runs it, and shuts down the DI container. |

The service returns a copy of its inventory slice from `Inventory()`, preventing callers from replacing entries in the internal slice. Storage itself remains in memory and is owned by one `Service` instance.

## Application of SOLID principles

The project applies SOLID where it helps separate the coursework responsibilities; it does not attempt to eliminate every future modification.

### Single Responsibility Principle

- Domain constructors validate entity state.
- `inspection.Inspector` decides whether a vehicle passes technical inspection.
- `service.Service` implements inventory use cases and queries.
- `console.Console` owns terminal interaction and presentation.
- `app.NewContainer` is responsible only for application composition.

This keeps inspection policy out of both entity constructors and console code, and keeps terminal concerns out of the service.

### Open/Closed Principle

Service reports depend on capabilities rather than concrete types. Any stored object implementing `EnergyConsumer` contributes to `TotalEnergyConsumption`, and any `Vehicle` participates in vehicle counting and beginner filtering. These algorithms do not need a concrete-type branch for a new implementation.

The console boundary is intentionally less open: making a new type creatable requires a new input handler and menu entry. `formatItem` also needs a type case to show a specialized label; otherwise the object is displayed as `Item` while its shared capabilities are still formatted.

### Liskov Substitution Principle

`Bicycle`, `Scooter`, and test doubles can be passed anywhere a `domain.Vehicle` is expected. The inspector and service operate only through `Name`, `Number`, `EaseOfUse`, and `Condition`. Similarly, scooters, charging stations, and energy-consuming test doubles are interchangeable for aggregation through `EnergyConsumer`.

### Interface Segregation Principle

The domain uses small capability interfaces:

- `InventoryItem` exposes only identity data;
- `Vehicle` extends it with condition and ease of use;
- `EnergyConsumer` exposes only energy consumption.

The console also depends on the narrow `NumberProvider` interface rather than the concrete generator. `console.InventoryService` groups exactly the use cases required by the CLI, while other packages are not forced to depend on that interface.

### Dependency Inversion Principle

`service.Service` depends on the `service.VehicleInspector` abstraction rather than `inspection.Inspector`. `console.Console` depends on `console.InventoryService` and `console.NumberProvider`, not on concrete service implementations. High-level workflows can therefore receive alternative implementations without importing infrastructure details.

## Dependency injection

`internal/app/container.go` uses `github.com/samber/do/v2` as the composition mechanism. `NewContainer` registers:

- injected `io.Reader` and `io.Writer` values;
- `inspection.Inspector` as `service.VehicleInspector`;
- `service.Service` and its `console.InventoryService` alias;
- `service.NumberGenerator` and its `console.NumberProvider` alias;
- `console.Console` assembled from those dependencies.

The concrete choices are confined to the composition root. The domain, inspection, service, and console packages do not import the DI library. `cmd/main.go` resolves `*console.Console`, runs it, and defers container shutdown.

Tests replace dependencies without changing production code. Service tests pass `inspectorMock` directly to `service.NewService`; console tests provide an in-memory service, `strings.Reader`, and `bytes.Buffer`; the container test verifies that the complete graph resolves and that repeated service resolution returns the same instance.

## Testing

The test suite includes:

- table-driven domain tests for valid and invalid entity construction;
- compile-time interface conformance checks;
- boundary tests for the inspection threshold;
- service tests using mocks for inspection, inventory capabilities, removal errors, ordering, counts, energy totals, and beginner filtering;
- number-generator sequence and independence tests;
- scripted console workflow tests covering validation retries, reports, pagination, removal, confirmation, and EOF handling;
- a DI-container integration test.

Run the standard checks with:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Generate combined coverage for all packages with:

```bash
go test ./... -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

At the time of this README update, `go test ./...`, `go test -race ./...`, and `go vet ./...` completed successfully. The combined statement coverage reported by the commands above was **83.7%**. No coverage threshold is enforced by the project.

## Extension scenarios

### New vehicle type

A new transport type can implement `domain.Vehicle`; if it consumes energy, it can additionally implement `domain.EnergyConsumer`. Inspection, vehicle counting, beginner filtering, and energy aggregation will work without changes to their business algorithms. Console creation still requires a form method, a menu entry, and optionally a specialized `formatItem` case.

### New equipment type

New equipment can implement `domain.InventoryItem` and be handled by `Service.AddItem`. Implementing `EnergyConsumer` automatically includes it in the total. Creating it interactively still requires console code, and a specialized display name requires formatter support.

### Employees

Employees are not inventory objects and are not implemented. A suitable extension would introduce a separate employee model and service rather than make employees implement `InventoryItem`. The composition root and root console menu could then connect the new module to the application.

### Spare-parts management

Spare parts and warehouse operations are not implemented. They can be introduced as a separate aggregate and service, with dedicated stock operations and console workflows. Existing fleet entities and inspection rules need not change unless a new use case explicitly coordinates both modules.

## Limitations and known discrepancies

- Inventory is stored only in memory; there is no database, file storage, loading, or saving.
- Application state is local to one process, with no authentication, authorization, networking, or multi-user synchronization.
- Inventory numbers are not persisted and restart from `1` after each application start.
- Generated numbers are not reused after rejection or deletion.
- Supported console creation forms are fixed to bicycles, scooters, helmets, and charging stations.
- Employees and spare-parts accounting are extension scenarios, not existing functionality.
- The business meaning of energy consumption is kWh/day, but the current report title displays `kW/h`. This is a presentation-unit discrepancy in the implementation; the README does not silently correct the program output.
- Domain constructors allow inventory number `0`, although the production generator and removal form use positive numbers starting at `1`.

## Use of AI tools

AI assistance was used during the coursework process to help extend and review tests, examine architectural and SOLID decisions, integrate and review dependency-injection wiring, identify focused implementation issues, and prepare documentation. This statement does not attribute the entire implementation to either AI or a particular person. The submitted project should be reviewed and understood by the student responsible for it.
