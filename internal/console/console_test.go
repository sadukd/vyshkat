package console_test

import (
	"bytes"
	"strings"
	"testing"

	"vyshkat/internal/console"
	"vyshkat/internal/domain"
	"vyshkat/internal/inspection"
	"vyshkat/internal/service"
)

func TestConsoleWorkflow(t *testing.T) {
	t.Parallel()

	input := strings.Join([]string{
		"1", "1", "City bicycle", "6", "7",
		"2", "Rejected scooter", "2.5", "8", "6",
		"2", "Good scooter", "2.5", "8", "8", "q",
		"2", "1", "Helmet", "2", "Station", "4.5", "q",
		"5", "1", "q", "2", "q", "3", "q", "4", "q", "5", "q", "q", "q",
	}, "\n") + "\n"

	svc, output := runConsole(t, input)

	if got := len(svc.Inventory()); got != 4 {
		t.Fatalf("inventory length = %d, want 4", got)
	}
	if got := svc.VehicleCount(); got != 2 {
		t.Errorf("VehicleCount() = %d, want 2", got)
	}
	if got := svc.TotalEnergyConsumption(); got != 7 {
		t.Errorf("TotalEnergyConsumption() = %v, want 7", got)
	}

	for _, want := range []string{
		"Added with inventory number 1.",
		"Not added: vehicle failed inspection",
		"Added with inventory number 3.",
		"#1 Bicycle: City bicycle, condition=7, ease=6",
		"#3 Scooter: Good scooter, condition=8, ease=8, energy=2.5",
		"#4 Helmet: Helmet",
		"#5 Charging station: Station, energy=4.5",
		"≡≡≡≡≡≡ Inventory count ≡≡≡≡≡≡\n4",
		"≡≡≡≡≡≡ Vehicle count ≡≡≡≡≡≡\n2",
		"≡≡≡≡≡≡ Total energy consumption (kW/h) ≡≡≡≡≡≡\n7",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output does not contain %q", want)
		}
	}
}

func TestConsoleRetriesInvalidFormInput(t *testing.T) {
	t.Parallel()

	input := "2\n2\n \nStation\ninvalid\n-1\n3.5\nq\nq\n"
	svc, output := runConsole(t, input)

	if got := len(svc.Inventory()); got != 1 {
		t.Fatalf("inventory length = %d, want 1", got)
	}
	if !strings.Contains(output, "Name cannot be empty.") {
		t.Error("missing empty-name validation message")
	}
	if got := strings.Count(output, "Enter a non-negative finite number."); got != 2 {
		t.Errorf("energy validation message count = %d, want 2", got)
	}
}

func TestConsoleRetriesInvalidVehicleValues(t *testing.T) {
	t.Parallel()

	input := "1\n1\nBicycle\ninvalid\n0\n6\n-1\n11\n7\nq\nq\n"
	svc, output := runConsole(t, input)

	if got := svc.VehicleCount(); got != 1 {
		t.Fatalf("VehicleCount() = %d, want 1", got)
	}
	if got := strings.Count(output, "Enter a whole number from 1 to 10."); got != 2 {
		t.Errorf("ease validation message count = %d, want 2", got)
	}
	if got := strings.Count(output, "Enter a whole number from 0 to 10."); got != 2 {
		t.Errorf("condition validation message count = %d, want 2", got)
	}
}

func TestConsoleShowsEmptyInventory(t *testing.T) {
	t.Parallel()

	_, output := runConsole(t, "5\n1\nq\nq\nq\n")
	want := "≡≡≡≡≡≡ All inventory ≡≡≡≡≡≡\n(empty)\n" + strings.Repeat("≡", 27) + "\n< previous"
	if !strings.Contains(output, want) {
		t.Error("empty inventory report is missing its empty state")
	}
}

func TestConsolePaginatesInventory(t *testing.T) {
	t.Parallel()

	svc := service.NewService(inspection.NewInspector())
	for number := 1; number <= 11; number++ {
		helmet, err := domain.NewHelmet("Helmet", number)
		if err != nil {
			t.Fatalf("NewHelmet() error = %v", err)
		}
		if err := svc.AddItem(helmet); err != nil {
			t.Fatalf("AddItem() error = %v", err)
		}
	}

	var output bytes.Buffer
	app := console.New(strings.NewReader("5\n1\n>\n<\nq\nq\nq\n"), &output, svc, service.NewNumberGenerator())
	if err := app.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	text := output.String()
	if !strings.Contains(text, "0. #10 Helmet: Helmet") {
		t.Error("first page does not use 0 for its tenth row")
	}
	if !strings.Contains(text, "1. #11 Helmet: Helmet") {
		t.Error("second page does not contain the eleventh inventory item")
	}
	if strings.Contains(text, "No next page.") || strings.Contains(text, "No previous page.") {
		t.Error("valid page navigation was rejected")
	}
}

func TestConsoleRemovesInventory(t *testing.T) {
	t.Parallel()

	input := strings.Join([]string{
		"1", "1", "Bicycle", "6", "7", "q",
		"2", "1", "Helmet", "q",
		"4", "1",
		"3", "2",
		"3", "1", "maybe", "Y",
		"4", "2", "n",
		"4", "2", "y",
		"4", "999",
		"q",
	}, "\n") + "\n"

	svc, output := runConsole(t, input)
	if got := svc.InventoryCount(); got != 0 {
		t.Fatalf("InventoryCount() = %d, want 0", got)
	}
	for _, want := range []string{
		"Not removed: vehicles must be removed through RemoveVehicle",
		"Not removed: inventory item is not a vehicle",
		"Selected: #1 Bicycle: Bicycle, condition=7, ease=6",
		"Enter y or n.",
		"Removed inventory number 1.",
		"Removal cancelled.",
		"Removed inventory number 2.",
		"Not removed: inventory item not found",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output does not contain %q", want)
		}
	}
}

func TestConsoleRetriesInvalidRemovalNumber(t *testing.T) {
	t.Parallel()

	input := "2\n1\nHelmet\nq\n4\ninvalid\n0\n1\nn\nq\n"
	svc, output := runConsole(t, input)
	if got := svc.InventoryCount(); got != 1 {
		t.Fatalf("InventoryCount() = %d, want 1", got)
	}
	if got := strings.Count(output, "Enter a whole number from 1 to"); got != 2 {
		t.Errorf("inventory-number validation count = %d, want 2", got)
	}
}

func TestConsoleEOFExitsCleanly(t *testing.T) {
	t.Parallel()

	svc := service.NewService(inspection.NewInspector())
	app := console.New(strings.NewReader(""), &bytes.Buffer{}, svc, service.NewNumberGenerator())
	if err := app.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func runConsole(t *testing.T, input string) (*service.Service, string) {
	t.Helper()

	svc := service.NewService(inspection.NewInspector())
	var output bytes.Buffer
	app := console.New(strings.NewReader(input), &output, svc, service.NewNumberGenerator())
	if err := app.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	return svc, output.String()
}
