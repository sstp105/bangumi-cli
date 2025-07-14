package formathandler

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mediafmt"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

var (
	videoFormats                         = []string{".mp4", ".mkv", ".flac"}
	fmtter       mediafmt.MediaFormatter = mediafmt.TVShowFormatter{}

	episodeMap map[string]int
)

func Run() {
	loadEpisodes()

	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("获取当前工作目录路径错误: %v", err)
		return
	}

	traverse(wd)
}

func loadEpisodes() {
	episodeMap = make(map[string]int)

	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return
	}

	if subscription == nil {
		log.Error("subscription config file is empty")
	}
	
	for _, s := range subscription {
		var b model.Bangumi
		err := path.ReadJSONConfigFile(s.ConfigFileName(), &b)
		if err != nil {
			continue
		}
		episodeMap[b.Name] = b.StartEpisode()
	}
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
		return
	}

	// sort media files by file name in increasing order
	sort.Strings(files)

	rename(files, dir)
}

func rename(files []string, dir string) {
	log.Infof("%s, 共 %d 个文件", dir, len(files))

	subject := GetFolderName(dir)

	var offset int
	if v, ok := episodeMap[subject]; ok {
		offset = v
	} else {
		offset = 1
	}

	// dry-run
	paths, err := mediafmt.FormatFiles(files, dir, offset, fmtter)
	if err != nil {
		log.Errorf("命名出现错误:%s", err)
		return
	}

	if proceed := prompt.Confirm("是否要继续这些命名?"); !proceed {
		log.Warn("命名已取消")
		return
	}

	for i, f := range files {
		if err := os.Rename(f, paths[i]); err != nil {
			log.Errorf("命名 %s -> %s 时错误: %v", f, paths[i], err)
		}
	}
}

func GetFolderName(dir string) string {
	return filepath.Base(filepath.Clean(dir))
}