package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeSet_HasChanged(t *testing.T) {
	t.Run("Empty ChangeSet", func(t *testing.T) {
		cs := ChangeSet{}
		require.False(t, cs.HasChanged(), "empty changeset should not have changes")
	})

	t.Run("Added Items", func(t *testing.T) {
		cs := ChangeSet{
			Added: []model.BangumiBase{{ID: "1"}},
		}
		require.True(t, cs.HasChanged(), "added items should indicate changes")
	})

	t.Run("Removed Items", func(t *testing.T) {
		cs := ChangeSet{
			Removed: []model.BangumiBase{{ID: "1"}},
		}
		require.True(t, cs.HasChanged(), "removed items should indicate changes")
	})

	t.Run("Only Unchanged Items", func(t *testing.T) {
		cs := ChangeSet{
			Unchanged: []model.BangumiBase{{ID: "1"}},
		}
		require.False(t, cs.HasChanged(), "unchanged items should not indicate changes")
	})
}

func TestNewChangeSet(t *testing.T) {
	t.Run("Normal Case", func(t *testing.T) {
		local := []model.BangumiBase{
			{ID: "1"},
			{ID: "2"},
		}
		remote := []model.BangumiBase{
			{ID: "2"},
			{ID: "3"},
		}

		cs := NewChangeSet(local, remote)

		require.Len(t, cs.Added, 1, "should have one added item")
		require.Equal(t, "3", cs.Added[0].ID)

		require.Len(t, cs.Removed, 1, "should have one removed item")
		require.Equal(t, "1", cs.Removed[0].ID)

		require.Len(t, cs.Unchanged, 1, "should have one unchanged item")
		require.Equal(t, "2", cs.Unchanged[0].ID)
	})

	t.Run("Both Local and Remote Empty", func(t *testing.T) {
		cs := NewChangeSet(nil, nil)

		require.Empty(t, cs.Added)
		require.Empty(t, cs.Removed)
		require.Empty(t, cs.Unchanged)
	})

	t.Run("All Remote Items New", func(t *testing.T) {
		var local []model.BangumiBase
		remote := []model.BangumiBase{{ID: "1"}, {ID: "2"}}

		cs := NewChangeSet(local, remote)

		require.Len(t, cs.Added, 2)
		require.Empty(t, cs.Removed)
		require.Empty(t, cs.Unchanged)
	})

	t.Run("All Local Items Removed", func(t *testing.T) {
		local := []model.BangumiBase{{ID: "1"}, {ID: "2"}}
		var remote []model.BangumiBase

		cs := NewChangeSet(local, remote)

		require.Len(t, cs.Removed, 2)
		require.Empty(t, cs.Added)
		require.Empty(t, cs.Unchanged)
	})
}
