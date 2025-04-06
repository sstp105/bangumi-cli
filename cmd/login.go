package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/login"
)

var loginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		login.Handler()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
