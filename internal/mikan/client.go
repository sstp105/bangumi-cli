package mikan

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

const (
	baseURL = "https://mikanani.me"

	myBangumiPath = "/Home/MyBangumi"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Referer":         "https://mikanani.me/",
		"Accept-Language": "en-US,en;q=0.9",
	})

	cookies := []*http.Cookie{
		{
			Name:  ".AspNetCore.Identity.Application",
			Value: "CfDJ8MyNMqFNaC9JmJW13PvY-93aPVQmcCyLPRZTnT3W3zvgOuH79yWk7oRY5ZgYqYbGGG4c-AtrSkMCzSBaOHz1dk1tZPZH4t140CQD6ENGsQyOY1djFGfHTdIh9YXYThxx880ZnZ1ylpBgMcDMnFAMHisU9fGhwaoaM1mkfF3LI4mCBFrcgH8TI3Cr4Bm6_6cycZPmqdFwcY8PTjaXIPT_eMhEe5_pVTzonHmbrnUz6jKrM5ChUZfZGM_MUmGyPrvFVNDxC1HD1Q_rQk7ydLTOfp5OeT7gHgfN7D5Yo-ur0lCs-5I5OdAATr3oEYNyQ-vol4NCu7StIGeBy5AEg3kDAP6mehsfNO8jOE7QBZNNP1g_UNN0p6IgouFts-ml0NR4wjr7bqTB3ctCOFZznRodSZqu6ep85HiY1OKRzK_S1_xwmomOA6berBqf_tNWxq479CpcurrWYP-fRtjnGTmFJnTwxMUIFB1fF9g7GiQrT417Ev4LpU6mlq_U9H_Je6NBd2nKgulve9A9t9WkLJUTU4n659Pd93luUvSLU1EX3h7YelsOtkF8ZGKkh3YlLQDWeYn2W4UEbU7jvpO6KJ9tHz-S2PpuQ7IJYMceVRhxJchiygy_hQFNikMJpaU4ze9C4vKBe2BIM6LIhcn6YFs8ssal5tDRl5YRxQsFWd4hCwmnFazgSH8DANsDbgR_4wmZFCbr5K9aFozqgumd3yKAjiA",
		},
	}
	c.SetCookies(cookies)
	return &Client{
		client: c,
	}
}

func (c *Client) GetSubscribedAnimation() (string, error) {
	resp, err := c.client.R().
		Get(myBangumiPath)

	if err != nil {
		return "", fmt.Errorf("failed to fetch subscribed animation: %w", err)
	}

	return resp.String(), nil
}

func (c *Client) GetBangumi(path string) (string, error) {
	resp, err := c.client.R().
		Get(path)

	if err != nil {
		return "", fmt.Errorf("failed to fetch bangumi page: %w", err)
	}

	return resp.String(), nil
}
