package bangumi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"net/http"
)

const (
	baseURL string = "https://api.bgm.tv"

	getUserCollectionPath   libs.APIPath = "/v0/users/%s/collections/%s"
	postUserCollectionPath  libs.APIPath = "/v0/users/-/collections/%s"
	patchUserCollectionPath libs.APIPath = "/v0/users/-/collections/%s"
	getUserCollectionsPath  libs.APIPath = "/v0/users/%s/collections"
	getEpisodePath          libs.APIPath = "/v0/episodes"
	postAccessToken         libs.APIPath = "/oauth/access_token"

	defaultPaginationLimit int = 30
)

var headers = map[string]string{
	"User-Agent":   "github.com/sstp105/bangumi-cli (CLI; Golang)",
	"Content-Type": "application/json",
}

type Client struct {
	client     *resty.Client
	credential OAuthCredential
}

type Option func(*Client)

func WithAuthorization(credential OAuthCredential) Option {
	return func(c *Client) {
		c.credential = credential
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = resty.NewWithClient(client)
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{}

	client := resty.New()
	c.client = client

	for _, opt := range opts {
		opt(c)
	}

	c.client.SetBaseURL(baseURL)
	c.client.SetHeaders(headers)

	return c
}

func (c *Client) GetUserCollection(username, subjectID string) (*UserSubjectCollection, error) {
	var collection UserSubjectCollection
	var errorResp ErrorResponse

	url := libs.FormatAPIPath(getUserCollectionPath, username, subjectID)

	resp, err := c.client.R().
		SetResult(&collection).
		SetError(&errorResp).
		Get(url)
	if err != nil {
		return nil, err
	}

	// subject is not collected by user
	if resp.StatusCode() == 404 {
		return nil, nil
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	return &collection, nil
}

func (c *Client) PostUserCollection(subjectID string, payload UserSubjectCollectionModifyPayload) error {
	var errorResp ErrorResponse

	url := libs.FormatAPIPath(postUserCollectionPath, subjectID)

	resp, err := c.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.credential.AccessToken)).
		SetBody(payload).
		SetError(&errorResp).
		Post(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return &errorResp
	}

	return nil
}

func (c *Client) PatchUserCollection(subjectID string, payload UserSubjectCollectionModifyPayload) error {
	var errorResp ErrorResponse

	url := libs.FormatAPIPath(patchUserCollectionPath, subjectID)

	resp, err := c.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.credential.AccessToken)).
		SetBody(payload).
		SetError(&errorResp).
		Patch(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return &errorResp
	}

	return nil
}

func (c *Client) GetUserCollections(username string, subjectType, collectionType int) ([]UserSubjectCollection, error) {
	return paginate(func(offset int) ([]UserSubjectCollection, int, error) {
		resp, err := c.GetPaginatedUserCollections(username, subjectType, collectionType, offset)
		if err != nil {
			return nil, 0, err
		}
		return resp.Data, resp.Total, nil
	})
}

func (c *Client) GetPaginatedUserCollections(username string, subjectType, collectionType, offset int) (*UserSubjectCollectionResponse, error) {
	var collections UserSubjectCollectionResponse
	var errorResp ErrorResponse

	params := map[string]string{
		"subject_type": fmt.Sprintf("%d", subjectType),
		"type":         fmt.Sprintf("%d", collectionType),
		"limit":        fmt.Sprintf("%d", defaultPaginationLimit),
		"offset":       fmt.Sprintf("%d", offset),
	}
	url := libs.FormatAPIPath(getUserCollectionsPath, username)

	resp, err := c.client.R().
		SetQueryParams(params).
		SetResult(&collections).
		SetError(&errorResp).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	return &collections, nil
}

func (c *Client) GetEpisodes(subjectID string) ([]Episode, error) {
	return paginate(func(offset int) ([]Episode, int, error) {
		resp, err := c.GetPaginatedEpisodes(subjectID, offset)
		if err != nil {
			return nil, 0, err
		}
		return resp.Data, resp.Total, nil
	})
}

func (c *Client) GetPaginatedEpisodes(subjectID string, offset int) (*EpisodesResponse, error) {
	var episodesResp EpisodesResponse
	var errorResp ErrorResponse

	params := map[string]string{
		"subject_id": subjectID,
		"type":       "0",
		"limit":      fmt.Sprintf("%d", defaultPaginationLimit),
		"offset":     fmt.Sprintf("%d", offset),
	}
	url := libs.FormatAPIPath(getEpisodePath)

	resp, err := c.client.R().
		SetQueryParams(params).
		SetResult(&episodesResp).
		SetError(&errorResp).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, &errorResp
	}

	return &episodesResp, nil
}
