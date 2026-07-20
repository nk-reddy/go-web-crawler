package main

import (
	"fmt"
	"net/url"
)

func sameDomain(rawCurrentURL, rawBaseURL string) bool {
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return false
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return false
	}

	return currentURL.Hostname() == baseURL.Hostname()
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !sameDomain(rawCurrentURL, rawBaseURL) {
		return
	}

	currentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if pages[currentURL] > 0 {
		pages[currentURL]++
		return
	}

	pages[currentURL] = 1
	fmt.Printf("Crawling: %s\n", currentURL)

	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error crawling %s: %v\n", currentURL, err)
		return
	}

	data := extractPageData(pageHTML, rawCurrentURL)
	for _, rawURL := range data.OutgoingLinks {
		crawlPage(rawBaseURL, rawURL, pages)
	}
}
