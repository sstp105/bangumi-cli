package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/loghandler"
)

var date string

var logCmd = &cobra.Command{
	Use: "log",
	Example: `
  bangumi log
  bangumi log -d 2025-04-10
`,
	Run: func(cmd *cobra.Command, args []string) {
		loghandler.Run()
	},
}

func init() {
	logCmd.Flags().StringVarP(&date, "date", "d", "", "查询指定日期的日志, 默认将显示当天的日志。")

	rootCmd.AddCommand(logCmd)
}
