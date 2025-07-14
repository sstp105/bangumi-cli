package mediafmt

type MediaFormatter interface {
	// FormatFileName formats a single file name using the provided metadata and directory.
	FormatFileName(f, dir string, metadata interface{}) (string, error)

	// FormatFiles formats a list of files in a given directory and returns a slice of formatted file names.
	FormatFiles(files []string, dir string, soffset int) ([]string, error)
}

func FormatFileName(f, dir string, metadata interface{}, fmtter MediaFormatter) (string, error) {
	return fmtter.FormatFileName(f, dir, metadata)
}

func FormatFiles(files []string, dir string, offset int, fmtter MediaFormatter) ([]string, error) {
	return fmtter.FormatFiles(files, dir, offset)
}
