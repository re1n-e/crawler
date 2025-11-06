package main

import "net/url"

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
