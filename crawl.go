package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var visited = map[string]bool{}

type dOut struct {
	body []byte
	link string
}

type extrOut struct {
	Words  []string
	Hrefs  []string
	Title  string
	Images []image
	body   []byte
	link   string
}

func getSitemapLinks(u string) TheUrls {
	pols := NewPolicies()

	link := pols.Sitemap(u)

	resp, err := http.Get(link)
	if err != nil {
		log.Println("getSitemaps err: ", err)
	}

	bod, er := io.ReadAll(resp.Body)

	if er != nil {
		log.Println("io error getsitemaps ", er)
	}

	var s SitemapIndex

	xml.Unmarshal(bod, &s)

	var links TheUrls
	for _, sm := range s.Sitemaps {
		resp, err = http.Get(sm.Loc)
		if err != nil {
			log.Println("get sitemap err: ", err)
		}

		bod, er = io.ReadAll(resp.Body)

		if er != nil {
			log.Println("io error getsitemaps ", er)
		}

		xml.Unmarshal(bod, &links)

	}
	return links
}

func fillQueue(links TheUrls, downIn chan string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()
	for _, link := range links.Locations {
		if link != "" && !visited[link] {
			visited[link] = true
			downIn <- link
		}
	}
	time.Sleep(1 * time.Second)
	mu.Unlock()
	close(downIn)
}

func crawlDb(url1 string, idx *dbSql, stopwords []string, wgc *sync.WaitGroup) {
	defer wgc.Done()
	downOut := make(chan dOut, 5000)
	extractOut := make(chan extrOut)
	downIn := make(chan string, 5000)
	downIn <- url1

	var mu sync.Mutex
	var wg sync.WaitGroup

	visited[url1] = true
	fmt.Println("crawling...")

	links := getSitemapLinks(url1)

	for {

		wg.Add(4)

		go download(downIn, downOut, &wg)

		go fillQueue(links, downIn, &mu, &wg)

		go extract(downOut, extractOut, &wg)

		go idx.Indexer(extractOut, stopwords, &wg)

		wg.Wait()

		break

	}

}
