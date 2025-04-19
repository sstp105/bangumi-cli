package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/model"
)

type ChangeSet struct {
	Added     []model.BangumiBase
	Unchanged []model.BangumiBase
	Removed   []model.BangumiBase
}

func (cs ChangeSet) HasChanged() bool {
	return len(cs.Added) != 0 || len(cs.Removed) != 0
}

func NewChangeSet(local, remote []model.BangumiBase) ChangeSet {
	var added []model.BangumiBase
	var removed []model.BangumiBase
	var common []model.BangumiBase

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
