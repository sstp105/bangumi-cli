package parser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// ParseHTML parses an HTML string and returns a *goquery.Document that can be used for further parsing.
func ParseHTML(s string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseSuffixID extracts an ID from a URL-like string (href) in the form of "/<type>/127791".
func ParseSuffixID(href string) (string, error) {
	if len(href) == 0 {
		return "", errors.New("href is empty, unable to parse the id")
	}

	parts := strings.Split(href, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("href is invalid, split by / returns %s", parts)
	}

	return parts[len(parts)-1], nil
}
