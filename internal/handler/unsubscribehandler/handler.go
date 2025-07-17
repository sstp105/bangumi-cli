package unsubscribehandler

import (
	"fmt"

	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
)

type Handler struct {
	id           int
	subscription []model.BangumiBase
}

func NewHandler(id int) (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to init handler: %w", err)
	}

	return &Handler{id, subscription}, nil
}

func (h *Handler) Run() {
	switch h.id {
	case -1:
		UnsubscribeAll(h.subscription)
	default:
		UnsubscribeByID(h.subscription, h.id)
	}
	log.Info("Successfully unsubscribed all animations, task ended")
}

func UnsubscribeAll(subscriptions []model.BangumiBase) {
	var failures []model.BangumiBase
	for _, sub := range subscriptions {
		if err := path.DeleteJSONConfigFile(sub.ConfigFileName()); err != nil {
			log.Warnf("Error unsubscribing %s: %s", sub.Name, err)
			failures = append(failures, sub)
			continue
		}
		log.Debugf("Unsubscribed: %s", sub.Name)
	}

	if err := path.DeleteJSONConfigFile(path.SubscriptionConfigFile); err != nil {
		log.Errorf("Failed to delete subscription list: %s", err)
	}

	if len(failures) != 0 {
		log.Warnf("%d anime subscriptions were not successfully removed, please retry or use --id to delete:", len(failures))
		for _, b := range failures {
			log.Debug(b.Name)
		}
		return
	}

	log.Successf("Successfully unsubscribed from %d anime", len(subscriptions))
}


func UnsubscribeByID(subscriptions []model.BangumiBase, id int) {
	updated := filterOutByID(subscriptions, id)
	if len(updated) == len(subscriptions) {
		log.Errorf("Subscription with ID %d not found", id)
		return
	}

	if err := path.SaveJSONConfigFile(path.SubscriptionConfigFile, updated); err != nil {
		log.Errorf("Failed to save subscription config: %s", err)
		return
	}

	fn := fmt.Sprintf("%d.json", id)
	if err := path.DeleteJSONConfigFile(fn); err != nil {
		log.Errorf("Failed to delete config file %s: %s", fn, err)
		return
	}

	log.Successf("Successfully unsubscribed: %d", id)
}

func filterOutByID(subscriptions []model.BangumiBase, id int) []model.BangumiBase {
	targetID := fmt.Sprintf("%d", id)

	var result []model.BangumiBase
	for _, s := range subscriptions {
		if s.ID != targetID {
			result = append(result, s)
		}
	}
	return result
}
