package boot

import (
	"errors"
	"net/url"
	"sync"

	"github.com/csrar/crawler/internal/models"
	"github.com/csrar/crawler/pkg/config"
	"github.com/csrar/crawler/pkg/store"
)

type Ibootstrap interface {
	BoostrapStore() (store.ICrawlerStore, error)
	BootsRootPage() (*url.URL, error)
	BootstrapChannels() *models.CommunitationChans
	StartWorkersQueue(workers chan int)
}

type boot struct {
	config config.IConfig
}

// NewBootstrap creates a new bootstrapper instance.
func NewBootstrap(config config.IConfig) Ibootstrap {
	return &boot{
		config: config,
	}
}

// BoostrapStore initializes and configures the crawler store.
func (b boot) BoostrapStore() (store.ICrawlerStore, error) {
	store := store.NewMemfileStore(&sync.Mutex{})
	err := store.StoreData("", models.SiteStore{
		Sites: map[string]bool{},
	})
	if err != nil {
		return nil, err
	}
	return store, nil
}

// BootsRootPage parses and validates the root page URL from the configuration.
func (b boot) BootsRootPage() (*url.URL, error) {
	page, err := url.Parse(b.config.GetConfig().WepPage)
	if err != nil {
		return nil, errors.New("invalid site URL provided")
	}
	return page, nil
}

// BootstrapChannels creates communication channels for the crawler.
func (b boot) BootstrapChannels() *models.CommunitationChans {
	return &models.CommunitationChans{
		Queue:    make(chan string, 100000),
		Workers:  make(chan int, b.config.GetConfig().Workers),
		Finished: make(chan int),
	}
}

// StartWorkersQueue initializes and starts worker queue with the specified number of workers.
func (b boot) StartWorkersQueue(workers chan int) {
	for i := 0; i < b.config.GetConfig().Workers; i++ {
		workers <- i + 1
	}
}
