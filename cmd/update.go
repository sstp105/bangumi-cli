package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/handler/updatehandler"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Example: `
  bangumi update。
`,
	Run: func(cmd *cobra.Command, args []string) {
		h, err := updatehandler.NewHandler()
		if err != nil {
			console.Errorf("updatehandler.NewHandler err: %v", err)
			return
		}
		h.Run()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
