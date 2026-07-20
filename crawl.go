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

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if pages[normalizedURL] > 0 {
		pages[normalizedURL]++
		return
	}

	pages[normalizedURL] = 1
	fmt.Printf("Crawling: %s\n", normalizedURL)

	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error crawling %s: %v\n", normalizedURL, err)
		return
	}

	data := extractPageData(pageHTML, rawCurrentURL)
	for _, rawURL := range data.OutgoingLinks {
		crawlPage(rawBaseURL, rawURL, pages)
	}
}

// re-implementing crawl to allow for concurrency
func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()

	if cfg.reachedMaxPages() {
		return
	}

	if !sameDomain(rawCurrentURL, cfg.baseURL.String()) {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	cfg.concurrencyControl <- struct{}{}
	fmt.Printf("Crawling: %s\n", normalizedURL)

	pageHTML, err := getHTML(rawCurrentURL)
	<-cfg.concurrencyControl

	if err != nil {
		fmt.Printf("Error crawling %s: %v\n", normalizedURL, err)
		return
	}

	pageData := extractPageData(pageHTML, rawCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	for _, rawURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(rawURL)
	}

}

func (cfg *config) reachedMaxPages() (reachedLimit bool) {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	return len(cfg.pages) >= cfg.maxPages
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}
