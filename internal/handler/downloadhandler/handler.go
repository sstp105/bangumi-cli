package downloadhandler

import (
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

type Handler struct {
	client       *torrent.QBittorrentClient
	subscription []model.BangumiBase
}

func NewHandler() (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	if len(subscription) == 0 {
		log.Infof("本地暂无番剧订阅, 任务结束")
		return nil, nil
	}

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
	var errs []error

	for _, s := range h.subscription {
		if err := h.download(s); err != nil {
			log.Errorf("%s 下载失败: %s", s.Name, err)
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		log.Successf("任务已全部添加至 %s, 任务完成", h.client.Name())
	}
}

func (h *Handler) download(bb model.BangumiBase) error {
	log.Infof("%s - 保存路径:%s", bb.Name, bb.SavePath())

	var b model.Bangumi
	if err := path.ReadJSONConfigFile(bb.ConfigFileName(), &b); err != nil {
		log.Errorf("读取 %s 配置文件错误:%s", bb.Name, err)
		return err
	}

	if err := h.client.Add(b.TorrentURLs(), b.SavePath()); err != nil {
		log.Errorf("添加任务错误:%s", err)
		return err
	}

	log.Successf("%d 个任务已成功添加到 %s!", len(b.Torrents), h.client.Name())
	return nil
}
