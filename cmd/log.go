package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/loghandler"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "查询日志",
	Long: `
Summary:
  log 命令用于显示调试日志。
	`,
	Example: `
  bangumi log 默认显示当天的日志。
`,
	Run: func(cmd *cobra.Command, args []string) {
		loghandler.Run()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
