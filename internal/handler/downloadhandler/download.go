package downloadhandler

import (
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

func Run() {
	var subscription []mikan.BangumiBase
	err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &subscription)
	if err != nil {
		console.Errorf("读取本地订阅配置文件错误:%s", err)
		return
	}

	client, err := torrent.NewQBittorrentClient(config.QBittorrentConfig())
	if err != nil {
		console.Errorf("初始化 bt 客户端错误: %s", err)
		return
	}

	for _, item := range subscription {
		console.Infof("%s - 保存路径:%s", item.Name, item.SavePath())

		var bangumi mikan.Bangumi
		if err := path.ReadJSONConfigFile(item.ConfigFileName(), &bangumi); err != nil {
			console.Errorf("读取 %s 配置文件错误: %s", item.Name, err)
		}

		console.Plainf("%s 读取成功, 共%d个文件准备下载", bangumi.Name, len(bangumi.Torrents))

		if err := client.Add(bangumi.TorrentURLs(), bangumi.SavePath()); err != nil {
			console.Errorf("添加任务错误: %s", err)
			return
		}

		console.Successf("任务已成功添加到%s!", client.Name())
	}
}
