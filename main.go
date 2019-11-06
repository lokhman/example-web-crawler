package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lokhman/example-web-crawler/crawler"
)

func main() {
	maxPages := flag.Uint("max-pages", 1_000, "maximum number of pages to process")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] url\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		flag.Usage()
		return
	}

	// if URL does not contain scheme, try to fix it
	// TODO: Make more accurate if required.
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}

	// assume that error output is collected by log aggregator
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	webCrawler := crawler.NewCrawler(*maxPages, logger)
	err := webCrawler.Stream(os.Stdout, url)
	if err != nil {
		logger.Println(err)
	}

	fmt.Fprintf(os.Stdout, "Total pages: %d\n", webCrawler.Count())
}
