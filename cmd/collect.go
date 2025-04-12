package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/handler/collecthandler"
	"github.com/sstp105/bangumi-cli/internal/log"
)

var username string
var collectionType int

var collectCmd = &cobra.Command{
	Use:     "collect",
	Example: "bangumi collect",
	Run: func(cmd *cobra.Command, args []string) {
		t := bangumi.SubjectCollectionType(collectionType)

		h, err := collecthandler.NewHandler(username, t)
		if err != nil {
			log.Fatal(err)
		}

		h.Run()
	},
}

func init() {
	collectCmd.Flags().StringVarP(&username, "username", "u", "", "bangumi 用户名")
	collectCmd.Flags().IntVarP(&collectionType, "type", "t", -1, "bangumi 收藏状态")

	rootCmd.AddCommand(collectCmd)
}
