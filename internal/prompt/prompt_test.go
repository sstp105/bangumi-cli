package prompt

import (
	"os"
	"testing"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"User presses 'n'", "n", false},
		{"User presses 'y'", "y", true},
		{"User presses Enter", "", true},
		{"User presses random key", "random", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStdin(tt.input, func() {
				result := Confirm("Do you want to proceed?")
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			})
		})
	}
}

func TestReadUserInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"User enter text", "CHT", "CHT"},
		{"User enter comma split text", "简体,1080P", "简体,1080P"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStdin(tt.input, func() {
				result := ReadUserInput("Please enter text: (hit enter to proceed)")
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			})
		})
	}
}

// mockStdin replaces os.Stdin with provided input and restores it after test
func mockStdin(input string, testFunc func()) {
	stdin := os.Stdin
	// restore original os.Stdin after test
	defer func() {
		os.Stdin = stdin
	}()

	r, w, _ := os.Pipe()
	w.Write([]byte(input + "\n")) // Simulate user input
	w.Close()
	os.Stdin = r

	testFunc()
}
