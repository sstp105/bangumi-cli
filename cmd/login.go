package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/loginhandler"
)

var loginCmd = &cobra.Command{
	Use: "login",
	Example: `
  bangumi login 授权 bangumi.tv。
`,
	Run: func(cmd *cobra.Command, args []string) {
		loginhandler.Run()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
