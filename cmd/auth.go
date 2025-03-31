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

/**
flow:

0. bangumi init
1. Get auth cookie from mikan 
2. Parse My Subscribed anime 
3. For each anime
	3.1 parse bangumi id
	3.2 retrieve the subscribed RSS
	3.3 create the folder, using the title from mikan
	3.4 create metadata file in the folder with the following information
		- bangumi id
		- also need to call bangumi API to get start episode
		- rss file 
4. can run bangumi download 
	- loop each folder under working dir
	- if rss does not exist -> skip
	- add torrent to qBittorrent but DO NOT START as it will block network speed
5. manually group files to folder
6. bangumi format to rename to jellyfin, plex format


bangumi collect
- bulk collect each anime in bangumi.tv to either hold or doing
*/