package main

import (
	"errors"
	"fmt"
	"sync"

	boot "github.com/csrar/crawler/internal/bootstrap"
	"github.com/csrar/crawler/internal/service"
	"github.com/csrar/crawler/pkg/config"
	"github.com/csrar/crawler/pkg/logger"
)

func main() {
	// Initialize logger, config, and bootstrapper
	log := logger.NewLogrusLogger()
	config := config.NewConfig()
	boot := boot.NewBootstrap(config)

	// Bootstrap root page
	page, err := boot.BootsRootPage()
	if err != nil {
		log.Error(err)
		return
	}

	// Bootstrap store and handle errors
	store, err := boot.BoostrapStore()
	if err != nil {
		log.Error(err)
		return
	}

	// Bootstrap channels
	channels := boot.BootstrapChannels()
	channels.Queue <- page.String()

	boot.StartWorkersQueue(channels.Workers)

	if err != nil {
		log.Error(errors.New("invalid site url provided"))
		return
	}

	// Initialize counters and synchronization objects
	found := 0
	processed := 0
	var mx sync.Mutex
	var wg sync.WaitGroup

	// Add a WaitGroup for tracking the completion of crawling
	wg.Add(1)

	// Create a CrawlerHandler to manage crawling
	handler := service.NewCrawlerHandler(&found, &processed, channels, log, store, &wg, &mx)

	// Start Goroutines for listening to new links and validating crawl finish
	go handler.ListenForNewLinks()
	go handler.ValidateCrawlFinish()

	wg.Wait()
	log.Info(fmt.Sprintf("finished crawling for %s, total liks explored: %d", page.String(), processed))

}
