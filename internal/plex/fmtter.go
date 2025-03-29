package plex

type FileNameFormatter interface {
	FormatFileName(f, dir string, metadata interface{}) (string, error)
	FormatFiles(files []string, dir string) ([]string, error)
}

func FormatFileName(f, dir string, metadata interface{}, fmtter FileNameFormatter) (string, error) {
	return fmtter.FormatFileName(f, dir, metadata)
}

func FormatFiles(files []string, dir string, fmtter FileNameFormatter) ([]string, error) {
	return fmtter.FormatFiles(files, dir)
}
