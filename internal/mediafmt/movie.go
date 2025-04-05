package mediafmt

import (
	"fmt"
	"path/filepath"
)

type MovieFormatter struct{}

// TODO
var _ = MovieFormatter{}

// FormatFiles formats/renames a list of movies based on the information derived from the directory name.
//
// Parameters:
//
//	files ([]string): A list of file names to be formatted.
//	dir (string): The directory path that contains the files to be formatted.
//
// Returns: A slice of renamed file paths and an error (if any).
//
//	If renaming fails, the process will be aborted and an error will be returned.
func (MovieFormatter) FormatFiles(files []string, dir string) ([]string, error) {
	// TODO
	return nil, nil
}

// FormatFileName generates a formatted file name for a movie based on the provided metadata.
//
// Parameters:
//
//	f (string): The original file name.
//	dir (string): The directory path of the file.
//	metadata (interface{}): The metadata for the movie, expected to be of type MovieMetadata.
//
// Returns:
//
//	string: The formatted file name, including the directory path.
//	error: An error if the metadata is invalid or the validation fails.
func (MovieFormatter) FormatFileName(f, dir string, metadata interface{}) (string, error) {
	data, ok := metadata.(MovieMetadata)
	if !ok {
		return "", ErrInvalidMovieMetadata
	}

	if err := data.validate(); err != nil {
		return "", err
	}

	var fn string
	fn = fmt.Sprintf("%s (%s)", *data.Title, *data.Year)

	ext := filepath.Ext(f)
	fn += ext

	return filepath.Join(dir, fn), nil
}
