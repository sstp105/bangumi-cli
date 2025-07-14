package mediafmt

import (
	"path/filepath"
	"testing"
)

func TestFormatFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		dir      string
		expected []string
		wantErr  bool
	}{
		{
			name:  "Test formatting multiple files",
			files: []string{"青之箱 (1).mp4", "青之箱 (2).mp4", "青之箱 (3).mp4"},
			dir:   "/青之箱",
			expected: []string{
				"/青之箱/青之箱 - S01E01.mp4",
				"/青之箱/青之箱 - S01E02.mp4",
				"/青之箱/青之箱 - S01E03.mp4",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmtter := TVShowFormatter{}

			got, err := fmtter.FormatFiles(tt.files, tt.dir, 1)

			// normalize path separators to forward slashes for comparison (windows \, unix /)
			for i, v := range got {
				got[i] = filepath.ToSlash(v)
			}

			expected := tt.expected
			for i, v := range tt.expected {
				expected[i] = filepath.ToSlash(v)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("FormatFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalStringSlices(got, expected) {
				t.Errorf("FormatFiles() = %v, want %v", got, expected)
			}
		})
	}
}

// Test case for the FormatFileName method
func TestFormatFileName(t *testing.T) {
	tests := []struct {
		name     string
		f        string
		dir      string
		metadata TVShowMetadata
		expected string
		wantErr  bool
	}{
		{
			name: "Test valid metadata with all fields",
			f:    "[Aharen-san wa Hakarenai][02][BDRip 1080p AVC AAC][CHS].mkv",
			dir:  "/tv_shows",
			metadata: TVShowMetadata{
				Title:        strPtr("测不准的阿波连同学"),
				Season:       intPtr(1),
				Episode:      intPtr(2),
				Year:         strPtr("2022"),
				EpisodeTitle: strPtr("我是被跟踪了吧？"),
			},
			expected: "/tv_shows/测不准的阿波连同学 (2022) - S01E02 - 我是被跟踪了吧？.mkv",
			wantErr:  false,
		},
		{
			name: "Test valid metadata without episode title",
			f:    "[Aharen-san wa Hakarenai][02][BDRip 1080p AVC AAC][CHS].mkv",
			dir:  "/tv_shows",
			metadata: TVShowMetadata{
				Title:   strPtr("测不准的阿波连同学"),
				Season:  intPtr(1),
				Episode: intPtr(2),
				Year:    strPtr("2022"),
			},
			expected: "/tv_shows/测不准的阿波连同学 (2022) - S01E02.mkv",
			wantErr:  false,
		},
		{
			name: "Test missing title in metadata",
			f:    "[Aharen-san wa Hakarenai][02][BDRip 1080p AVC AAC][CHS].mkv",
			dir:  "/tv_shows",
			metadata: TVShowMetadata{
				Season:  intPtr(1),
				Episode: intPtr(2),
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tvFormatter := TVShowFormatter{}
			got, err := tvFormatter.FormatFileName(tt.f, tt.dir, tt.metadata)

			// normalize path separators to forward slashes for comparison (windows \, unix /)
			got = filepath.ToSlash(got)
			expected := filepath.ToSlash(tt.expected)

			if (err != nil) != tt.wantErr {
				t.Errorf("FormatFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != expected {
				t.Errorf("FormatFileName() = %v, want %v", got, expected)
			}
		})
	}
}

func TestParseFolderName(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantTitle  string
		wantSeason int
	}{
		{
			name:       "No season - default season",
			input:      "好想告诉你",
			wantTitle:  "好想告诉你",
			wantSeason: 1,
		},
		{
			name:       "Chinese season format - full character",
			input:      "好想告诉你 第三季",
			wantTitle:  "好想告诉你",
			wantSeason: 3,
		},
		{
			name:       "Chinese season format - numeric",
			input:      "好想告诉你 第3季",
			wantTitle:  "好想告诉你",
			wantSeason: 3,
		},
		{
			name:       "Chinese period format - full character",
			input:      "好想告诉你 第三期",
			wantTitle:  "好想告诉你",
			wantSeason: 3,
		},
		{
			name:       "Chinese period format - numeric",
			input:      "好想告诉你 第3期",
			wantTitle:  "好想告诉你",
			wantSeason: 3,
		},
		{
			name:       "Chinese season > 10 - full character",
			input:      "好想告诉你 第十五季",
			wantTitle:  "好想告诉你",
			wantSeason: 15,
		},
		{
			name:       "Chinese period > 10 - full character",
			input:      "好想告诉你 第十五期",
			wantTitle:  "好想告诉你",
			wantSeason: 15,
		},
		{
			name:       "Title with subtitle before season",
			input:      "我独自升级 第二季 -起于暗影-",
			wantTitle:  "我独自升级",
			wantSeason: 2,
		},
		{
			name:       "Title with subtitle after season",
			input:      "我独自升级 -起于暗影- 第二季",
			wantTitle:  "我独自升级 -起于暗影-",
			wantSeason: 2,
		},
		{
			name:       "English season format",
			input:      "我独自升级 S02",
			wantTitle:  "我独自升级",
			wantSeason: 2,
		},
		{
			name:       "English season format with subtitle",
			input:      "我独自升级 -起于暗影- S03",
			wantTitle:  "我独自升级 -起于暗影-",
			wantSeason: 3,
		},
		{
			name:       "English season format double-digit",
			input:      "我独自升级 -起于暗影- S13",
			wantTitle:  "我独自升级 -起于暗影-",
			wantSeason: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotSeason := parseFolderName(tt.input)
			if gotTitle != tt.wantTitle || gotSeason != tt.wantSeason {
				t.Errorf("parseFolderName(%q) = (%q, %q), want (%q, %q)",
					tt.input, gotTitle, gotSeason, tt.wantTitle, tt.wantSeason)
			}
		})
	}
}

// Helpers
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
