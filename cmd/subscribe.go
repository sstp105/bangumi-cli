package cmd

import (
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use: "subscribe",
	Run: func(cmd *cobra.Command, args []string) {
		//client := mikan.NewClient()
		//
		//var list []mikan.BangumiBase
		//err := libs.ReadJSONConfigFile(libs.SubscribedBangumiConfigFile, &list)
		//if err != nil {
		//	log.Fatalf("read config file error: %s", err)
		//}
		//
		//if list == nil {
		//	log.Infof("no subscribed bangumi found, fetching from mikan...")
		//	list, err = fetchSubscribedBangumi(client)
		//	if err != nil {
		//		log.Fatalf("fetch bangumi list error: %s", err)
		//	}
		//}
		//
		//for _, item := range list {
		//	html, err := client.GetBangumi(item.Link)
		//	if err != nil {
		//		log.Fatalf("failed to fetch bangumi list %s: %s", item.Link, err)
		//	}
		//
		//	bangumiID, err := mikan.ParseBangumiID(html)
		//	if err != nil {
		//		log.Fatalf("failed to parse bangumi id %s: %s", html, err)
		//	}
		//
		//	rssLink, err := mikan.ParseSubscribedRSSLink(html)
		//	if err != nil {
		//		log.Fatalf("failed to parse subscribed rs link %s: %s", html, err)
		//	}
		//
		//	log.Infof("bangumiID: %s", bangumiID)
		//	log.Debugf("Fetching RSS link: %s...", rssLink)
		//
		//	rss, err := mikan.LoadRSS(rssLink)
		//	if err != nil {
		//		log.Fatalf("failed to load bangumi rss: %s", err)
		//	}
		//
		//	for _, r := range rss.Channel.Items {
		//		log.Infof("%s", r.Title)
		//	}
		//
		//	var input string
		//	input = libs.ReadInput("Please enter word that must include, split by comma, hit enter to proceed: ")
		//
		//	include := libs.SplitToSlice(input)
		//	includeItems := rss.FilterInclude(include)
		//	log.Infof("Updated list: ")
		//	for _, item := range includeItems {
		//		log.Infof("%s", item.Title)
		//	}
		//
		//	proceed := libs.Confirm("Do you want to save?")
		//	if !proceed {
		//		continue
		//	}
		//
		//	var torrents []string
		//	for _, item := range includeItems {
		//		torrents = append(torrents, item.Enclosure.URL)
		//	}
		//
		//	cfgFile := mikan.Bangumi{
		//		BangumiBase: item,
		//		BangumiID:   bangumiID,
		//		RSSLink:     rssLink,
		//		Torrents:    torrents,
		//		Filters:     mikan.Filters{},
		//	}
		//
		//	data, err := libs.MarshalJSONIndented(cfgFile)
		//	if err != nil {
		//		log.Fatalf("failed to marshal bangumi config: %s", err)
		//	}
		//
		//	wd, err := os.Getwd()
		//	if err != nil {
		//		log.Fatalf("failed to get current working directory: %s", err)
		//	}
		//
		//	dirPath := filepath.Join(wd, item.Name)
		//	err = os.MkdirAll(dirPath, os.ModePerm)
		//	if err != nil {
		//		log.Fatalf("failed to create directory: %s", err)
		//	}
		//
		//	filePath := filepath.Join(dirPath, ".config.json")
		//
		//	err = os.WriteFile(filePath, data, 0644)
		//	if err != nil {
		//		log.Fatalf("failed to save config: %s", err)
		//	}
		//}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

//func fetchSubscribedBangumi(client *mikan.Client) ([]mikan.BangumiBase, error) {
//	html, err := client.GetMyBangumi()
//	if err != nil {
//		return nil, err
//	}
//
//	list, err := mikan.ParseMyBangumiList(html)
//	if err != nil {
//		return nil, err
//	}
//
//	err = libs.SaveJSONConfigFile(libs.SubscribedBangumiConfigFile, list)
//	if err != nil {
//		return nil, err
//	}
//
//	return list, nil
//}
