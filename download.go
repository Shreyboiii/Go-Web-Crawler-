package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func download(downIn chan string, downOut chan dOut, wg *sync.WaitGroup) {
	pols := NewPolicies()
	defer wg.Done()

	for link := range downIn {

		if pols.Disallowed(link) {
			continue
		}

		if dur := pols.Delay(link); dur > 0 {
			time.Sleep(dur * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
		}

		resp, err := http.Get(link)
		if err != nil {
			fmt.Println("func downlaod; Http.get err: ", err)
		}
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("func downlaod; Readall err: ", err)
		}

		downOut <- dOut{
			bts,
			link,
		}
	}

	close(downOut)

}
