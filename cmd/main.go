package main

import (
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/csrar/crawler/pkg/config"
	"github.com/csrar/crawler/pkg/crawler"
	"github.com/csrar/crawler/pkg/logger"
	"github.com/csrar/crawler/pkg/store"
)

func startWorkersQueue(maxWorkers int, workers chan int) {
	for i := 0; i < maxWorkers; i++ {
		workers <- i + 1
	}
}

func main() {
	log := logger.NewLogrusLogger()
	var wg sync.WaitGroup
	queue := make(chan string, 1000)

	maxWorkers := config.NewConfig().GetConfig().Workers
	workers := make(chan int, maxWorkers)
	startWorkersQueue(maxWorkers, workers)
	store := store.NewMemfileStore()
	config := config.NewConfig()

	page, err := url.Parse(config.GetConfig().WepPage)
	if err != nil {
		log.Error(errors.New("invalid site url provided"))
		return
	}

	go func() {
		for {
			select {
			case link := <-queue:
				workerID := <-workers
				crawl, err := crawler.NewCrawler(workerID, link, &wg, queue, log, workers, store)
				if err != nil {
					log.Error(err)
					return
				}
				wg.Add(1)
				go crawl.SpinUpCrawler()
			}
		}
	}()
	queue <- page.String()
	time.Sleep(time.Second * 1)
	wg.Wait()

}
