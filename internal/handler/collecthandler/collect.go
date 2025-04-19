package collecthandler

import (
	"errors"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/console"
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
		return nil, errors.New("invalid collection type %d")
	}

	var subscription []model.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &subscription); err != nil {
		return nil, err
	}

	var credential bangumi.OAuthCredential
	if err := path.ReadJSONConfigFile(path.BangumiCredentialConfigFile, &credential); err != nil {
		return nil, err
	}

	client := bangumi.NewClient(bangumi.WithAuthorization(credential))

	return &Handler{
		username:       username,
		collectionType: collectionType,
		subscription:   subscription,
		client:         client,
	}, nil
}

func (h *Handler) Run() {
	for _, s := range h.subscription {
		if err := h.process(s); err != nil {
			console.Errorf("处理 %s 时出错:%s", s.Name, err)
		}
	}
}

func (h *Handler) process(s model.BangumiBase) error {
	id, err := read(s.ConfigFileName())
	if err != nil {
		return err
	}

	if err = h.collect(id); err != nil {
		return err
	}

	return nil
}

func read(fn string) (string, error) {
	var subject model.Bangumi
	if err := path.ReadJSONConfigFile(fn, &subject); err != nil {
		return "", err
	}

	return subject.BangumiID, nil
}

func (h *Handler) collect(id string) error {
	payload := bangumi.UserSubjectCollectionModifyPayload{
		CollectionType: h.collectionType,
	}

	collection, err := h.client.GetUserCollection(h.username, id)
	if err != nil {
		return err
	}

	if collection == nil {
		return h.client.PostUserCollection(id, payload)
	}

	return h.client.PatchUserCollection(id, payload)
}
