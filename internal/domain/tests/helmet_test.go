package tests

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

func TestHelmet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		helmet  string
		id      int
		wantErr error
	}{
		{name: "valid", helmet: "Helmet", id: 1},
		{name: "negative number", helmet: "Helmet", id: -1, wantErr: domain.ErrNegativeItemNumber},
		{name: "blank name", helmet: " \t", id: 1, wantErr: domain.ErrEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			helmet, err := domain.NewHelmet(tt.helmet, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewHelmet() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if helmet != nil {
					t.Errorf("NewHelmet() = %#v, want nil", helmet)
				}
				return
			}
			if helmet.Name() != tt.helmet || helmet.Number() != tt.id {
				t.Error("NewHelmet() did not preserve its input values")
			}
		})
	}
}
