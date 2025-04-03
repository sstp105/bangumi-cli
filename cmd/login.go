package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/login"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "验证用户 Bangumi 凭证并授权获取 Bangumi API 访问令牌。",
	Run: func(cmd *cobra.Command, args []string) {
		login.Handler()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
