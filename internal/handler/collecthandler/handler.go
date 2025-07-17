package collecthandler

import (
	"errors"
	"fmt"

	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
)

type Handler struct {
	username       string
	collectionType bangumi.SubjectCollectionType
	subscription   []model.BangumiBase
	client         *bangumi.Client
}

func NewHandler(username string, t int) (*Handler, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

	collectionType := bangumi.SubjectCollectionType(t)
	if !collectionType.IsValid() {
		return nil, fmt.Errorf("invalid collection type %d", collectionType)
	}

	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return nil, errors.New("subscription config file is empty")
	}

	credential, err := path.ReadBangumiCredentialConfigFile()
	if err != nil {
		return nil, err
	}

	if credential == nil {
		return nil, errors.New("credential config file is empty")
	}

	client := bangumi.NewClient(bangumi.WithAuthorization(*credential))

	return &Handler{
		username:       username,
		collectionType: collectionType,
		subscription:   subscription,
		client:         client,
	}, nil
}

func (h *Handler) Run() error {
	var errs model.ProcessErrors

	for _, s := range h.subscription {
		if err := h.process(s); err != nil {
			log.Errorf("Error occurred while processing [%s]: %s", s.Name, err)
			errs = append(errs, model.ProcessError{
				Name: s.Name,
				Err:  err,
			})
		}
	}

	if errs != nil {
		log.Errorf("Total of %d subjects failed to process:\n%s", len(errs), errs.String())
		return errs
	}

	log.Successf("Successfully synced Mikan subscribed subjects to Bangumi as %s, task completed!", h.collectionType.String())

	return nil
}

func (h *Handler) process(s model.BangumiBase) error {
	id, err := getBangumiID(s.ConfigFileName())
	if err != nil {
		return fmt.Errorf("failed to get bangumi id:%w", err)
	}

	if err = h.collect(id); err != nil {
		return fmt.Errorf("failed to collect bangumi %s:%w", id, err)
	}

	log.Infof("Successfully collected [%s] as [%s]", s.Name, h.collectionType.String())

	return nil
}

func (h *Handler) collect(id string) error {
	payload := bangumi.UserSubjectCollectionModifyPayload{
		CollectionType: h.collectionType,
	}

	collection, err := h.client.GetUserCollection(h.username, id)
	if err != nil {
		return fmt.Errorf("failed to fetch collection status:%w", err)
	}

	// if user has not collected before, create the collection
	if collection == nil {
		log.Debugf("User %s has not collected the subject %s yet, adding collection...", h.username, id)
		return h.client.PostUserCollection(id, payload)
	}

	log.Debugf("%s has already collected the subject %s, updating collection status to %s...", h.username, id, h.collectionType)

	return h.client.PatchUserCollection(id, payload)
}

func getBangumiID(fn string) (string, error) {
	var b model.Bangumi
	if err := path.ReadJSONConfigFile(fn, &b); err != nil {
		return "", fmt.Errorf("failed to read bangumi config file:%w", err)
	}

	id := b.BangumiID
	if id == "" {
		return "", errors.New("bangumi id is empty")
	}

	return id, nil
}
