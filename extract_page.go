package main

import "net/url"

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func extractPageData(html, pageURL string) PageData {
	heading := getHeadingFromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)

	// Parse the page URL once
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		// If it's invalid, bail gracefully with minimal data
		return PageData{
			URL:            pageURL,
			Heading:        heading,
			FirstParagraph: firstParagraph,
			OutgoingLinks:  nil,
			ImageURLs:      nil,
		}
	}

	outgoingLinks, err := getURLsFromHTML(html, parsedURL)
	if err != nil {
		outgoingLinks = nil
	}

	imageURLs, err := getImagesFromHTML(html, parsedURL)
	if err != nil {
		imageURLs = nil
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
