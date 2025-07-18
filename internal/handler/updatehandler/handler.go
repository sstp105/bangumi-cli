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
			log.Errorf("Error updating %s: %v", s.Name, err)
		}
	}
	log.Success("Local subscriptions synchronization completed!")
}

func (h *Handler) update(bb model.BangumiBase) error {
	b, err := path.ReadBangumiConfigFile(bb)
	if err != nil {
		return err
	}

	log.Infof("Updating: %s, total episode: %d, added torrents:", bb.Name, len(b.Episodes))
	for _, v := range b.Torrents {
		log.Debug(v.Title)
	}

	rss, err := h.loadRSS(*b)
	if err != nil {
		return err
	}

	d := diff(*rss, b.Filters, b.Torrents)
	if len(d) == 0 {
		log.Debug("RSS is already synced, no new torrents to add")
		return nil
	}

	added := promptAdd(d)
	if len(added) == 0 {
		return nil
	}

	if err := save(*b, added); err != nil {
		return err
	}

	log.Successf("%s update completed!", bb.Name)

	return nil
}

func (h *Handler) loadRSS(b model.Bangumi) (*mikan.RSS, error) {
	log.Debugf("Checking RSS for new torrents...")

	rss, err := h.mikanClient.LoadRSS(b.RSSLink)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS from %s: %w", b.RSSLink, err)
	}
	return rss, nil
}

func diff(rss mikan.RSS, filters model.Filters, torrents []model.Torrent) map[string]string {
	r := rss.Filter(filters)

	mp := make(map[string]string) // key:hash, value:bangumi name
	for _, item := range r.Channel.Items {
		mp[item.Enclosure.URL] = item.Title
	}

	for _, item := range torrents {
		if _, ok := mp[item.Link]; ok {
			delete(mp, item.Link)
		}
	}

	return mp
}

func promptAdd(diff map[string]string) []model.Torrent {
	log.Infof("There are %d new torrents available to add:", len(diff))

	var added []model.Torrent
	for k, v := range diff {
		log.Debug(v)
		added = append(added, model.Torrent{
			Link: k,
			Title: v,
		})
	}

	proceed := prompt.Confirm("Do you want to add them?")
	if !proceed {
		return nil
	}

	return added
}

func save(b model.Bangumi, added []model.Torrent) error {
	log.Debugf("Added %d new torrents", len(added))
	b.Torrents = append(b.Torrents, added...)

	if err := path.SaveJSONConfigFile(b.ConfigFileName(), b); err != nil {
		return fmt.Errorf("save bangumi config file error: %w", err)
	}

	log.Debugf("%s config file saved successfully", b.Name)
	return nil
}
