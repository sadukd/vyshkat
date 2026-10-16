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
	if got := svc.InventoryCount(); got != 3 {
		t.Errorf("InventoryCount() = %d, want 3", got)
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

func TestRemoveInventory(t *testing.T) {
	t.Parallel()

	newService := func(t *testing.T) *service.Service {
		t.Helper()
		svc := service.NewService(&inspectorMock{result: true})
		for _, item := range []domain.InventoryItem{
			&itemMock{name: "helmet", number: 1},
			&vehicleMock{itemMock: itemMock{name: "scooter", number: 2}, easeOfUse: 6, energy: 2.5},
			&energyItemMock{itemMock: itemMock{name: "station", number: 3}, energy: 10},
		} {
			var err error
			if vehicle, ok := item.(domain.Vehicle); ok {
				err = svc.AddVehicle(vehicle)
			} else {
				err = svc.AddItem(item)
			}
			if err != nil {
				t.Fatalf("prepare inventory: %v", err)
			}
		}
		return svc
	}

	t.Run("finds and removes items while preserving order", func(t *testing.T) {
		svc := newService(t)
		item, err := svc.FindInventoryItem(2)
		if err != nil || item.Number() != 2 {
			t.Fatalf("FindInventoryItem(2) = (%v, %v)", item, err)
		}
		if err := svc.RemoveVehicle(2); err != nil {
			t.Fatalf("RemoveVehicle(2) error = %v", err)
		}
		if err := svc.RemoveItem(1); err != nil {
			t.Fatalf("RemoveItem(1) error = %v", err)
		}
		inventory := svc.Inventory()
		if len(inventory) != 1 || inventory[0].Number() != 3 {
			t.Fatalf("Inventory() = %v, want only number 3", inventory)
		}
		if svc.VehicleCount() != 0 || svc.InventoryCount() != 1 || svc.TotalEnergyConsumption() != 10 {
			t.Error("inventory reports were not updated after removal")
		}
		if len(svc.BeginnerVehicles()) != 0 {
			t.Error("removed vehicle remains in beginner report")
		}
	})

	t.Run("rejects missing numbers and wrong categories", func(t *testing.T) {
		svc := newService(t)
		before := svc.Inventory()
		checks := []struct {
			name string
			err  error
			want error
		}{
			{name: "find missing", err: func() error { _, err := svc.FindInventoryItem(99); return err }(), want: service.ErrInventoryItemNotFound},
			{name: "remove missing", err: svc.RemoveItem(99), want: service.ErrInventoryItemNotFound},
			{name: "item as vehicle", err: svc.RemoveVehicle(1), want: service.ErrNotVehicle},
			{name: "vehicle as item", err: svc.RemoveItem(2), want: service.ErrVehicleRemovalRequired},
		}
		for _, check := range checks {
			if !errors.Is(check.err, check.want) {
				t.Errorf("%s error = %v, want %v", check.name, check.err, check.want)
			}
		}
		after := svc.Inventory()
		if len(after) != len(before) {
			t.Fatalf("failed removal changed inventory length to %d, want %d", len(after), len(before))
		}
		for i := range before {
			if after[i] != before[i] {
				t.Errorf("failed removal changed inventory at index %d", i)
			}
		}
	})

	t.Run("removes final item", func(t *testing.T) {
		svc := newService(t)
		if err := svc.RemoveItem(3); err != nil {
			t.Fatalf("RemoveItem(3) error = %v", err)
		}
		inventory := svc.Inventory()
		if len(inventory) != 2 || inventory[0].Number() != 1 || inventory[1].Number() != 2 {
			t.Fatalf("Inventory() numbers changed after removing final item: %v", inventory)
		}
	})
}
