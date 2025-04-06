package cmd

import (
	"github.com/spf13/cobra"
)

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "取消订阅番剧",
	Long: `
Summary:
  unsubscribe 命令用于
`,
	Example: `
  bangumi unsubscribe 
  bangumi unsubscribe --id 3513
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(unsubscribeCmd)
}
