package console

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestRunMenuPaginatesAndDispatches(t *testing.T) {
	t.Parallel()

	selected := -1
	options := make([]menuOption, 12)
	for i := range options {
		index := i
		options[i] = menuOption{
			label: fmt.Sprintf("Option %d", i+1),
			action: func() error {
				selected = index
				return nil
			},
		}
	}

	var output bytes.Buffer
	app := &Console{
		input:  bufioScanner(">\n2\n<\n<\n>\n>\ninvalid\nq\n"),
		output: &output,
	}
	if err := app.runMenu("Many options", options, "Back"); err != nil {
		t.Fatalf("runMenu() error = %v", err)
	}

	if selected != 11 {
		t.Errorf("selected option = %d, want 11", selected)
	}
	text := output.String()
	for _, want := range []string{
		"≡≡≡≡≡≡ Many options ≡≡≡≡≡≡",
		"0. Option 10",
		"1. Option 11",
		"2. Option 12",
		strings.Repeat("≡", 26) + "\n< previous",
		"< previous | > next | q back",
		"No previous page.",
		"No next page.",
		"Invalid choice.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("output does not contain %q", want)
		}
	}
}

func TestRunMenuReturnsActionError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("action failed")
	app := &Console{
		input:  bufioScanner("1\n"),
		output: &bytes.Buffer{},
	}
	err := app.runMenu("Menu", []menuOption{{
		label: "Fail",
		action: func() error {
			return wantErr
		},
	}}, "Back")
	if !errors.Is(err, wantErr) {
		t.Fatalf("runMenu() error = %v, want %v", err, wantErr)
	}
}

func bufioScanner(input string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(input))
}
