package formathandler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTraverse(t *testing.T) {
	tmpDir := t.TempDir()

	root := filepath.Join(tmpDir, "测不准的阿波连同学")
	if err := os.Mkdir(root, 0755); err != nil {
		t.Fatalf("failed to create root folder: %v", err)
	}

	seasons := []string{"第一季", "第二季"}
	expectedFiles := make([]string, 0)

	for idx, season := range seasons {
		seasonDir := filepath.Join(root, fmt.Sprintf("测不准的阿波连同学 %s", season))
		if err := os.Mkdir(seasonDir, 0755); err != nil {
			t.Fatalf("failed to create season folder: %v", err)
		}

		seasonNum := fmt.Sprintf("S%02d", idx+1)
		for i := 1; i <= 12; i++ {
			originalName := fmt.Sprintf(
				"[喵萌奶茶屋] Aharen-san wa Hakarenai %s - %02d [WebRip 1080p HEVC-10bit AAC][简繁内封字幕].mp4",
				season, i,
			)
			filePath := filepath.Join(seasonDir, originalName)
			if err := writeVideoFile(filePath); err != nil {
				t.Fatalf("failed to create video file: %v", err)
			}

			expectedName := fmt.Sprintf("测不准的阿波连同学 - %sE%02d.mp4", seasonNum, i)
			expectedFiles = append(expectedFiles, filepath.Join(seasonDir, expectedName))
		}
	}

	mockStdin("y", func() {
		traverse(root)
	})

	for _, expected := range expectedFiles {
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("expected renamed file not found: %s", expected)
		}
	}
}

func mockStdin(input string, testFunc func()) {
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	r, w, _ := os.Pipe()
	w.Write([]byte(input + "\n"))
	w.Close()
	os.Stdin = r

	testFunc()
}

func writeVideoFile(path string) error {
	return os.WriteFile(path, []byte("dummy video content"), 0644)
}
