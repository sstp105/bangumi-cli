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
		log.Errorf("Error getting current working directory path: %v", err)
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
		log.Error("Subscription config file is empty")
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
		log.Errorf("Error reading path %s: %v", dir, err)
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
	log.Infof("%s, total %d files", dir, len(files))

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
		log.Errorf("Renaming error occurred: %s", err)
		return
	}

	if proceed := prompt.Confirm("Do you want to proceed with these names?"); !proceed {
		log.Warn("Renaming cancelled")
		return
	}

	for i, f := range files {
		if err := os.Rename(f, paths[i]); err != nil {
			log.Errorf("Error renaming %s -> %s: %v", f, paths[i], err)
		}
	}
}

func GetFolderName(dir string) string {
	return filepath.Base(filepath.Clean(dir))
}