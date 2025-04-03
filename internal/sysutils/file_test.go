package sysutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFindFiles(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "[Title][01][BDRip 1080p AVC AAC][CHS].mp4")
	file2 := filepath.Join(tempDir, "[Title][02][BDRip 1080p AVC AAC][CHS].mp4")
	file3 := filepath.Join(tempDir, "[Title][OVA1][BDRip 1080p AVC AAC][CHS].mkv")
	_ = os.WriteFile(file1, []byte(""), 0644)
	_ = os.WriteFile(file2, []byte(""), 0644)
	_ = os.WriteFile(file3, []byte(""), 0644)

	tests := []struct {
		name     string
		dir      string
		formats  []string
		wantErr  error
		expected []string
	}{
		{
			name:     "Valid case - Find .mp4 files",
			dir:      tempDir,
			formats:  []string{".mp4"},
			wantErr:  nil,
			expected: []string{file1, file2},
		},
		{
			name:     "Valid case - Find .mkv files",
			dir:      tempDir,
			formats:  []string{".mkv"},
			wantErr:  nil,
			expected: []string{file3},
		},
		{
			name:     "No matching files",
			dir:      tempDir,
			formats:  []string{".flac"},
			wantErr:  fmt.Errorf("no files found with formats %v", []string{".flac"}),
			expected: nil,
		},
		{
			name:     "Empty formats list",
			dir:      tempDir,
			formats:  []string{},
			wantErr:  ErrEmptyFormats,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := FindFiles(tt.dir, tt.formats)

			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				}
			} else if err == nil && tt.wantErr != nil {
				t.Errorf("expected error: %v, got nil", tt.wantErr)
			}

			if len(files) != len(tt.expected) {
				t.Errorf("expected files: %v, got: %v", tt.expected, files)
			}
		})
	}
}
