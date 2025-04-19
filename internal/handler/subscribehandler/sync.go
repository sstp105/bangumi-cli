package subscribehandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func (h *Handler) sync(remote []model.BangumiBase) []model.BangumiBase {
	cs := NewChangeSet(h.subscription, remote)

	if !cs.HasChanged() {
		console.Infof("本地订阅列表与 mikan 一致，无需同步")
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
	console.Infof("有 %d 部新的番剧在 mikan 订阅:", len(added))
	for _, item := range added {
		console.Plain(item.Name)
	}

	proceed := prompt.Confirm(fmt.Sprint("是否要在本地订阅?"))
	if proceed {
		items := h.process(added)
		subscribed = append(subscribed, items...)
	}
	return subscribed
}

func (h *Handler) syncUnsubscribe(subscribed, removed []model.BangumiBase) []model.BangumiBase {
	console.Infof("有 %d 部番剧在 mikan 取消了订阅:", len(removed))
	for _, item := range removed {
		console.Plain(item.Name)
	}

	proceed := prompt.Confirm(fmt.Sprint("是否也要在本地取消订阅?"))
	if proceed {
		unsubscribe(removed)
	} else {
		subscribed = append(subscribed, removed...)
	}
	return subscribed
}
