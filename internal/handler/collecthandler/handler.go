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

func NewHandler(username string, collectionType bangumi.SubjectCollectionType) (*Handler, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

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

func (h *Handler) Run() {
	var errs []error

	for _, s := range h.subscription {
		if err := h.process(s); err != nil {
			log.Errorf("处理 %s 时出错:%s", s.Name, err)
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		log.Successf("已同步 mikan 订阅的番剧到 bangumi %s, 任务完成!", h.collectionType.String())
	}
}

func (h *Handler) process(s model.BangumiBase) error {
	id, err := getBangumiID(s.ConfigFileName())
	if err != nil {
		log.Errorf("error getting bangumi id: %v", err)
		return err
	}

	if err = h.collect(id); err != nil {
		log.Errorf("error collecting subject %s", id)
		return err
	}

	log.Infof("收藏 %s 成功 (%s)", s.Name, h.collectionType.String())

	return nil
}

func (h *Handler) collect(id string) error {
	payload := bangumi.UserSubjectCollectionModifyPayload{
		CollectionType: h.collectionType,
	}

	collection, err := h.client.GetUserCollection(h.username, id)
	if err != nil {
		log.Errorf("error fetching user %s collection status for %s:%s", h.username, id, err)
		return err
	}

	// if user has not collected before, create the collection
	if collection == nil {
		log.Debugf("user %s has not collected %s before, creating collection", h.username, id)
		return h.client.PostUserCollection(id, payload)
	}

	log.Debugf("user %s already collected %s, updating collection status to %s", h.username, id, h.collectionType)
	return h.client.PatchUserCollection(id, payload)
}

func getBangumiID(fn string) (string, error) {
	var b model.Bangumi
	if err := path.ReadJSONConfigFile(fn, &b); err != nil {
		return "", err
	}

	return b.BangumiID, nil
}
