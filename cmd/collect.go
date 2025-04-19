package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/handler/collecthandler"
	"github.com/sstp105/bangumi-cli/internal/log"
)

var username string
var collectionType int

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "批量收藏bangumi.tv条目、 同步mikan订阅列表",
	Long: `
Summary:
  collect 命令用于批量收藏用户在 Mikan 上订阅的番剧，并将其同步到 Bangumi.tv，设置为指定的收藏状态（例如“在看”、“想看”等）。
  由于该任务需要用户 bangumi.tv 的授权才能更新条目收藏状态, 请先运行 login 命令。
	`,
	Example: `
  bangumi collect -u username -t 3 将批量同步 mikan 订阅的番剧到bangumi.tv 用户 username 在看状态。
`,
	Run: func(cmd *cobra.Command, args []string) {
		t := bangumi.SubjectCollectionType(collectionType)

		h, err := collecthandler.NewHandler(username, t)
		if err != nil {
			log.Fatalf("failed to init collect handler: %v", err)
		}

		h.Run()
	},
}

func init() {
	collectCmd.Flags().StringVarP(&username, "username", "u", "", "bangumi.tv 用户名")
	collectCmd.Flags().IntVarP(&collectionType, "type", "t", -1,
		"1 = 想看 (wish)\n"+
			"2 = 看过 (collect)\n"+
			"3 = 在看 (do)\n"+
			"4 = 搁置 (on_hold)\n"+
			"5 = 抛弃 (dropped)")

	rootCmd.AddCommand(collectCmd)
}
