package libs

import (
	"fmt"
	"net/http"
)

type APIPath string

func FormatAPIPath(p APIPath, args ...interface{}) string {
	return fmt.Sprintf(string(p), args...)
}

func GetCookie(cookies []*http.Cookie, k string) (string, error) {
	for _, c := range cookies {
		if c.Name == k {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("cookie %s not found", k)
}
