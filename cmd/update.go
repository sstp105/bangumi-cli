package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/handler/updatehandler"
	"github.com/sstp105/bangumi-cli/internal/log"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "同步订阅的 RSS 并更新下载任务",

	Long: `
Summary:
  update 命令会读取用户在蜜柑订阅的 RSS 链接，与本地已下载的种子进行比较，用户可选择是否添加新的种子任务到下载队列中。
`,
	Example: `
  bangumi update 查询 Mikan 订阅的 RSS 与本地进行对比并更新。
`,
	Run: func(cmd *cobra.Command, args []string) {
		h, err := updatehandler.NewHandler(config.MikanClientConfig())
		if err != nil {
			log.Errorf("updatehandler.NewHandler err: %v", err)
			return
		}
		h.Run()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
