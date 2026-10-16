package app_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/samber/do/v2"

	"vyshkat/internal/app"
	"vyshkat/internal/console"
	"vyshkat/internal/service"
)

func TestContainerResolvesApplicationGraph(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	injector := app.NewContainer(strings.NewReader("q\n"), &output)

	application, err := do.Invoke[*console.Console](injector)
	if err != nil {
		t.Fatalf("resolve Console: %v", err)
	}
	if err := application.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(output.String(), "Inventory") {
		t.Error("resolved console did not write its menu to the injected output")
	}

	first, err := do.Invoke[*service.Service](injector)
	if err != nil {
		t.Fatalf("resolve Service: %v", err)
	}
	second, err := do.Invoke[*service.Service](injector)
	if err != nil {
		t.Fatalf("resolve Service again: %v", err)
	}
	if first != second {
		t.Error("container returned different service instances")
	}
}
