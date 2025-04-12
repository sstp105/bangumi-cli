package season

import (
	"fmt"
	"testing"
)

func Test_Season(t *testing.T) {
	tests := []struct {
		id       ID
		expected Season
		wantErr  bool
	}{
		{1, Winter, false},
		{2, Spring, false},
		{3, Summer, false},
		{4, Autumn, false},
		{0, "", true},
		{5, "", true},
		{-1, "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("ID_%d", tt.id), func(t *testing.T) {
			season, err := tt.id.Season()

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if season != tt.expected {
				t.Errorf("expected season: %v, got: %v", tt.expected, season)
			}
		})
	}
}
