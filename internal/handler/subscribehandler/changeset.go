package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/mikan"
)

type ChangeSet struct {
	Added     []mikan.BangumiBase
	Unchanged []mikan.BangumiBase
	Removed   []mikan.BangumiBase
}

func (cs ChangeSet) HasChanged() bool {
	return len(cs.Added) != 0 || len(cs.Removed) != 0
}

func NewChangeSet(local, remote []mikan.BangumiBase) ChangeSet {
	var added []mikan.BangumiBase
	var removed []mikan.BangumiBase
	var common []mikan.BangumiBase

	localSet := libs.NewSet[string]()
	remoteSet := libs.NewSet[string]()

	for _, item := range remote {
		remoteSet.Add(item.ID)
	}
	for _, item := range local {
		localSet.Add(item.ID)
	}

	for _, item := range remote {
		if !localSet.Contains(item.ID) {
			added = append(added, item)
		} else {
			common = append(common, item)
		}
	}

	for _, item := range local {
		if !remoteSet.Contains(item.ID) {
			removed = append(removed, item)
		}
	}

	return ChangeSet{
		Added:     added,
		Unchanged: common,
		Removed:   removed,
	}
}
