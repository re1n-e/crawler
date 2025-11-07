package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	for _, page := range pages {
		outgoingLinks := strings.Join(page.OutgoingLinks, ";")
		imageUrls := strings.Join(page.ImageURLs, ";")
		writer.Write([]string{page.URL, page.H1, page.FirstParagraph, outgoingLinks, imageUrls})
	}
	return nil
}
