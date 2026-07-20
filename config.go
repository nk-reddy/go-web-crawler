package main

import (
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	maxPages           int
	baseURL            *url.URL
	mu                 *sync.RWMutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}
