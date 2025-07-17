package subscribehandler

import (
	"fmt"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/prompt"
	"github.com/sstp105/bangumi-cli/internal/season"
)

type Handler struct {
	bangumiClient *bangumi.Client
	mikanClient   *mikan.Client
	year          int
	seasonID      season.ID
	season        season.Season
	subscription  []model.BangumiBase
}

func NewHandler(year, seasonID int) (*Handler, error) {
	sid := season.ID(seasonID)
	s, err := sid.Season()
	if err != nil {
		return nil, err
	}

	b := bangumi.NewClient()

	m, err := mikan.NewClient(config.MikanClientConfig())
	if err != nil {
		return nil, err
	}

	// load local subscription state, used for comparing the diff
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	return &Handler{
		bangumiClient: b,
		mikanClient:   m,
		year:          year,
		seasonID:      sid,
		season:        s,
		subscription:  subscription,
	}, nil
}

func (h *Handler) Run() {
	remote, err := h.fetch() // user's bangumi subscription on mikan (latest)
	if err != nil {
		log.Errorf("Error reading user's Mikan subscribed anime list (year=%d, seasonID=%d): %s", h.year, h.seasonID, err)
		return
	}

	if h.hasLocalSubscription() {
		h.update(remote)
	} else {
		h.add(remote)
	}

	if err := h.save(); err != nil {
		log.Errorf("Error saving anime subscription config file: %s, please try again", err)
		return
	}

	log.Success("Local subscriptions have been synced with Mikan subscriptions!")
}

func (h *Handler) add(subscription []model.BangumiBase) {
	log.Infof("No local anime subscriptions found, preparing to subscribe to Mikan %d %d subscription list", h.year, h.seasonID)
	h.subscription = h.process(subscription)
}

func (h *Handler) update(subscription []model.BangumiBase) {
	log.Infof("Local subscriptions already exist, preparing to sync Mikan %d %d subscription list", h.year, h.seasonID)
	h.subscription = h.sync(subscription)
}

func (h *Handler) process(data []model.BangumiBase) []model.BangumiBase {
	var subscribed []model.BangumiBase

	for _, item := range data {
		proceed, err := h.subscribe(item)
		if err != nil {
			log.Error("%s subscription error: %s", item.Name, err)
			continue
		}
		if !proceed {
			log.Warnf("Saving this subscription has been canceled. If you want to subscribe, you can run the subscribe command again later.")
			continue
		}
		log.Successf("%s subscription successful!", item.Name)
		subscribed = append(subscribed, item)
	}
	return subscribed
}

func (h *Handler) subscribe(bb model.BangumiBase) (bool, error) {
	b, err := h.parse(bb)
	if err != nil {
		return false, err
	}

	if !confirm() {
		return false, nil
	}

	if err := saveBangumi(*b); err != nil {
		return false, err
	}

	return true, nil
}

func unsubscribe(items []model.BangumiBase) {
	for _, item := range items {
		if err := path.DeleteJSONConfigFile(item.ConfigFileName()); err != nil {
			log.Errorf("Unsubscribe error for %s: %s", item.Name, err)
		}
	}
}

func (h *Handler) fetch() ([]model.BangumiBase, error) {
	year := h.year
	s := h.season

	log.Infof("Reading Mikan %d %s user's subscribed anime list...", year, s.String())

	resp, err := h.mikanClient.GetMyBangumi(h.year, h.season)
	if err != nil {
		return nil, err
	}

	html, err := libs.ParseHTML(resp)
	if err != nil {
		return nil, err
	}

	list, err := mikan.ParseMyBangumiList(html)
	if err != nil {
		return nil, err
	}

	log.Success("Successfully parsed the user's subscribed anime list:")

	for _, item := range list {
		log.Debug(item.Name)
	}

	return list, nil
}

func (h *Handler) parse(b model.BangumiBase) (*model.Bangumi, error) {
	log.Infof("Starting to parse anime: %s, id: %s", b.Name, b.ID)

	resp, err := h.mikanClient.GetBangumi(b.ID)
	if err != nil {
		return nil, err
	}

	html, err := libs.ParseHTML(resp)
	if err != nil {
		return nil, err
	}

	bangumiID, err := mikan.ParseBangumiID(html)
	if err != nil {
		return nil, err
	}

	episodes, err := h.bangumiClient.GetEpisodes(bangumiID)
	if err != nil {
		return nil, err
	}

	rssLink, err := mikan.ParseSubscribedRSSLink(html)
	if err != nil {
		return nil, err
	}

	rss, err := h.mikanClient.LoadRSS(rssLink)
	if err != nil {
		return nil, err
	}

	torrents, filters := filterRSS(*rss)

	return &model.Bangumi{
		BangumiBase: b,
		BangumiID:   bangumiID,
		RSSLink:     rssLink,
		Torrents:    torrents,
		Filters:     *filters,
		Episodes:    episodes,
	}, nil
}

func confirm() bool {
	return prompt.Confirm("Do you want to save this subscription?")
}

func (h *Handler) save() error {
	if err := path.SaveJSONConfigFile(path.SubscriptionConfigFile, h.subscription); err != nil {
		return err
	}
	return nil
}

func saveBangumi(b model.Bangumi) error {
	if err := path.SaveJSONConfigFile(b.ConfigFileName(), b); err != nil {
		return fmt.Errorf("save bangumi config file error: %s", err)
	}
	return nil
}

func (h *Handler) hasLocalSubscription() bool {
	return h.subscription != nil
}
