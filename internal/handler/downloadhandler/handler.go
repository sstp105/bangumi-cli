package downloadhandler

import (
	"fmt"
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
		return nil, fmt.Errorf("fail to read subscription config file:%w", err)
	}

	client, err := torrent.NewQBittorrentClient(config.QBittorrentConfig())
	if err != nil {
		return nil, fmt.Errorf("fail to create torrent client:%w", err)
	}

	return &Handler{
		client:       client,
		subscription: subscription,
	}, nil
}

func (h *Handler) Run() error {
	if h.subscription == nil {
		log.Warn("本地暂无任何订阅，任务结束。")
		return nil
	}

	var errs model.ProcessErrors

	for _, s := range h.subscription {
		if err := h.download(s); err != nil {
			log.Errorf("%s 添加任务失败: %s", s.Name, err)
			errs = append(errs, model.ProcessError{Name: s.Name, Err: err})
		}
	}

	if len(errs) != 0 {
		log.Errorf("共有 %d 个任务处理失败: \n%s", len(errs), errs.String())
		return errs
	}

	log.Successf("任务已全部添加至 %s, 任务完成！", h.client.Name())
	return nil
}

func (h *Handler) download(bb model.BangumiBase) error {
	log.Debugf("%s - 保存路径:%s", bb.Name, bb.SavePath())

	b, err := path.ReadBangumiConfigFile(bb)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", bb.ConfigFileName(), err)
	}

	if err := h.client.Add(b.TorrentURLs(), b.SavePath()); err != nil {
		return fmt.Errorf("failed to add torrent urls: %w", err)
	}

	log.Successf("%d 个任务已成功添加到 %s!", len(b.Torrents), h.client.Name())
	return nil
}
