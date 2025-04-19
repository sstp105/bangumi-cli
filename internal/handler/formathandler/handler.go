package formathandler

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mediafmt"
	"github.com/sstp105/bangumi-cli/internal/prompt"
	"os"
	"path/filepath"
	"sort"
)

var (
	videoFormats                         = []string{".mp4", ".mkv", ".flac"}
	fmtter       mediafmt.MediaFormatter = mediafmt.TVShowFormatter{}
)

func Run() {
	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("获取当前工作目录路径错误: %v", err)
		return
	}

	traverse(wd)
}

func traverse(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Errorf("读取路径 %s 时错误:%v", dir, err)
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
		log.Errorf("%s 查找文件错误: %v", dir, err)
		return
	}

	// sort media files by file name in increasing order
	sort.Strings(files)

	rename(files, dir)
}

func rename(files []string, dir string) {
	log.Debugf("处理:%s, 共 %d 个文件", dir, len(files))

	// dry-run
	paths, err := mediafmt.FormatFiles(files, dir, fmtter)
	if err != nil {
		log.Errorf("命名时出现错误, 路径: %s:%s", dir, err)
		return
	}

	if proceed := prompt.Confirm("是否要继续这些命名?"); !proceed {
		log.Info("命名已取消")
		return
	}

	for i, f := range files {
		if err := os.Rename(f, paths[i]); err != nil {
			log.Errorf("命名 %s -> %s 时错误: %v\n", f, paths[i], err)
		}
	}
}
