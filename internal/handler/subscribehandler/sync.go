package subscribehandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func (h *Handler) sync(remote []mikan.BangumiBase) []mikan.BangumiBase {
	cs := NewChangeSet(h.subscription, remote)

	if !cs.HasChanged() {
		console.Infof("本地订阅列表与 mikan 一致，无需同步")
		return h.subscription
	}

	added, unchanged, removed := cs.Added, cs.Unchanged, cs.Removed
	var subscribed = unchanged

	if len(added) > 0 {
		h.syncSubscribe(subscribed, added)
	}

	if len(removed) > 0 {
		h.syncUnsubscribe(subscribed, removed)
	}

	return subscribed
}

func (h *Handler) syncSubscribe(subscribed, added []mikan.BangumiBase) {
	console.Infof("有 %d 部新的番剧在 mikan 订阅:", len(added))
	for _, item := range added {
		console.Plain(item.Name)
	}

	proceed := prompt.Confirm(fmt.Sprint("是否要在本地订阅?"))
	if proceed {
		items := h.process(added)
		subscribed = append(subscribed, items...)
	}
}

func (h *Handler) syncUnsubscribe(subscribed, removed []mikan.BangumiBase) {
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
}
