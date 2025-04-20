package unsubscribehandler

import (
	"testing"

	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_filterOutByID(t *testing.T) {
	subscriptions := []model.BangumiBase{
		{ID: "1", Name: "Bangumi One", Link: "https://example.com/1"},
		{ID: "2", Name: "Bangumi Two", Link: "https://example.com/2"},
		{ID: "3", Name: "Bangumi Three", Link: "https://example.com/3"},
	}

	t.Run("remove existing id", func(t *testing.T) {
		result := filterOutByID(subscriptions, 2)

		assert.Len(t, result, 2)
		assert.NotContains(t, result, model.BangumiBase{ID: "2", Name: "Bangumi Two", Link: "https://example.com/2"})
	})

	t.Run("remove non-existing id", func(t *testing.T) {
		result := filterOutByID(subscriptions, 999)

		assert.Len(t, result, 3)
		assert.Equal(t, subscriptions, result)
	})

	t.Run("remove only one item", func(t *testing.T) {
		subs := []model.BangumiBase{
			{ID: "42", Name: "Only", Link: "https://example.com/only"},
		}
		result := filterOutByID(subs, 42)

		assert.Len(t, result, 0)
	})
}
