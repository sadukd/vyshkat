package tests

import (
	"testing"

	"vyshkat/internal/inspection"
)

type vehicleMock struct {
	condition int
}

func (v vehicleMock) Name() string { return "vehicle" }

func (v vehicleMock) Number() int { return 1 }

func (v vehicleMock) EaseOfUse() int { return 1 }

func (v vehicleMock) Condition() int { return v.condition }

func TestInspectorInspect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		condition int
		want      bool
	}{
		{name: "minimum condition", condition: 0, want: false},
		{name: "below acceptance threshold", condition: inspection.MinAcceptedCondition - 1, want: false},
		{name: "at acceptance threshold", condition: inspection.MinAcceptedCondition, want: true},
		{name: "maximum condition", condition: 10, want: true},
	}

	inspector := inspection.NewInspector()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := inspector.Inspect(vehicleMock{condition: tt.condition}); got != tt.want {
				t.Errorf("Inspect() = %t for condition %d, want %t", got, tt.condition, tt.want)
			}
		})
	}
}
