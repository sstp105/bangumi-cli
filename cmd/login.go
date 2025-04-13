package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/loginhandler"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录 bangumi.tv 并获取 API 访问令牌",
	Long: `
Summary:
  login 命令用于登录 bangumi.tv 用户，并获取用于访问 bangumi API 的令牌。
  该令牌的有效期为 7 天。在有效期内，login 命令会自动使用 refresh token 刷新 access token。
	`,
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
