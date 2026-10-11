package tests

import (
	"testing"

	"vyshkat/internal/domain"
)

func TestChargingStationPreservesEnergyConsumption(t *testing.T) {
	t.Parallel()

	station, err := domain.NewChargingStation("Station", 1, 4.5)
	if err != nil {
		t.Fatalf("NewChargingStation() error = %v", err)
	}
	if got := station.EnergyConsumption(); got != 4.5 {
		t.Errorf("EnergyConsumption() = %v, want 4.5", got)
	}
}
