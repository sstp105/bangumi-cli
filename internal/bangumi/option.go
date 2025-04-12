package bangumi

type Option func(*Client)

func WithAuthorization(credential OAuthCredential) Option {
	return func(c *Client) {
		c.credential = credential
	}
}
