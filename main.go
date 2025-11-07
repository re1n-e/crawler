package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func printUsage() {
	fmt.Println("usage: ./crawler URL maxConcurrency maxPages")
}

func main() {
	args := os.Args

	if len(args) != 4 {
		printUsage()
		os.Exit(1)
	}
	base_url, err := url.Parse(args[1])
	if err != nil {
		fmt.Printf("err parsing base url: %v\n", err)
		return
	}
	max_concurreny, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("err parsing max concurrecny arg: %v\n", err)
		return
	}
	max_pages, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("err parsing max pages arg: %v\n", err)
		return
	}
	fmt.Printf("starting crawl of: %s\n", base_url)
	cfg := &config{
		maxPages:           max_pages,
		pages:              make(map[string]PageData),
		baseURL:            base_url,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, max_concurreny), // limit to 10 workers
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(args[1])

	cfg.wg.Wait()

	writeCSVReport(cfg.pages, "report.csv")
}
