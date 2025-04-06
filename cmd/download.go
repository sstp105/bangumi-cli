package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

var downloadCmd = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		var local []mikan.BangumiBase
		err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &local)
		if err != nil {
			console.Error(err.Error())
			return
		}

		client, err := torrent.NewQBittorrentClient(config.QBittorrentConfig())
		if err != nil {
			console.Error(err.Error())
			return
		}

		for _, item := range local {
			var bangumi mikan.Bangumi
			err := path.ReadJSONConfigFile(item.ConfigFileName(), &bangumi)
			if err != nil {
				console.Error(err.Error())
			}

			for _, t := range bangumi.Torrents {
				err := client.Add(t)
				if err != nil {
					console.Error(err.Error())
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
