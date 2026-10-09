package tests

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

func TestBicycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		bicycle   string
		id        int
		easeOfUse int
		condition int
		wantErr   error
	}{
		{name: "valid minimum condition", bicycle: "Road bicycle", id: 7, easeOfUse: 6, condition: 0},
		{name: "valid maximum condition", bicycle: "Road bicycle", id: 7, easeOfUse: 6, condition: 10},
		{name: "invalid ease of use", bicycle: "Road bicycle", id: 7, easeOfUse: 0, condition: 5, wantErr: domain.ErrInvalidEaseOfUse},
		{name: "condition above maximum", bicycle: "Road bicycle", id: 7, easeOfUse: 6, condition: 11, wantErr: domain.ErrInvalidCondition},
		{name: "negative item number", bicycle: "Road bicycle", id: -1, easeOfUse: 6, condition: 5, wantErr: domain.ErrNegativeItemNumber},
		{name: "blank name", bicycle: " \t\n", id: 7, easeOfUse: 6, condition: 5, wantErr: domain.ErrEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bicycle, err := domain.NewBicycle(tt.bicycle, tt.id, tt.easeOfUse, tt.condition)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewBicycle() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if bicycle != nil {
					t.Errorf("NewBicycle() = %#v, want nil", bicycle)
				}
				return
			}

			if bicycle.Name() != tt.bicycle || bicycle.Number() != tt.id ||
				bicycle.EaseOfUse() != tt.easeOfUse || bicycle.Condition() != tt.condition {
				t.Errorf("NewBicycle() did not preserve its input values")
			}
		})
	}
}
