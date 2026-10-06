package domain

import (
	"errors"
	"testing"
)

func TestNewScooter(t *testing.T) {
	t.Parallel()

	scooter, err := NewScooter("City scooter", 42, 1.25, 8)
	if err != nil {
		t.Fatalf("NewScooter() error = %v", err)
	}

	if got, want := scooter.Name(), "City scooter"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
	if got, want := scooter.Number(), 42; got != want {
		t.Errorf("Number() = %d, want %d", got, want)
	}
	if got, want := scooter.EnergyConsumption(), 1.25; got != want {
		t.Errorf("EnergyConsumption() = %v, want %v", got, want)
	}
	if got, want := scooter.EaseOfUse(), 8; got != want {
		t.Errorf("EaseOfUse() = %d, want %d", got, want)
	}
}

func TestNewScooterValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		scooter   string
		id        int
		energy    float64
		easeOfUse int
		wantErr   error
	}{
		{
			name:      "ease of use below minimum",
			scooter:   "City scooter",
			id:        1,
			energy:    1,
			easeOfUse: 0,
			wantErr:   ErrInvalidEaseOfUse,
		},
		{
			name:      "ease of use above maximum",
			scooter:   "City scooter",
			id:        1,
			energy:    1,
			easeOfUse: 11,
			wantErr:   ErrInvalidEaseOfUse,
		},
		{
			name:      "negative item number",
			scooter:   "City scooter",
			id:        -1,
			energy:    1,
			easeOfUse: 5,
			wantErr:   ErrNegativeItemNumber,
		},
		{
			name:      "negative energy consumption",
			scooter:   "City scooter",
			id:        1,
			energy:    -1,
			easeOfUse: 5,
			wantErr:   ErrInvalidEnergyConsumption,
		},
		{
			name:      "empty name",
			scooter:   " \t\n",
			id:        1,
			energy:    1,
			easeOfUse: 5,
			wantErr:   ErrEmptyName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scooter, err := NewScooter(
				tt.scooter,
				tt.id,
				tt.energy,
				tt.easeOfUse,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewScooter() error = %v, want %v", err, tt.wantErr)
			}
			if scooter != nil {
				t.Errorf("NewScooter() scooter = %#v, want nil", scooter)
			}
		})
	}
}

func TestNewScooterBoundaryValues(t *testing.T) {
	t.Parallel()

	for _, easeOfUse := range []int{1, 10} {
		scooter, err := NewScooter("City scooter", 0, 0, easeOfUse)
		if err != nil {
			t.Errorf("NewScooter() with easeOfUse %d returned error: %v", easeOfUse, err)
			continue
		}
		if got := scooter.EaseOfUse(); got != easeOfUse {
			t.Errorf("EaseOfUse() = %d, want %d", got, easeOfUse)
		}
	}
}
