package bangumi

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
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

// Client wraps a resty client for interacting with the bangumi API.
type Client struct {
	client     *resty.Client
	credential OAuthCredential
}

// NewClient creates and returns a new instance of the Client.
func NewClient(opts ...Option) *Client {
	c := &Client{}

	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeaders(headers)

	for _, opt := range opts {
		opt(c)
	}

	c.client = client

	return c
}

// GetUserCollection retrieves the collection status of a specific subject for a given username.
// If the subject has not been collected by the user, the error is nil.
//
// Parameters:
//   - username: bangumi username (unique)
//   - subjectID: the ID of the subject
//
// Returns:
//   - *UserSubjectCollection: the user's collection data for the subject
//   - error: non 200 api response or http errors
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

// GetUserCollections retrieves all the collections for a user from the bangumi API by paginating through the results.
//
// Parameters:
// - `username`: The username of the user whose collections are to be fetched.
// - `subjectType`: The type of subject (e.g., 2 for 动画) to filter the collections.
// - `collectionType`: The type of collection (e.g., 1 for 在看) to filter the results.
//
// Returns:
// - A slice of `UserSubjectCollection`, which contains all the user's collections.
// - An `error` if there was an issue with making the request or parsing the response.
func (c *Client) GetUserCollections(username string, subjectType, collectionType int) ([]UserSubjectCollection, error) {
	var collections []UserSubjectCollection
	total := 1 // initially set as 1 for first request

	for offset := 0; offset < total; offset += defaultPaginationLimit {
		resp, err := c.GetPaginatedUserCollections(username, subjectType, collectionType, offset)
		if err != nil {
			return nil, err
		}
		collections = append(collections, resp.Data...)
		total = resp.Total
	}

	return collections, nil
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

// GetPaginatedUserCollections fetches the paginated collections for a user from the bangumi API.
//
// Parameters:
// - `username`: The username of the user whose collections are to be fetched.
// - `subjectType`: The type of subject (e.g., 2 for 动画) to filter the collections.
// - `collectionType`: The type of collection (e.g., 1 for 在看) to filter the results.
// - `offset`: The offset for pagination.
//
// Returns:
// - A pointer to `UserSubjectCollectionResponse`, which contains the paginated list of collections for the user.
// - An `error` if there was an issue with making the request or parsing the response.
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

func paginate[T any](fetch func(offset int) ([]T, int, error)) ([]T, error) {
	var result []T
	total := 1

	for offset := 0; offset < total; offset += defaultPaginationLimit {
		data, pageTotal, err := fetch(offset)
		if err != nil {
			return nil, err
		}
		result = append(result, data...)
		total = pageTotal
	}

	return result, nil
}
