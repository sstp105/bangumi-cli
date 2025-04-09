package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/update"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Example: `
  bangumi update。
`,
	Run: func(cmd *cobra.Command, args []string) {
		update.Run()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
