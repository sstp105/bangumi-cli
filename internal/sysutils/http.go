package sysutils

import (
	"fmt"
	"net/http"
)

func GetCookie(cookies []*http.Cookie, k string) (string, error) {
	for _, c := range cookies {
		if c.Name == k {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("cookie %s not found", k)
}
