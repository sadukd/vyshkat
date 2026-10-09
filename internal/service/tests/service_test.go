package tests

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/service"
)

func TestAddVehicle(t *testing.T) {
	t.Parallel()

	t.Run("adds vehicle after successful inspection", func(t *testing.T) {
		inspector := &inspectorMock{result: true}
		svc := service.NewService(inspector)
		vehicle := &vehicleMock{itemMock: itemMock{name: "bicycle", number: 1}, easeOfUse: 7}

		if err := svc.AddVehicle(vehicle); err != nil {
			t.Fatalf("AddVehicle() error = %v", err)
		}
		if len(inspector.calls) != 1 || inspector.calls[0] != vehicle {
			t.Fatalf("Inspect() calls = %v, want [%v]", inspector.calls, vehicle)
		}
		if got := svc.VehicleCount(); got != 1 {
			t.Errorf("VehicleCount() = %d, want 1", got)
		}
	})

	t.Run("rejects failed inspection", func(t *testing.T) {
		inspector := &inspectorMock{result: false}
		svc := service.NewService(inspector)

		err := svc.AddVehicle(&vehicleMock{itemMock: itemMock{number: 1}})
		if !errors.Is(err, service.ErrInspectionFailed) {
			t.Fatalf("AddVehicle() error = %v, want %v", err, service.ErrInspectionFailed)
		}
		if len(svc.Inventory()) != 0 {
			t.Error("failed vehicle was added to inventory")
		}
	})

	t.Run("rejects duplicate before inspection", func(t *testing.T) {
		inspector := &inspectorMock{result: true}
		svc := service.NewService(inspector)
		vehicle := &vehicleMock{itemMock: itemMock{number: 1}}

		if err := svc.AddVehicle(vehicle); err != nil {
			t.Fatalf("first AddVehicle() error = %v", err)
		}
		err := svc.AddVehicle(&vehicleMock{itemMock: itemMock{number: 1}})
		if !errors.Is(err, service.ErrDuplicateNumber) {
			t.Fatalf("second AddVehicle() error = %v, want %v", err, service.ErrDuplicateNumber)
		}
		if got := len(inspector.calls); got != 1 {
			t.Errorf("Inspect() called %d times, want 1", got)
		}
	})
}

func TestAddItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(*service.Service)
		item    domain.InventoryItem
		wantErr error
	}{
		{
			name: "adds regular item",
			item: &itemMock{name: "helmet", number: 1},
		},
		{
			name:    "rejects vehicle",
			item:    &vehicleMock{itemMock: itemMock{name: "bicycle", number: 1}},
			wantErr: service.ErrVehicleAsItem,
		},
		{
			name: "rejects duplicate number",
			prepare: func(svc *service.Service) {
				if err := svc.AddItem(&itemMock{name: "helmet", number: 1}); err != nil {
					t.Fatalf("prepare AddItem() error = %v", err)
				}
			},
			item:    &itemMock{name: "lock", number: 1},
			wantErr: service.ErrDuplicateNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inspector := &inspectorMock{result: true}
			svc := service.NewService(inspector)
			if tt.prepare != nil {
				tt.prepare(svc)
			}

			err := svc.AddItem(tt.item)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddItem() error = %v, want %v", err, tt.wantErr)
			}
			if len(inspector.calls) != 0 {
				t.Errorf("AddItem() unexpectedly called Inspect()")
			}
		})
	}
}

func TestInventoryQueries(t *testing.T) {
	t.Parallel()

	inspector := &inspectorMock{result: true}
	svc := service.NewService(inspector)
	beginner := &vehicleMock{
		itemMock:  itemMock{name: "beginner bicycle", number: 1},
		easeOfUse: 6,
		energy:    2.5,
	}
	advanced := &vehicleMock{
		itemMock:  itemMock{name: "advanced bicycle", number: 2},
		easeOfUse: 5,
		energy:    1.5,
	}
	station := &energyItemMock{
		itemMock: itemMock{name: "charging station", number: 3},
		energy:   10,
	}

	for _, vehicle := range []domain.Vehicle{beginner, advanced} {
		if err := svc.AddVehicle(vehicle); err != nil {
			t.Fatalf("AddVehicle() error = %v", err)
		}
	}
	if err := svc.AddItem(station); err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	if got := svc.VehicleCount(); got != 2 {
		t.Errorf("VehicleCount() = %d, want 2", got)
	}
	if got := svc.TotalEnergyConsumption(); got != 14 {
		t.Errorf("TotalEnergyConsumption() = %v, want 14", got)
	}
	if got := svc.BeginnerVehicles(); len(got) != 1 || got[0] != beginner {
		t.Errorf("BeginnerVehicles() = %v, want [%v]", got, beginner)
	}

	inventory := svc.Inventory()
	inventory[0] = station
	if got := svc.Inventory()[0]; got != beginner {
		t.Error("Inventory() exposed the service's internal slice")
	}
}
