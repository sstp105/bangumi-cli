package downloadhandler

import (
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

type Handler struct {
	client       *torrent.QBittorrentClient
	subscription []mikan.BangumiBase
}

func NewHandler() (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	if len(subscription) == 0 {
		console.Infof("本地暂无番剧订阅, 任务结束")
		return nil, nil
	}

	// TODO: factory method that returns torrent.Client based on user input
	client, err := torrent.NewQBittorrentClient(config.QBittorrentConfig())
	if err != nil {
		return nil, err
	}

	return &Handler{
		client:       client,
		subscription: subscription,
	}, nil
}

func (h *Handler) Run() {
	for _, s := range h.subscription {
		if err := h.download(s); err != nil {
			console.Errorf("%s 下载失败: %s", s.Name, err)
		}
	}
}

func (h *Handler) download(bb mikan.BangumiBase) error {
	console.Infof("%s - 保存路径:%s", bb.Name, bb.SavePath())

	var b mikan.Bangumi
	if err := path.ReadJSONConfigFile(bb.ConfigFileName(), &b); err != nil {
		console.Errorf("读取 %s 配置文件错误:%s", bb.Name, err)
		return err
	}

	if err := h.client.Add(b.TorrentURLs(), b.SavePath()); err != nil {
		console.Errorf("添加任务错误:%s", err)
		return err
	}

	console.Successf("%d 个任务已成功添加到 %s!", len(b.Torrents), h.client.Name())
	return nil
}
