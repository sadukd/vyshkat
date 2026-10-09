package tests

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

func TestScooter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		scooter   string
		id        int
		energy    float64
		easeOfUse int
		condition int
		wantErr   error
	}{
		{name: "valid minimum condition", scooter: "City scooter", id: 42, energy: 1.25, easeOfUse: 8, condition: 0},
		{name: "valid maximum condition", scooter: "City scooter", id: 42, energy: 1.25, easeOfUse: 8, condition: 10},
		{name: "invalid ease of use", scooter: "City scooter", id: 42, energy: 1.25, easeOfUse: 11, condition: 5, wantErr: domain.ErrInvalidEaseOfUse},
		{name: "condition below minimum", scooter: "City scooter", id: 42, energy: 1.25, easeOfUse: 8, condition: -1, wantErr: domain.ErrInvalidCondition},
		{name: "negative item number", scooter: "City scooter", id: -1, energy: 1.25, easeOfUse: 8, condition: 5, wantErr: domain.ErrNegativeItemNumber},
		{name: "negative energy consumption", scooter: "City scooter", id: 42, energy: -1, easeOfUse: 8, condition: 5, wantErr: domain.ErrInvalidEnergyConsumption},
		{name: "blank name", scooter: " \t\n", id: 42, energy: 1.25, easeOfUse: 8, condition: 5, wantErr: domain.ErrEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scooter, err := domain.NewScooter(tt.scooter, tt.id, tt.energy, tt.easeOfUse, tt.condition)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewScooter() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if scooter != nil {
					t.Errorf("NewScooter() = %#v, want nil", scooter)
				}
				return
			}

			if scooter.Name() != tt.scooter || scooter.Number() != tt.id ||
				scooter.EnergyConsumption() != tt.energy || scooter.EaseOfUse() != tt.easeOfUse ||
				scooter.Condition() != tt.condition {
				t.Errorf("NewScooter() did not preserve its input values")
			}
		})
	}
}
