package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/downloadhandler"
)

var downloadCmd = &cobra.Command{
	Use: "download",
	Example: `
  bangumi download 默认下载到当前工作目录。
  bangumi download -o /path 下载到指定的目录。
`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadhandler.Run()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
