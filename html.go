package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) (string, error) {
	// Load the html cotent
	reader := strings.NewReader(html)
	docc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("-Failed pasring h1 tag: %v", err)
	}

	var h1Text string

	docc.Find("h1").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		h1Text = strings.TrimSpace(s.Text())
		return false
	})

	return h1Text, nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed parsing HTML: %v", err)
	}

	// First try: <main> p
	if mainSel := doc.Find("main p").First(); mainSel.Length() > 0 {
		return strings.TrimSpace(mainSel.Text()), nil
	}

	// Fallback: first <p>
	if p := doc.Find("p").First(); p.Length() > 0 {
		return strings.TrimSpace(p.Text()), nil
	}

	return "", nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to parse html body: %v", err)
	}

	var urls []string

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		// Parse the URL
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		// Convert relative → absolute
		u = baseURL.ResolveReference(u)

		urls = append(urls, u.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to parse html body: %v", err)
	}

	var urls []string

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}

		if src == "" {
			return
		}

		//Parse url
		u, err := url.Parse(src)
		if err != nil {
			return
		}

		// Convert realtive to absolute
		u = baseURL.ResolveReference(u)
		urls = append(urls, u.String())
	})

	return urls, nil
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate GET request: %v", err)
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("bad status code (not OK): %v", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "text/html") {
		return "", fmt.Errorf("expected text/html, got: %v", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)
	}

	return string(body), nil
}
