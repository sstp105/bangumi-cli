package mediafmt

// MediaFormatter defines methods for formatting media file names and directories.
type MediaFormatter interface {
	// FormatFileName formats a single file name using the provided metadata and directory.
	FormatFileName(f, dir string, metadata interface{}) (string, error)

	// FormatFiles formats a list of files in a given directory.
	// It returns a slice of formatted file names.
	FormatFiles(files []string, dir string) ([]string, error)
}

// FormatFileName uses the provided MediaFormatter to format a single file name.
func FormatFileName(f, dir string, metadata interface{}, fmtter MediaFormatter) (string, error) {
	return fmtter.FormatFileName(f, dir, metadata)
}

// FormatFiles uses the provided MediaFormatter to format a list of file names.
func FormatFiles(files []string, dir string, fmtter MediaFormatter) ([]string, error) {
	return fmtter.FormatFiles(files, dir)
}
