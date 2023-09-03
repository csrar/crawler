package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/csrar/crawler/pkg/logger"
	"github.com/csrar/crawler/pkg/store"
	"golang.org/x/net/html"
)

type Crawler struct {
	ID      int
	page    *url.URL
	wg      *sync.WaitGroup
	queue   chan string
	logger  logger.Ilogger
	workers chan int
	store   store.ICrawlerStore
}

//go:generate mockgen -source=crawler.go -destination=mocks/crawler_mock.go
type ICrawler interface {
	ExtractLinks() error
	SpinUpCrawler()
}

func NewCrawler(ID int, webPage string, waitGroup *sync.WaitGroup,
	queue chan string, log logger.Ilogger, workers chan int, store store.ICrawlerStore) (ICrawler, error) {
	linkURL, err := url.Parse(webPage)

	if err != nil {
		return nil, fmt.Errorf("error parsing the provided URL: %v", err)
	}
	return &Crawler{
		ID:      ID,
		page:    linkURL,
		wg:      waitGroup,
		queue:   queue,
		logger:  log,
		workers: workers,
		store:   store,
	}, nil
}

func (c *Crawler) ExtractLinks() error {
	defer c.returnWorker()
	defer c.wg.Done()

	pageBody, err := http.Get(c.page.String())
	if err != nil {
		return fmt.Errorf("worker: %d - error visiting page: %s", c.ID, err)
	}
	c.logger.Info(fmt.Sprintf("worker: %d - visiting page: %s", c.ID, c.page))
	tokenizer := html.NewTokenizer(pageBody.Body)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, finish method
				return nil
			}
			return fmt.Errorf("worker: %d error tokenizing HTML: %v", c.ID, tokenizer.Err())
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			link, err := c.extractTagLink(token)
			if err != nil {
				return err
			}
			if link != nil {
				c.queue <- *link
			}
		}
	}
}

func (c Crawler) returnWorker() {
	c.workers <- c.ID
}

func (c *Crawler) SpinUpCrawler() {
	// c.wg.Add(1)
	err := c.ExtractLinks()
	if err != nil {
		c.logger.Error(err)
	}
}

func (c *Crawler) extractTagLink(token html.Token) (*string, error) {
	var link *string
	if token.Data == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" {
				if linkURL, err := url.Parse(attr.Val); err == nil {
					if c.checkURL(linkURL) {
						c.logger.Info(fmt.Sprintf("worker: %d - Found link: %s", c.ID, linkURL.String()))
						tmp := linkURL.String()
						link = &tmp
						break
					}

				} else {
					fmt.Println("some error")
				}
			}
		}
	}
	return link, nil
}

func (c *Crawler) checkURL(url *url.URL) bool {
	if !url.IsAbs() {
		return false
	}
	if url.Path == "/" || url.Path == "" {
		return false
	}

	return true
}
