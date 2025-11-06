package main

import (
	"fmt"
	"net/url"
)

func normalizeURL(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	normalizedURL := parsedURL.Hostname() + parsedURL.Path

	if normalizedURL[len(normalizedURL)-1] == '/' {
		normalizedURL = normalizedURL[:len(normalizedURL)-1]
	}

	return normalizedURL, nil
}
