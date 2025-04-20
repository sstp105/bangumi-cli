package updatehandler

import (
	"errors"
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

type Handler struct {
	mikanClient  *mikan.Client
	subscription []model.BangumiBase
}

func NewHandler(cfg mikan.ClientConfig) (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return nil, errors.New("no subscription config found")
	}

	client, err := mikan.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Handler{
		mikanClient:  client,
		subscription: subscription,
	}, nil
}

func (h *Handler) Run() {
	for _, s := range h.subscription {
		if err := h.update(s); err != nil {
			log.Errorf("更新 %s 出错: %v", s.Name, err)
		}
	}
	log.Success("本地订阅同步完成!")
}

func (h *Handler) update(bb model.BangumiBase) error {
	log.Infof("更新:%s", bb.Name)

	b, err := path.ReadBangumiConfigFile(bb)
	if err != nil {
		return err
	}

	rss, err := h.loadRSS(*b)
	if err != nil {
		return err
	}

	d := diff(*rss, b.Filters, b.Torrents)
	if len(d) == 0 {
		log.Debug("已同步 RSS, 暂无新的种子可添加")
		return nil
	}

	added := promptAdd(d)
	if len(added) == 0 {
		return nil
	}

	if err := save(*b, added); err != nil {
		return err
	}

	log.Successf("%s 更新完成!", bb.Name)

	return nil
}

func (h *Handler) loadRSS(b model.Bangumi) (*mikan.RSS, error) {
	log.Infof("查询 RSS 是否有新的种子...")

	rss, err := h.mikanClient.LoadRSS(b.RSSLink)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS from %s:%w", b.RSSLink, err)
	}
	return rss, nil
}

func diff(rss mikan.RSS, filters model.Filters, torrents []string) map[string]string {
	r := rss.Filter(filters)

	mp := make(map[string]string) // key:hash, value:bangumi name
	for _, item := range r.Channel.Items {
		mp[item.Enclosure.URL] = item.Title
	}

	for _, item := range torrents {
		if _, ok := mp[item]; ok {
			delete(mp, item)
		}
	}

	return mp
}

func promptAdd(diff map[string]string) []string {
	log.Infof("有 %d 个新的种子可添加:", len(diff))

	var added []string
	for k, v := range diff {
		log.Debug(v)
		added = append(added, k)
	}

	proceed := prompt.Confirm("是否要添加?")
	if !proceed {
		return nil
	}

	return added
}

func save(b model.Bangumi, added []string) error {
	log.Debugf("添加了 %d 个新的种子", len(added))
	b.Torrents = append(b.Torrents, added...)

	if err := path.SaveJSONConfigFile(b.ConfigFileName(), b); err != nil {
		return fmt.Errorf("save bangumi config file error: %w", err)
	}

	log.Debug("%s 配置文件保存成功", b.Name)
	return nil
}
