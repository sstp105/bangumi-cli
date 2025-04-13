package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/log"
)

var rootCmd = &cobra.Command{
	Use: "bangumi-cli",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error executing root command: %v", err)
	}
}
