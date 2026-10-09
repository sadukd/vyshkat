package tests

import (
	"testing"

	"vyshkat/internal/service"
)

func TestNumberGeneratorSequence(t *testing.T) {
	t.Parallel()

	generator := service.NewNumberGenerator()
	for want := 1; want <= 100; want++ {
		got, err := generator.Next()
		if err != nil {
			t.Fatalf("Next() call %d returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("Next() call %d = %d, want %d", want, got, want)
		}
	}
}

func TestNumberGeneratorsAreIndependent(t *testing.T) {
	t.Parallel()

	first := service.NewNumberGenerator()
	second := service.NewNumberGenerator()

	if _, err := first.Next(); err != nil {
		t.Fatalf("first generator Next() returned error: %v", err)
	}
	if _, err := first.Next(); err != nil {
		t.Fatalf("first generator Next() returned error: %v", err)
	}

	got, err := second.Next()
	if err != nil {
		t.Fatalf("second generator Next() returned error: %v", err)
	}
	if got != 1 {
		t.Errorf("second generator Next() = %d, want 1", got)
	}
}
