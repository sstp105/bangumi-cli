package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/auth"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate user and obtain banggumi API access token.",
	Run: func(cmd *cobra.Command, args []string) {
		auth.Handler()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
