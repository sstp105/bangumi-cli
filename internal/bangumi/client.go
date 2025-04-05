package bangumi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
)

const (
	// baseURL is the base URL for the bangumi API.
	baseURL string = "https://api.bgm.tv"

	// getUserCollectionsPath is path for fetching a user's collections.
	getUserCollectionsPath libs.APIPath = "/v0/users/%s/collections"

	// postAccessToken is path for requesting bangumi access token
	postAccessToken libs.APIPath = "/oauth/access_token"

	// defaultPaginationLimit is the default number of items per page when paginating through bangumi API.
	defaultPaginationLimit int = 30
)

var (
	// headers defines the default HTTP headers used for making requests to the bangumi API.
	headers map[string]string = map[string]string{
		"User-Agent":   "github.com/sstp105/bangumi-cli (CLI; Golang)",
		"Content-Type": "application/json",
	}
)

// UserSubjectCollectionResponse represents the response for fetching user subject collections.
type UserSubjectCollectionResponse struct {
	// Total is the total number of subject collections available under the subject and collection type.
	Total int `json:"total"`

	// Limit is the pagination limit as requested.
	Limit int `json:"limit"`

	// Offset is the current position for the returned items, used for pagination.
	Offset int `json:"offset"`

	// Data contains the list of UserSubjectCollection objects representing the user's subject collections.
	Data []UserSubjectCollection `json:"data"`
}

// UserSubjectCollection represents a single subject collection owned by a user.
type UserSubjectCollection struct {
	// Subject holds the details of the subject in the collection.
	Subject SlimSubject `json:"subject"`
}

// SlimSubject represents a simplified version of a subject.
type SlimSubject struct {
	// ID is the unique identifier for the subject.
	ID int `json:"id"`

	// Name is the name of the subject in the original language.
	Name string `json:"name"`

	// NameCN is the name of the subject in Chinese (if available).
	NameCN string `json:"name_cn"`
}

// ErrorResponse represents a generic API error response from the bangumi.
type ErrorResponse struct {
	// Title is a brief summary of the error.
	Title string `json:"title"`

	// Description provides a detailed explanation of the error.
	Description string `json:"description"`

	// Details contains additional information about the error.
	Details string `json:"details"`
}

// Error implements the error interface for ErrorResponse.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("bangumi api error: %s - %s: %s", e.Title, e.Description, e.Details)
}

// Client wraps a resty client for interacting with the bangumi API.
type Client struct {
	client *resty.Client
}

// NewClient creates and returns a new instance of the Client.
func NewClient() *Client {
	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetHeaders(headers)

	return &Client{
		client: c,
	}
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
