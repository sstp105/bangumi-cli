package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/unsubscribehandler"
)

var id int

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "取消订阅番剧",
	Long: `
Summary:
  unsubscribe 命令用于取消订阅本地订阅的番剧。
`,
	Example: `
  bangumi unsubscribe 取消所有本地订阅。
  bangumi unsubscribe --id 3513 取消指定番剧的订阅。
`,
	Run: func(cmd *cobra.Command, args []string) {
		unsubscribehandler.Run(id)
	},
}

func init() {
	unsubscribeCmd.Flags().IntVarP(&id, "id", "i", -1, "取消订阅的番剧 ID")

	rootCmd.AddCommand(unsubscribeCmd)
}
