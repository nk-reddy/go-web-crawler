package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	cfg := config{
		pages: map[string]PageData{},
		mu:    &sync.RWMutex{},
		wg:    &sync.WaitGroup{},
	}

	if len(os.Args) < 4 {
		fmt.Println("usage: web-crawler <url> <max-concurrency> <max-pages>")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("invalid max concurrency %q: must be an integer\n", os.Args[2])
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("invalid max pages %q: must be an integer\n", os.Args[3])
		os.Exit(1)
	}

	cfg.baseURL = parsedBaseURL
	cfg.concurrencyControl = make(chan struct{}, maxConcurrency)
	cfg.maxPages = int(maxPages)

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()

	writeJSONReport(cfg.pages, "report.json")
}
