package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/downloadhandler"
	"github.com/sstp105/bangumi-cli/internal/path"
)

var output string

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "批量处理 RSS，并将下载任务发送到 qBitTorrent 队列",
	Long: `
Summary:
  download 会加载已订阅的番剧种子任务，并将其添加到 qBitTorrent 进行下载。可选参数允许指定种子的保存目录。`,
	Example: `
  bangumi download           下载到当前工作目录
`,
	Run: func(cmd *cobra.Command, args []string) {
		h, err := downloadhandler.NewHandler(output)
		if err != nil {
			return
		}

		_ = h.Run()
	},
}

func init() {
	defaultOutputPath, err := path.DownloadPath()
	if err != nil {
		os.Exit(1)
	}
	downloadCmd.Flags().StringVarP(&output, "output", "o", defaultOutputPath, "保存路径")

	rootCmd.AddCommand(downloadCmd)
}
