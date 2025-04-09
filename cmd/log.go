package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/loghandler"
)

var logCmd = &cobra.Command{
	Use: "log",
	Example: `
  bangumi log。
`,
	Run: func(cmd *cobra.Command, args []string) {
		loghandler.Run()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
