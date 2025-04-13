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
  unsubscribe 命令用于
`,
	Example: `
  bangumi unsubscribe 
  bangumi unsubscribe --id 3513
`,
	Run: func(cmd *cobra.Command, args []string) {
		unsubscribehandler.Run(id)
	},
}

func init() {
	unsubscribeCmd.Flags().IntVarP(&id, "id", "i", -1, "指定取消订阅的番剧 ID")

	rootCmd.AddCommand(unsubscribeCmd)
}
