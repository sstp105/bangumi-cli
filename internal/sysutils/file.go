package sysutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
)

// ErrEmptyFormats is returned when no file formats are provided.
var ErrEmptyFormats = errors.New("file formats cannot be empty")

// FindFiles searches for files in the specified directory that match any of the given formats.
//
// It scans the directory for files with the specified extensions and returns the first set of matching files found.
// If no files match, it returns an error.
//
// Parameters:
//   - dir: The directory path to search for files.
//   - formats: A list of file extensions to match (e.g., []string{".mp4", ".mkv"}).
//
// Returns:
//   - A slice of matching file paths if found.
//   - An error if the formats list is empty, if no matching files are found, or if there's an issue with filepath.Glob.
func FindFiles(dir string, formats []string) ([]string, error) {
	if len(formats) == 0 {
		return nil, ErrEmptyFormats
	}

	for _, ext := range formats {
		files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
		if err != nil {
			return nil, fmt.Errorf("error searching for %s files: %s", ext, err)
		}

		if len(files) > 0 {
			return files, nil
		}
	}

	return nil, fmt.Errorf("no files found with formats %v", formats)
}

// MarshalJSONIndented serializes the given data (v) into a pretty-printed JSON format.
//
// Parameters:
//   - v: Any Go data type that needs to be serialized into JSON.
//
// Returns:
//   - []byte: The indented JSON byte slice.
//   - error: An error if the marshaling process fails.
func MarshalJSONIndented(v any) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}
