package service

import (
	"sync"

	"github.com/csrar/crawler/internal/models"
	"github.com/csrar/crawler/pkg/crawler"
	"github.com/csrar/crawler/pkg/logger"
	"github.com/csrar/crawler/pkg/store"
)

// crawlerHandler handles the crawling process, tracks found and processed links, and manages synchronization.
type crawlerHandler struct {
	found     *int
	processed *int
	channels  *models.CommunitationChans
	log       logger.Ilogger
	store     store.ICrawlerStore
	wg        *sync.WaitGroup
	mx        *sync.Mutex
}

type ICrawlerHandler interface {
	ListenForNewLinks()
	ValidateCrawlFinish()
}

// NewCrawlerHandler creates a new crawlerHandler instance.
func NewCrawlerHandler(found *int, processed *int, channels *models.CommunitationChans, log logger.Ilogger,
	store store.ICrawlerStore, wg *sync.WaitGroup, mx *sync.Mutex) ICrawlerHandler {
	return &crawlerHandler{
		found:     found,
		processed: processed,
		channels:  channels,
		log:       log,
		store:     store,
		wg:        wg,
		mx:        mx,
	}
}

// ListenForNewLinks listens for new links in the queue and initiates crawling.
func (c *crawlerHandler) ListenForNewLinks() {
	for {
		select {
		case link := <-c.channels.Queue:
			c.mx.Lock()
			*c.found++
			c.mx.Unlock()
			workerID := <-c.channels.Workers
			crawl, err := crawler.NewCrawler(workerID, link, c.channels, c.log, c.store)
			if err != nil {
				c.log.Error(err)
				return
			}
			go crawl.SpinUpCrawler()
		}
	}
}

// ValidateCrawlFinish validates if all found links have been processed and signals WaitGroup.
func (c *crawlerHandler) ValidateCrawlFinish() {
	for {
		select {
		case <-c.channels.Finished:
			*c.processed++
		}
		c.mx.Lock()
		if *c.found == *c.processed {
			c.wg.Done()
		}
		c.mx.Unlock()

	}
}
