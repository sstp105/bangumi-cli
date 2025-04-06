package format

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/mediafmt"
	"github.com/sstp105/bangumi-cli/internal/prompt"
	"os"
	"path/filepath"
	"sort"
)

var videoFormats []string = []string{".mp4", ".mkv", ".flac"}
var fmtter mediafmt.MediaFormatter = mediafmt.TVShowFormatter{}

func Run() {
	wd, err := os.Getwd()
	if err != nil {
		console.Errorf("获取当前工作目录路径错误: %v", err)
		return
	}

	traverse(wd)
}

func traverse(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("warning: skipping directory %s due to error: %v", dir, err)
		return
	}

	process(dir)

	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			traverse(subdir)
		}
	}
}

func process(dir string) {
	files, err := libs.FindFiles(dir, videoFormats)
	if err != nil {
		fmt.Printf("warning: no media files can be renamed in %s:%s\n", dir, err)
		return
	}

	// sort media files by file name in increasing order
	sort.Strings(files)

	rename(files, dir)
}

func rename(files []string, dir string) {
	fmt.Printf("Processing directory:%s\n", dir)

	// dry-run
	paths, err := mediafmt.FormatFiles(files, dir, fmtter)
	if err != nil {
		fmt.Printf("error formatting files at %s:%s", dir, err)
		return
	}

	// ask user before rename files
	if ok := prompt.Confirm("Do you want to proceed with renaming these files?"); !ok {
		fmt.Printf("Cancelled rename process for %s", dir)
		return
	}

	// rename each file
	for i, f := range files {
		if err := os.Rename(f, paths[i]); err != nil {
			fmt.Printf("error renaming %s -> %s: %v\n", f, paths[i], err)
		}
	}
}
