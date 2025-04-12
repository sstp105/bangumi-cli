package bangumi

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient_Default(t *testing.T) {
	client := NewClient()

	assert.NotNil(t, client)
	assert.NotNil(t, client.client)

	assert.Equal(t, baseURL, client.client.BaseURL)

	for k, v := range headers {
		assert.Equal(t, v, client.client.Header.Get(k))
	}
}

func TestNewClient_WithAuthorization(t *testing.T) {
	accessToken := "access_token"
	credential := OAuthCredential{
		AccessToken: accessToken,
	}
	client := NewClient(WithAuthorization(credential))

	assert.Equal(t, accessToken, client.credential.AccessToken)
}

func TestGetUserCollection_Error(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/bangumi-cli/collections/123",
			httpmock.NewErrorResponder(errors.New("context canceled")),
		)

		resp, err := c.GetUserCollection("bangumi-cli", "123")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetUserCollection_NotFound(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/bangumi-cli/collections/123",
			httpmock.NewJsonResponderOrPanic(404, map[string]string{
				"title":       "Not Found",
				"description": "The subject is not collected by user",
			}))

		resp, err := c.GetUserCollection("bangumi-cli", "123")
		assert.NoError(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetUserCollection_ServerError(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/bangumi-cli/collections/123",
			httpmock.NewJsonResponderOrPanic(502, map[string]string{
				"title":       "Bad Gateway",
				"description": "The server is under maintenance",
			}),
		)

		resp, err := c.GetUserCollection("bangumi-cli", "123")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetUserCollection_Success(t *testing.T) {
	collectionType := 4

	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"GET",
			"https://api.bgm.tv/v0/users/bangumi-cli/collections/123",
			httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
				"subject": map[string]interface{}{
					"name":    "メダリスト",
					"name_cn": "金牌得主",
					"type":    2,
					"id":      430699,
				},
				"type": collectionType,
			}),
		)

		resp, err := c.GetUserCollection("bangumi-cli", "123")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.CollectionType, SubjectCollectionType(collectionType))
	})
}

func TestPostUserCollection_Error(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"POST",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewErrorResponder(errors.New("context canceled")))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(3),
		}
		err := c.PostUserCollection("123", payload)

		assert.Error(t, err)
	})
}

func TestPostUserCollection_BadRequest(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"POST",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewJsonResponderOrPanic(502, map[string]string{
				"title":       "Bad Request",
				"description": "The request is not valid",
			}))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(99),
		}
		err := c.PostUserCollection("123", payload)

		assert.Error(t, err)
	})
}

func TestPostUserCollection_Success(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"POST",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewStringResponder(200, ""))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(2),
		}
		err := c.PostUserCollection("123", payload)

		assert.NoError(t, err)
	})
}

func TestPatchUserCollection_Error(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"PATCH",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewErrorResponder(errors.New("context canceled")))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(3),
		}
		err := c.PatchUserCollection("123", payload)

		assert.Error(t, err)
	})
}

func TestPatchUserCollection_BadRequest(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"PATCH",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewJsonResponderOrPanic(502, map[string]string{
				"title":       "Bad Request",
				"description": "The request is not valid",
			}))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(99),
		}
		err := c.PatchUserCollection("123", payload)

		assert.Error(t, err)
	})
}

func TestPatchUserCollection_Success(t *testing.T) {
	withMockClient(t, func(c *Client) {
		httpmock.RegisterResponder(
			"PATCH",
			"https://api.bgm.tv/v0/users/-/collections/123",
			httpmock.NewStringResponder(200, ""))

		payload := UserSubjectCollectionModifyPayload{
			CollectionType: SubjectCollectionType(2),
		}
		err := c.PatchUserCollection("123", payload)

		assert.NoError(t, err)
	})
}

func newMockClient() *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeaders(headers)

	httpmock.ActivateNonDefault(client.GetClient())

	return &Client{
		client: client,
	}
}

func withMockClient(t *testing.T, testFn func(c *Client)) {
	client := newMockClient()
	defer httpmock.DeactivateAndReset()
	testFn(client)
}
