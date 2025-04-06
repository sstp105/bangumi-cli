package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/command/subscribe"
	"github.com/sstp105/bangumi-cli/internal/season"
	"time"
)

var year int
var seasonID int

/*
subscribeCmd will perform the following tasks:

If no subscribed bangumi config is found locally, it fetches the user's subscribed bangumi list from Mikan, parses and saves the config.
For each bangumi, it requests the corresponding Mikan bangumi page, parses the bangumi.tv ID and the user-subscribed fan-sub RSS link,
prompts the user to filter desired torrent files, and stores the config locally for further processing.

If a subscribed bangumi config already exists, it fetches the latest subscribed list from Mikan,
compares it with the local config, and prompts the user to add any new bangumi.
If a bangumi is no longer present on Mikan, it prompts the user to confirm whether to unsubscribe it locally.
*/
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "订阅你在 Mikan 上关注的番剧",
	Long: `
Summary:
  subscribe 命令用于处理用户对番剧的订阅。如果本地未找到订阅番剧配置，它将从 Mikan 获取用户订阅的番剧列表，解析并保存配置。
  对于每一部番剧，它会请求对应的 Mikan 番剧页面，解析 bangumi.tv 的 ID 和用户订阅的字幕组 RSS 链接，提示用户筛选想要的种子文件，并将配置保存在本地以供后续处理。
  如果本地已有订阅配置，则会从 Mikan 获取最新的订阅列表，与本地配置进行对比，并提示用户是否添加新增的番剧。
  如果某部番剧已从 Mikan 上移除，则会提示用户是否取消订阅该番剧。
`,
	Example: `
  bangumi subscribe 默认将使用当前年份和季度。
  bangumi subscribe --season 2 --year 2025 指定读取 2025 年，春季番剧订阅列表。
`,
	Run: func(cmd *cobra.Command, args []string) {
		subscribe.Run(year, season.ID(seasonID))
	},
}

func init() {
	defaultYear := time.Now().Year()
	defaultSeasonID := season.Now().ID()

	subscribeCmd.Flags().IntVarP(&year, "year", "y", defaultYear, "选择指定年的番剧订阅列表")
	subscribeCmd.Flags().IntVarP(
		&seasonID,
		"season",
		"s", int(defaultSeasonID),
		"选择指定季度的番剧订阅列表, 可选值为：1, 2, 3, 4，分别对应 冬,春,夏,秋")

	rootCmd.AddCommand(subscribeCmd)
}
