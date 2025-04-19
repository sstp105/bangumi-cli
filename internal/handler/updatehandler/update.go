package updatehandler

import (
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

type Handler struct {
	mikanClient  *mikan.Client
	subscription []model.BangumiBase
}

func NewHandler() (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	client, err := mikan.NewClient(config.MikanClientConfig())
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
			console.Errorf("更新 %s 出错: %v", s.Name, err)
		}
	}
	console.Successf("本地订阅同步完成!")
}

func (h *Handler) update(bb model.BangumiBase) error {
	console.Infof("更新:%s", bb.Name)

	var b model.Bangumi
	if err := path.ReadJSONConfigFile(bb.ConfigFileName(), &b); err != nil {
		return err
	}

	console.Infof("查询 RSS 是否有新的种子...")
	rss, err := h.mikanClient.LoadRSS(b.RSSLink)
	if err != nil {
		return err
	}

	d, err := h.diff(*rss, b.Filters, b.Torrents)

	added, err := h.promptAdd(d)
	if err != nil {
		return err
	}

	if len(added) == 0 {
		return nil
	}

	b.Torrents = append(b.Torrents, added...)

	if err := path.SaveJSONConfigFile(b.ConfigFileName(), b); err != nil {
		return err
	}

	console.Successf("%s 更新完成", bb.Name)

	return nil
}

func (h *Handler) diff(rss mikan.RSS, filters model.Filters, torrents []string) (map[string]string, error) {
	r := rss.Filter(filters)

	mp := make(map[string]string) // key:hash, value:name
	for _, item := range r.Channel.Items {
		mp[item.Enclosure.URL] = item.Title
	}

	for _, item := range torrents {
		if _, ok := mp[item]; ok {
			delete(mp, item)
		}
	}
	return mp, nil
}

func (h *Handler) promptAdd(diff map[string]string) ([]string, error) {
	sz := len(diff)

	if sz == 0 {
		console.Plain("已同步 RSS, 暂无新的种子可添加")
		return nil, nil
	}

	console.Infof("有 %d 个新的种子可添加:", sz)
	var added []string
	for k, v := range diff {
		console.Plain(v)
		added = append(added, k)
	}

	proceed := prompt.Confirm("是否要添加?")
	if !proceed {
		console.Plain("更新已取消")
		return nil, nil
	}

	return added, nil
}
