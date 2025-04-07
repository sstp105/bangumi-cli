package subscribe

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/parser"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
	"github.com/sstp105/bangumi-cli/internal/season"
	"os"
	"path/filepath"
)

func Run(year int, seasonID season.ID) {
	client, err := mikan.NewClient(config.MikanClientConfig())
	if err != nil {
		console.Errorf("初始化 mikan 客户端错误:%s", err)
		return
	}

	subscription, err := fetchSubscription(client, year, seasonID)
	if err != nil {
		console.Errorf("读取 mikan 用户订阅的番剧列表(year=%d,seasonID=%d)错误:%s", year, seasonID, err)
		return
	}

	var localSubscription []mikan.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &localSubscription); err != nil && !os.IsNotExist(err) {
		console.Errorf("读取本地番剧订阅配置文件错误:%s", err)
		return
	}

	if localSubscription == nil {
		localSubscription = subscription
		subscribe(client, localSubscription)
	} else {
		localSubscription = sync(client, localSubscription, subscription)
	}

	err = saveSubscriptionConfig(localSubscription)
	if err != nil {
		console.Errorf("保存番剧订阅配置文件错误:%s", err)
	}

	console.Successf("订阅任务结束")
}

func sync(client *mikan.Client, local, remote []mikan.BangumiBase) []mikan.BangumiBase {
	var added []mikan.BangumiBase   // subscribed on mikan but does not exist locally
	var removed []mikan.BangumiBase // removed on mikan but appears locally

	localSet := libs.NewSet[string]()
	remoteSet := libs.NewSet[string]()

	for _, item := range remote {
		remoteSet.Add(item.ID)
	}
	for _, item := range local {
		localSet.Add(item.ID)
	}

	for _, item := range remote {
		if !localSet.Contains(item.ID) {
			added = append(added, item)
		}
	}

	for _, item := range local {
		if !remoteSet.Contains(item.ID) {
			removed = append(removed, item)
		}
	}

	if len(added) == 0 && len(removed) == 0 {
		console.Infof("本地订阅列表与 mikan 一致，无需同步")
		return local
	}

	if len(added) > 0 {
		proceed := prompt.Confirm(fmt.Sprintf("有 %d 部新的番剧在 mikan 订阅, 是否要在本地订阅?", len(added)))
		if proceed {
			subscribe(client, added)
			local = append(local, added...)
		}
	}

	if len(removed) > 0 {
		proceed := prompt.Confirm(fmt.Sprintf("有 %d 部番剧在 mikan 取消了订阅, 是否也要在本地取消订阅?", len(removed)))
		if proceed {
			local = libs.RemoveElements(local, removed)

			for _, item := range removed {
				if err := path.DeleteJSONConfigFile(item.ConfigFileName()); err != nil {
					console.Errorf("取消订阅:%s 错误:%s", item.Name, err)
				}
			}
		}
	}

	console.Successf("本地订阅已和 mikan 订阅同步！")
	return local
}

func subscribe(client *mikan.Client, data []mikan.BangumiBase) {
	for _, item := range data {
		if err := subscribeBangumi(client, item); err != nil {
			console.Error("%s 订阅错误:%s", item.Name, err)
		}
		console.Successf("%s 订阅成功!", item.Name)
	}
}

func fetchSubscription(client *mikan.Client, year int, seasonID season.ID) ([]mikan.BangumiBase, error) {
	s, err := seasonID.Season()
	if err != nil {
		return nil, err
	}

	console.Infof("抓取 mikan %d %s 用户订阅番剧列表...", year, s.String())

	resp, err := client.GetMyBangumi(mikan.WithYearAndSeason(year, s))
	if err != nil {
		return nil, err
	}

	html, err := parser.ParseHTML(resp)
	if err != nil {
		return nil, err
	}

	list, err := mikan.ParseMyBangumiList(html)
	if err != nil {
		return nil, err
	}

	console.Success("成功解析用户订阅的番剧列表:")
	for _, item := range list {
		console.Infof("%s", item.Name)
	}

	return list, nil
}

func subscribeBangumi(client *mikan.Client, bangumiBase mikan.BangumiBase) error {
	id := bangumiBase.ID
	console.Infof("开始解析番剧:%s, id:%s", bangumiBase.Name, id)

	resp, err := client.GetBangumi(id)
	if err != nil {
		return fmt.Errorf("error calling mikan bangumi page: %s", err)
	}

	html, err := parser.ParseHTML(resp)
	if err != nil {
		return fmt.Errorf("failed to parse mikan bangumi to html: %s", err)
	}

	bangumiID, err := mikan.ParseBangumiID(html)
	if err != nil {
		return fmt.Errorf("failed to parse mikan bangumi bangumi.tv id: %s", err)
	}

	rssLink, err := mikan.ParseSubscribedRSSLink(html)
	if err != nil {
		return fmt.Errorf("failed to parse mikan bangumi subscribed fan-sub rss link: %s", err)
	}

	rss, err := client.ReadRSS(rssLink)
	if err != nil {
		return fmt.Errorf("error reading rss link %s: %s", rssLink, err)
	}

	torrents, err := filterTorrents(*rss)
	if err != nil {
		return fmt.Errorf("error filtering bangumi torrents: %s", err)
	}

	bangumi := mikan.Bangumi{
		BangumiBase: bangumiBase,
		BangumiID:   bangumiID,
		RSSLink:     rssLink,
		Torrents:    torrents,
		Filters:     mikan.Filters{},
	}

	if err := saveBangumiConfig(bangumi); err != nil {
		return fmt.Errorf("error saving bangumi config: %s", err)
	}

	return nil
}

func filterTorrents(rss mikan.RSS) ([]string, error) {
	console.Info("当前订阅的RSS包含以下结果:")
	for _, item := range rss.Channel.Items {
		console.Info(item.Title)
	}

	var input string
	input = prompt.ReadUserInput("请输入想要包含的关键词，多个请以英文逗号(,)隔开, 完成后请按 enter")

	include := libs.SplitToSlice(input, ",")
	filters := mikan.Filters{
		Include: include,
	}
	filteredItems := rss.Filter(filters)

	console.Infof("根据关键词:%v, 筛选后的结果如下:", include)
	for _, item := range filteredItems {
		console.Info(item.Title)
	}

	proceed := prompt.Confirm("是否要保存该订阅? (按 n 取消, 任意键继续)")
	if !proceed {
		return nil, nil
	}

	var torrents []string
	for _, item := range filteredItems {
		torrents = append(torrents, item.Enclosure.URL)
	}
	return torrents, nil
}

func createBangumiDir(bangumi mikan.Bangumi) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %s", err)
	}

	dir := filepath.Join(wd, bangumi.Name)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create bangumi folder: %s", err)
	}

	return nil
}

func saveSubscriptionConfig(data []mikan.BangumiBase) error {
	if err := path.SaveJSONConfigFile(path.SubscriptionConfigFile, data); err != nil {
		return err
	}
	return nil
}

func saveBangumiConfig(bangumi mikan.Bangumi) error {
	fn := fmt.Sprintf("%s.json", bangumi.ID)
	if err := path.SaveJSONConfigFile(fn, bangumi); err != nil {
		return fmt.Errorf("save bangumi config file error: %s", err)
	}

	return nil
}
