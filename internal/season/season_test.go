package season

import (
	"testing"
	"time"
)

func Test_ID(t *testing.T) {
	tests := []struct {
		season   Season
		expected ID
	}{
		{Winter, 1},
		{Spring, 2},
		{Summer, 3},
		{Autumn, 4},
		{"Unknown", -1},
	}

	for _, tt := range tests {
		t.Run(tt.season.String(), func(t *testing.T) {
			id := tt.season.ID()
			if id != tt.expected {
				t.Errorf("expected ID: %d, got: %d", tt.expected, id)
			}
		})
	}
}

func Test_nowAt(t *testing.T) {
	tests := []struct {
		date     string
		expected Season
	}{
		{"2025-01-01", Winter},
		{"2025-04-10", Spring},
		{"2025-07-15", Summer},
		{"2025-11-01", Autumn},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			parsed, _ := time.Parse("2006-01-02", tt.date)
			got := nowAt(parsed)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
