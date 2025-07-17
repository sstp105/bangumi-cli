package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func (h *Handler) sync(remote []model.BangumiBase) []model.BangumiBase {
	cs := NewChangeSet(h.subscription, remote)

	if !cs.HasChanged() {
		log.Infof("Local subscription list matches Mikan, no sync needed")
		return h.subscription
	}

	added, unchanged, removed := cs.Added, cs.Unchanged, cs.Removed
	subscribed := unchanged

	if len(added) > 0 {
		subscribed = h.syncSubscribe(subscribed, added)
	}

	if len(removed) > 0 {
		subscribed = h.syncUnsubscribe(subscribed, removed)
	}

	return subscribed
}

func (h *Handler) syncSubscribe(subscribed, added []model.BangumiBase) []model.BangumiBase {
	log.Infof("There are %d new anime in the Mikan subscription:", len(added))
	for _, item := range added {
		log.Debug(item.Name)
	}

	proceed := prompt.Confirm("Do you want to subscribe locally?")
	if proceed {
		items := h.process(added)
		subscribed = append(subscribed, items...)
	}
	return subscribed
}

func (h *Handler) syncUnsubscribe(subscribed, removed []model.BangumiBase) []model.BangumiBase {
	log.Infof("There are %d anime unsubscribed from Mikan:", len(removed))

	for _, item := range removed {
		log.Debug(item.Name)
	}

	proceed := prompt.Confirm("Do you also want to unsubscribe locally?")
	if proceed {
		unsubscribe(removed)
	} else {
		subscribed = append(subscribed, removed...)
	}
	return subscribed
}
