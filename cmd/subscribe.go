package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/command/subscribe"
)

var subscribeCmd = &cobra.Command{
	Use: "subscribe",
	Run: func(cmd *cobra.Command, args []string) {
		subscribe.Handler()
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
