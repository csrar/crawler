package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/csrar/crawler/internal/models"
	"github.com/csrar/crawler/pkg/logger"
	"github.com/csrar/crawler/pkg/store"
	"golang.org/x/net/html"
)

type Crawler struct {
	ID       int
	page     *url.URL
	queue    chan string
	logger   logger.Ilogger
	workers  chan int
	finished chan int
	store    store.ICrawlerStore
}

//go:generate mockgen -source=crawler.go -destination=mocks/crawler_mock.go
type ICrawler interface {
	ExtractLinks() error
	SpinUpCrawler()
}

func NewCrawler(ID int, webPage string, channels *models.CommunitationChans, log logger.Ilogger, store store.ICrawlerStore) (ICrawler, error) {
	linkURL, err := url.Parse(webPage)

	if err != nil {
		return nil, fmt.Errorf("error parsing the provided URL: %v", err)
	}
	return &Crawler{
		ID:       ID,
		page:     linkURL,
		queue:    channels.Queue,
		logger:   log,
		workers:  channels.Workers,
		finished: channels.Finished,
		store:    store,
	}, nil
}

// ExtractLinks extracts links from a web page.
func (c *Crawler) ExtractLinks() error {
	defer c.returnWorker()

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

// returnWorker signals that a worker has finished.
func (c Crawler) returnWorker() {
	c.workers <- c.ID
	c.finished <- 1
}

// SpinUpCrawler initiates the crawling process.
func (c *Crawler) SpinUpCrawler() {
	err := c.ExtractLinks()
	if err != nil {
		c.logger.Error(err)
	}
}

// extractTagLink extracts links from HTML tokens.
func (c *Crawler) extractTagLink(token html.Token) (*string, error) {
	var link *string
	if token.Data == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" {
				if linkURL, err := c.parseURL(attr.Val); err == nil {
					if c.checkURL(linkURL) {
						visited, err := c.store.WasAlreadyVisited(linkURL.String())
						if err != nil {
							return nil, err
						}
						if !visited {
							c.logger.Info(fmt.Sprintf("worker: %d - Found link: %s", c.ID, linkURL.String()))
							tmp := linkURL.String()
							link = &tmp
							break
						}
					}

				} else {
					c.logger.Error(fmt.Errorf("error worker: %d - found invalid link:%s", c.ID, attr.Val))
				}
			}
		}
	}
	return link, nil
}

// parseURL parses a string URL into a *url.URL object.
func (c *Crawler) parseURL(strUrl string) (*url.URL, error) {
	url, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}
	if url.Host == "" && url.Scheme == "" {
		url.Host = c.page.Host
		url.Scheme = c.page.Scheme
	}

	return url, nil
}

// checkURL checks if a URL is valid for crawling.
func (c *Crawler) checkURL(url *url.URL) bool {
	if !url.IsAbs() {
		return false
	}
	if url.Path == "/" || url.Path == "" {
		return false
	}
	if url.Host != c.page.Host {
		return false
	}

	return true
}
