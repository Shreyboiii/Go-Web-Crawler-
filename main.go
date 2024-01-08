package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

func main() {

	stopwords, err := LoadStopwords("stopwords-en.json")
	if err != nil {
		log.Fatal(err)
	}

	index := NewDb()
	go DbSearch(&index)

	indexFlag := flag.String("index", "disk", "Specify the index type (e.g., inmem, disk)")

	// Parse the command-line arguments
	flag.Parse()

	// Access the value of the "index" flag
	indexType := *indexFlag

	if indexType == "crawl" {

		fmt.Println("Creating DB..")
		var wgc sync.WaitGroup
		wgc.Add(1)
		go crawlDb("https://www.ucsc.edu", &index, stopwords, &wgc)
		wgc.Wait()
		fmt.Println("Done crawling")

	} else {
		fmt.Println("Opening DB...")
	}

	// Block the main thread to keep the server running
	select {}
}
