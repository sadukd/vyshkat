package tests

import (
	"errors"
	"math"
	"testing"

	"vyshkat/internal/domain"
)

func TestChargingStation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		station string
		id      int
		energy  float64
		wantErr error
	}{
		{name: "valid", station: "Station", id: 1, energy: 4.5},
		{name: "negative number", station: "Station", id: -1, energy: 4.5, wantErr: domain.ErrNegativeItemNumber},
		{name: "blank name", station: " \t", id: 1, energy: 4.5, wantErr: domain.ErrEmptyName},
		{name: "negative energy", station: "Station", id: 1, energy: -1, wantErr: domain.ErrInvalidEnergyConsumption},
		{name: "NaN energy", station: "Station", id: 1, energy: math.NaN(), wantErr: domain.ErrInvalidEnergyConsumption},
		{name: "infinite energy", station: "Station", id: 1, energy: math.Inf(1), wantErr: domain.ErrInvalidEnergyConsumption},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			station, err := domain.NewChargingStation(tt.station, tt.id, tt.energy)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewChargingStation() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if station != nil {
					t.Errorf("NewChargingStation() = %#v, want nil", station)
				}
				return
			}
			if station.Name() != tt.station || station.Number() != tt.id || station.EnergyConsumption() != tt.energy {
				t.Error("NewChargingStation() did not preserve its input values")
			}
		})
	}
}
