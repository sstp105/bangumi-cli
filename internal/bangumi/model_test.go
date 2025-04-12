package bangumi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubjectCollectionType_IsValid(t *testing.T) {
	tests := []struct {
		input    SubjectCollectionType
		expected bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{0, false},
		{6, false},
		{-1, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.input.IsValid(), "input: %d", tt.input)
	}
}

func TestSubjectCollectionType_String(t *testing.T) {
	tests := []struct {
		input    SubjectCollectionType
		expected string
	}{
		{1, "想看"},
		{2, "看过"},
		{3, "在看"},
		{4, "搁置"},
		{5, "抛弃"},
		{0, ""},
		{6, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.input.String(), "input: %d", tt.input)
	}
}
