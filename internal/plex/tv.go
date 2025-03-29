package plex

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	cn "github.com/aiialzy/chinese-number"
)

type TVShowFormatter struct{}

var tvfmtter = TVShowFormatter{}

// FormatFiles formats/renames a list of TV shows based on the information derived from the directory name.
//
// Parameters:
//
//	files ([]string): A list of file names to be formatted.
//	dir (string): The directory path that contains the files to be formatted. 
//
// Returns: A slice of renamed file paths and an error (if any).
//
//	If renaming fails, the process will be aborted and an error will be returned.
func (TVShowFormatter) FormatFiles(files []string, dir string) ([]string, error) {

	basename := filepath.Base(dir)
	title, season := parseFolderName(basename)
	fmt.Printf("Parsed title:%s, season:%d from %s", title, season, basename)

	paths := make([]string, len(files))

	fmt.Println("The following files will be renamed:")
	for i, f := range files {
		episode := i + 1 // start episode with 1, 0 normally represents speical (特别篇).
		meta := TVShowMetadata{
			Title:   &title,
			Season:  &season,
			Episode: &episode,
		}

		fn, err := FormatFileName(f, dir, meta, tvfmtter)
		if err != nil {
			fmt.Printf("error formatting file name:%s", err)
			return nil, err
		}

		paths[i] = fn
		fmt.Printf("%s -> %s\n", f, paths[i])
	}

	return paths, nil
}

// FormatFileName generates a formatted file name for a TV show based on the provided metadata.
//
// Parameters:
//   f (string): The original file name.
//   dir (string): The directory path of the file.
//   metadata (interface{}): The metadata for the TV show, expected to be of type TVShowMetadata.
//
// Returns:
//   string: The formatted file name, including the directory path.
//   error: An error if the metadata is invalid or the validation fails.
func (TVShowFormatter) FormatFileName(f, dir string, metadata interface{}) (string, error) {
	data, ok := metadata.(TVShowMetadata)
	if !ok {
		return "", ErrInvalidTVShowMetadata
	}

	if err := data.validate(); err != nil {
		return "", err
	}

	var fn string

	fn = *data.Title
	if data.Year != nil {
		fn = fmt.Sprintf("%s (%s)", fn, *data.Year)
	}

	fn += fmt.Sprintf(" - S%02dE%02d", *data.Season, *data.Episode)

	if data.EpisodeTitle != nil {
		fn += fmt.Sprintf(" - %s", *data.EpisodeTitle)
	}

	ext := filepath.Ext(f)
	fn += ext

	return filepath.Join(dir, fn), nil
}

// parseFolderName extracts the title and seasonfrom a given folder name (basename).
//
// The function supports the following season formats:
//   - "第N季", where N can be arabic or chinese number. (e.g., "第1季", "第十季")
//   - "第N期", where N can be arabic or chinese number. (e.g., "第2期", "第十期")
//   - "SNN", represents season in 2 digit format. (e.g., "S01", "S10")
func parseFolderName(basename string) (string, int) {

	patterns := []string{
		`第([\d一二三四五六七八九十]+)季`,
		`第([\d一二三四五六七八九十]+)期`,
		`第(\d+)季`,
		`第(\d+)期`,
		`S(\d{2})\b`,
	}

	for _, p := range patterns {
		reg := regexp.MustCompile(p)
		match := reg.FindStringSubmatch(basename)

		if match == nil {
			continue
		}

		title := strings.TrimSpace(basename[:strings.Index(basename, match[0])])

		season := match[1]
		n64, err := cn.Parse(season)
		if err == nil {
			return title, int(n64)
		}

		n, err := strconv.Atoi(season)
		if err == nil {
			return title, n
		}
	}

	// use folder name as title, and default season 1
	return basename, 1
}
