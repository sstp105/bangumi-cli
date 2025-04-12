package update

import (
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func Run() {
	var subscription []mikan.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &subscription); err != nil {
		console.Errorf("本地没有任何订阅, 请先运行 bangumi subscribehandler: %v", err)
		return
	}

	client, err := mikan.NewClient(config.MikanClientConfig())
	if err != nil {
		console.Errorf("初始化 mikan 客户端错误: %v", err)
	}

	for _, s := range subscription {
		if err := update(client, s); err != nil {
			console.Errorf("更新 %s 出错: %v", s.Name, err)
		}
		console.Successf("%s 更新完成\n", s.Name)
	}

	console.Successf("本地订阅同步完成!")
}

func update(client *mikan.Client, b mikan.BangumiBase) error {
	var bangumi mikan.Bangumi
	if err := path.ReadJSONConfigFile(b.ConfigFileName(), &bangumi); err != nil {
		return err
	}

	rss, err := client.ReadRSS(bangumi.RSSLink)
	if err != nil {
		return err
	}

	r := mikan.Filter(*rss, bangumi.Filters)
	mp := make(map[string]string)
	for _, item := range r.Channel.Items {
		mp[item.Enclosure.URL] = item.Title
	}

	for _, item := range bangumi.Torrents {
		if _, ok := mp[item]; ok {
			delete(mp, item)
		}
	}

	if len(mp) == 0 {
		console.Plainf("%s 已同步 RSS, 暂无新的种子可添加", bangumi.Name)
		return nil
	}

	console.Infof("%s 有 %d 个新的种子可添加", bangumi.Name, len(mp))
	var added []string
	for k, v := range mp {
		console.Plain(v)
		added = append(added, k)
	}

	proceed := prompt.Confirm("是否要添加?")
	if !proceed {
		console.Infof("%s 更新已取消", bangumi.Name)
		return nil
	}

	bangumi.Torrents = append(bangumi.Torrents, added...)

	console.Infof("已成功添加 %d 个种子", len(added))

	if err := path.SaveJSONConfigFile(bangumi.ConfigFileName(), bangumi); err != nil {
		return err
	}

	return nil
}
