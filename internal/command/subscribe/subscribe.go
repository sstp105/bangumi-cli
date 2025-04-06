package subscribe

import (
	"errors"
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/parser"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
	"os"
	"path/filepath"
)

func Handler() {
	client, err := mikan.NewClient(config.MikanClientConfig())
	if err != nil {
		log.Fatalf("error creating mikan client:%s", err)
	}

	var list []mikan.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscribedBangumiConfigFile, &list); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("read config file error: %s", err)
		}
	}

	if list == nil {
		console.Infof("本地暂无新番订阅记录, 准备从mikan抓取订阅列表")
		list, err = fetchSubscribedBangumi(client)
		if err != nil {
			log.Fatalf("fetch mikan user subscribed bangumi list error: %s", err)
		}
	}

	for _, item := range list {
		err = processBangumi(client, item)
		if err != nil {
			console.Errorf("解析 %s 失败: %s", item.Name, err)
		}
	}
}

func fetchSubscribedBangumi(client *mikan.Client) ([]mikan.BangumiBase, error) {
	resp, err := client.GetMyBangumi()
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

	console.Infof("成功解析用户订阅的番剧列表:")
	for _, item := range list {
		console.Infof("%s", item.Name)
	}

	err = path.SaveJSONConfigFile(path.SubscribedBangumiConfigFile, list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func processBangumi(client *mikan.Client, bangumiBase mikan.BangumiBase) error {
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

	if err := createBangumiDir(bangumi); err != nil {
		return fmt.Errorf("error creating bangumi dir: %s", err)
	}

	if err := saveBangumiConfig(bangumi); err != nil {
		return fmt.Errorf("error saving bangumi config: %s", err)
	}

	return nil
}

func filterTorrents(rss mikan.RSS) ([]string, error) {
	console.Info("当前订阅的RSS包含以下结果: ")
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

	proceed := prompt.Confirm("是否要保存该订阅? (y/n)")
	if !proceed {
		return nil, errors.New("user aborted")
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

func saveBangumiConfig(bangumi mikan.Bangumi) error {
	fn := fmt.Sprintf("%s.json", bangumi.ID)
	if err := path.SaveJSONConfigFile(fn, bangumi); err != nil {
		return fmt.Errorf("save bangumi config file error: %s", err)
	}

	return nil
}
