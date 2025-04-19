package libs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSuffixID(t *testing.T) {
	tests := []struct {
		name        string
		href        string
		expectedID  string
		expectError bool
	}{
		{
			name:        "Valid href",
			href:        "/subject/127791",
			expectedID:  "127791",
			expectError: false,
		},
		{
			name:        "Empty href",
			href:        "",
			expectedID:  "",
			expectError: true,
		},
		{
			name:        "Invalid href - not enough parts",
			href:        "/123",
			expectedID:  "123",
			expectError: false,
		},
		{
			name:        "Valid href with extra slashes",
			href:        "/collection/subject/123456",
			expectedID:  "123456", // The ID is the last part of the URL
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseSuffixID(tt.href)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
