package cmd

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/mediafmt"
)

var videoFormats []string = []string{".mp4", ".mkv", ".flac"}
var fmtter mediafmt.FileNameFormatter = mediafmt.TVShowFormatter{}

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename video files in the current and sub directories",
	Long: `The 'rename' command will rename video files in the current directory by appending the title, season, and episode number to each file name. 
	It will process files with common media formats (e.g., .mp4, .mkv) and recursively handle subdirectories. 
	Currently only supports Plex media content formats.
	The season and episode numbering will be determined based on the folder's name, and the files will be renamed accordingly.
	If the folder contains the season in custom format, e.g. 夏目友人帐 柒, which represents season 7, it's recommended to rename the folder to
	夏目友人帐 柒 - S07, 夏目友人帐 柒 第七季, or 夏目友人帐 第七期.
	Other examples: 东京喰种√A, CLANNAD 〜AFTER STORY〜, Code Geass 反叛的鲁路修R2.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		handler()
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func handler() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting current working directory:", err)
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

	// recursively process sub-directories
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
	if ok := libs.Confirm("Do you want to proceed with renaming these files?"); !ok {
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
