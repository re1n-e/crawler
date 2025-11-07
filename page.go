package main

import (
	"log"
	"net/url"
	"sync"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) (PageData, error) {
	baseUrl, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}, err
	}
	H1, err := getH1FromHTML(html)
	if err != nil {
		return PageData{}, err
	}
	FirstParagraph, err := getFirstParagraphFromHTML(html)
	if err != nil {
		return PageData{}, err
	}
	OutgoingLinks, err := getURLsFromHTML(html, baseUrl)
	if err != nil {
		return PageData{}, err
	}
	ImageURLs, err := getImagesFromHTML(html, baseUrl)
	if err != nil {
		return PageData{}, err
	}
	return PageData{
		URL:            pageURL,
		H1:             H1,
		FirstParagraph: FirstParagraph,
		OutgoingLinks:  OutgoingLinks,
		ImageURLs:      ImageURLs,
	}, nil
}

type config struct {
	maxPages           int
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	// acquire concurrency slot
	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	rawCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("can't parse url: %s err: %v", rawCurrentURL, err)
		return
	}

	// Stay on same hostname
	if cfg.baseURL.Hostname() != rawCurrent.Hostname() {
		return
	}

	normalized, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("failed to normalize: %s", rawCurrentURL)
		return
	}

	// dedupe
	if !cfg.addPageVisit(normalized) {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("failed to get html: %v", err)
		return
	}

	pageData, err := extractPageData(html, rawCurrentURL)
	if err != nil {
		log.Printf("can't extract page data: %v", err)
		return
	}

	// Store extracted metadata
	cfg.mu.Lock()
	cfg.pages[normalized] = pageData
	cfg.mu.Unlock()

	// Crawl outgoing links
	for _, link := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}
