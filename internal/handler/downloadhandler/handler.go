package downloadhandler

import (
	"fmt"

	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

type Handler struct {
	client       *torrent.QBittorrentClient
	subscription []model.BangumiBase
	output string
}

func NewHandler(output string) (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, fmt.Errorf("fail to read subscription config file:%w", err)
	}

	client, err := torrent.NewQBittorrentClient(config.QBittorrentConfig())
	if err != nil {
		return nil, fmt.Errorf("fail to create torrent client:%w", err)
	}

	return &Handler{
		client:       client,
		subscription: subscription,
		output: output,
	}, nil
}

func (h *Handler) Run() error {
	if h.subscription == nil {
		log.Warn("No local subscriptions found. Task ended.")
		return nil
	}

	var errs model.ProcessErrors

	for _, s := range h.subscription {
		if err := h.download(s, h.output); err != nil {
			log.Errorf("Failed to add task for %s: %s", s.Name, err)
			errs = append(errs, model.ProcessError{Name: s.Name, Err: err})
		}
	}

	if len(errs) != 0 {
		log.Errorf("A total of %d tasks failed to process:\n%s", len(errs), errs.String())
		return errs
	}

	log.Successf("All tasks have been added to %s, task completed!", h.client.Name())

	return nil
}

func (h *Handler) download(bb model.BangumiBase, output string) error {
	b, err := path.ReadBangumiConfigFile(bb)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", bb.ConfigFileName(), err)
	}

	dest := output + `\` + bb.Name
	log.Debugf("%s - Save path: %s", bb.Name, dest)
	
	if err := h.client.Add(b.TorrentURLs(), dest); err != nil {
		return fmt.Errorf("failed to add torrent urls: %w", err)
	}
	
	log.Successf("%d tasks have been successfully added to %s!", len(b.Torrents), h.client.Name())

	return nil
}
