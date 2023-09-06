package crawler

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/csrar/crawler/internal/models"
	mock_logger "github.com/csrar/crawler/pkg/logger/mocks"
	mock_store "github.com/csrar/crawler/pkg/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCheckURL(t *testing.T) {
	// Define test cases as a table.
	tests := []struct {
		name           string
		url            *url.URL
		pageHost       string
		expectedResult bool
	}{
		{
			name:           "TestCheckURL_AbsURL_Valid",
			url:            &url.URL{Scheme: "https", Host: "example.com", Path: "/path"},
			pageHost:       "example.com",
			expectedResult: true,
		},
		{
			name:           "TestCheckURL_AbsURL_DifferentHost",
			url:            &url.URL{Scheme: "https", Host: "example2.com", Path: "/path"},
			pageHost:       "example.com",
			expectedResult: false,
		},
		{
			name:           "TestCheckURL_AbsURL_RootPath",
			url:            &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
			pageHost:       "example.com",
			expectedResult: false,
		},
		{
			name:           "TestCheckURL_RelativePath",
			url:            &url.URL{Path: "relative/path"},
			pageHost:       "example.com",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Crawler instance with the specified page host.
			c := &Crawler{
				page: &url.URL{
					Host: tc.pageHost,
				},
			}

			// Call the checkURL method and compare the result.
			result := c.checkURL(tc.url)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestParseURL(t *testing.T) {
	// Define test cases as a table.
	tests := []struct {
		name             string
		strURL           string
		pageHost         string
		pageScheme       string
		expectedParseURL *url.URL
		expectedParseErr error
	}{
		{
			name:             "TestParseURL_ValidURL",
			strURL:           "https://example.com/path",
			pageHost:         "",
			pageScheme:       "",
			expectedParseURL: &url.URL{Scheme: "https", Host: "example.com", Path: "/path"},
			expectedParseErr: nil,
		},
		{
			name:             "TestParseURL_InvalidURL",
			strURL:           ":invalid",
			pageHost:         "",
			pageScheme:       "",
			expectedParseURL: nil,
			expectedParseErr: errors.New(`parse ":invalid": missing protocol scheme`),
		},
		{
			name:             "TestParseURL_RelativeURL",
			strURL:           "/relative/path",
			pageHost:         "example.com",
			pageScheme:       "https",
			expectedParseURL: &url.URL{Scheme: "https", Host: "example.com", Path: "/relative/path"},
			expectedParseErr: nil,
		},
		{
			name:             "TestParseURL_EmptyURL",
			strURL:           "",
			pageHost:         "example.com",
			pageScheme:       "https",
			expectedParseURL: &url.URL{Scheme: "https", Host: "example.com"},
			expectedParseErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Crawler instance with the specified page host and scheme.
			c := &Crawler{
				page: &url.URL{
					Host:   tc.pageHost,
					Scheme: tc.pageScheme,
				},
			}

			// Call the parseURL method and check the result.
			parsedURL, parseErr := c.parseURL(tc.strURL)

			// Use the testify/assert library for assertions.
			if parseErr != nil {
				assert.Equal(t, tc.expectedParseErr.Error(), parseErr.Error())
			}
			assert.Equal(t, tc.expectedParseURL, parsedURL)

		})
	}
}

func TestReturnWorker(t *testing.T) {
	// Define test cases as a table.
	tests := []struct {
		name             string
		crawler          Crawler
		expectedWorkerID int
	}{
		{
			name: "TestReturnWorker_Valid",
			crawler: Crawler{
				ID:       1,
				workers:  make(chan int, 1),
				finished: make(chan int, 1),
			},
			expectedWorkerID: 1,
		},
		{
			name: "TestReturnWorker_EmptyChannels",
			crawler: Crawler{
				ID:       2,
				workers:  make(chan int),
				finished: make(chan int),
			},
			expectedWorkerID: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the returnWorker method.
			go tc.crawler.returnWorker()

			// Check if the worker ID is received as expected.
			workerID := <-tc.crawler.workers
			assert.Equal(t, tc.expectedWorkerID, workerID)

			// Check if the finished channel receives a value.
			<-tc.crawler.finished
		})
	}
}

type extractLinksCh struct {
	queue    []string
	workers  []int
	finished []int
}

func TestExtractLinks(t *testing.T) {
	// Define test cases as a table.
	tests := []struct {
		name                             string
		mockHttpResponse                 string
		expectedbuildErr                 error
		expectedError                    error
		mockLogInfoCalls                 int
		mockLogErrorCalls                int
		mockStoreWasAlreadyVisitedCalls  int
		mockStoreWasAlreadyVisitedResult bool
		mockStoreWasAlreadyVisitedError  error
		prefixURL                        string
		ch                               *models.CommunitationChans
		expectedCH                       extractLinksCh
	}{
		{
			name:                             "Successfull link extraction",
			mockHttpResponse:                 `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Example HTML Page</title></head><body><header><h1>mock Website</h1><nav><ul><li><a href="%host%">Home</a></li><li><a href="%host%/about">About Us</a></li><li><a href="%host%/contact">Contact</a></li></ul></nav></header><section><h2>About Us</h2><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla eget risus eu purus efficitur ullamcorper.</p></section><section><h2>Our Services</h2><ul><li><a href="%host%/test">Service 1</a></li><li><a href="%host%/demo">Service 2</a></li><li><a href="%host%/foo">Service 3</a></li></ul></section><footer><p>&copy; 2023 My mock website. All rights reserved.</p></footer></body></html>`,
			expectedError:                    nil,
			mockLogInfoCalls:                 6,
			mockStoreWasAlreadyVisitedCalls:  5,
			mockStoreWasAlreadyVisitedResult: false,
			ch: &models.CommunitationChans{
				Queue:    make(chan string, 5),
				Workers:  make(chan int, 1),
				Finished: make(chan int, 1),
			},
			expectedCH: extractLinksCh{
				workers:  []int{1},
				finished: []int{1},
				queue:    []string{"%host%/about", "%host%/contact", "%host%/test", "%host%/demo", "%host%/foo"},
			},
		},
		{
			name:                             "fail getting error from was alreadyVisited",
			mockHttpResponse:                 `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Example HTML Page</title></head><body><header><h1>mock Website</h1><nav><ul><li><a href="%host%">Home</a></li><li><a href="%host%/about">About Us</a></li><li><a href="%host%/contact">Contact</a></li></ul></nav></header><section><h2>About Us</h2><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla eget risus eu purus efficitur ullamcorper.</p></section><section><h2>Our Services</h2><ul><li><a href="%host%/test">Service 1</a></li><li><a href="%host%/demo">Service 2</a></li><li><a href="%host%/foo">Service 3</a></li></ul></section><footer><p>&copy; 2023 My mock website. All rights reserved.</p></footer></body></html>`,
			expectedError:                    nil,
			mockLogInfoCalls:                 1,
			mockLogErrorCalls:                1,
			mockStoreWasAlreadyVisitedCalls:  1,
			mockStoreWasAlreadyVisitedResult: false,
			mockStoreWasAlreadyVisitedError:  errors.New("mock-error"),
			ch: &models.CommunitationChans{
				Queue:    make(chan string, 5),
				Workers:  make(chan int, 1),
				Finished: make(chan int, 1),
			},
			expectedCH: extractLinksCh{
				workers:  []int{1},
				finished: []int{1},
				queue:    []string{},
			},
		},
		{
			name:                             "successfult with invalid link",
			mockHttpResponse:                 `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Example HTML Page</title></head><body><header><h1>mock Website</h1><nav><ul><li><a href="%host%">Home</a></li><li><a href="%host%/about">About Us</a></li><li><a href="www.mock.com%%2">Contact</a></li></ul></nav></header><section><h2>About Us</h2><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla eget risus eu purus efficitur ullamcorper.</p></section><section><h2>Our Services</h2><ul><li><a href="%host%/test">Service 1</a></li><li><a href="%host%/demo">Service 2</a></li><li><a href="%host%/foo">Service 3</a></li></ul></section><footer><p>&copy; 2023 My mock website. All rights reserved.</p></footer></body></html>`,
			expectedError:                    nil,
			mockLogInfoCalls:                 5,
			mockLogErrorCalls:                1,
			mockStoreWasAlreadyVisitedCalls:  4,
			mockStoreWasAlreadyVisitedResult: false,
			ch: &models.CommunitationChans{
				Queue:    make(chan string, 5),
				Workers:  make(chan int, 1),
				Finished: make(chan int, 1),
			},
			expectedCH: extractLinksCh{
				workers:  []int{1},
				finished: []int{1},
				queue:    []string{"%host%/about", "%host%/test", "%host%/demo", "%host%/foo"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test HTTP server.
			l, err := net.Listen("tcp", "127.0.0.1:")
			host := fmt.Sprintf("http://%s", l.Addr().String())
			handler := func(w http.ResponseWriter, r *http.Request) {
				// Set response headers.
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)

				// Write the static HTML content to the response.
				htmlContent := strings.ReplaceAll(tc.mockHttpResponse, "%host%", host)
				fmt.Fprintln(w, htmlContent)
			}

			if err != nil {
				log.Fatal(err)
			}

			// Create and configure the test HTTP server.
			testServer := httptest.NewUnstartedServer(http.HandlerFunc(handler))
			testServer.Listener.Close()
			testServer.Listener = l

			// Start the server.
			testServer.Start()
			// Stop the server on return from the function.
			defer testServer.Close()

			// Create and configure mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logMock := mock_logger.NewMockIlogger(ctrl)
			logMock.EXPECT().Info(gomock.Any()).Times(tc.mockLogInfoCalls)
			logMock.EXPECT().Error(gomock.Any()).Times(tc.mockLogErrorCalls)

			storeMock := mock_store.NewMockICrawlerStore(ctrl)
			storeMock.EXPECT().WasAlreadyVisited(gomock.Any()).Return(tc.mockStoreWasAlreadyVisitedResult, tc.mockStoreWasAlreadyVisitedError).Times(tc.mockStoreWasAlreadyVisitedCalls)

			// Create and run the crawler.
			crawler, _ := NewCrawler(1, testServer.URL, tc.ch, logMock, storeMock)

			crawler.SpinUpCrawler()
			close(tc.ch.Finished)
			close(tc.ch.Workers)
			close(tc.ch.Queue)

			// Compare expected and actual results.
			expectedqueue := []string{}
			for _, row := range tc.expectedCH.queue {
				expectedqueue = append(expectedqueue, strings.Replace(row, "%host%", host, 1))

			}
			resultQueue := []string{}
			for queue := range tc.ch.Queue {
				resultQueue = append(resultQueue, queue)
			}
			assert.Equal(t, expectedqueue, resultQueue)

			resultWorker := []int{}
			for worker := range tc.ch.Workers {
				resultWorker = append(resultWorker, worker)
			}
			assert.Equal(t, tc.expectedCH.workers, resultWorker)

			resultFinished := []int{}
			for finished := range tc.ch.Finished {
				resultFinished = append(resultFinished, finished)
			}
			assert.Equal(t, tc.expectedCH.finished, resultFinished)
		})
	}
}

func BenchmarkSpinUpCrawler(b *testing.B) {
	b.StopTimer()

	l, err := net.Listen("tcp", "127.0.0.1:")
	host := fmt.Sprintf("http://%s", l.Addr().String())
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set response headers.
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		// Write the static HTML content to the response.
		body := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Example HTML Page</title></head><body><header><h1>mock Website</h1><nav><ul><li><a href="%host%">Home</a></li><li><a href="%host%/about">About Us</a></li><li><a href="%host%/contact">Contact</a></li></ul></nav></header><section><h2>About Us</h2><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla eget risus eu purus efficitur ullamcorper.</p></section><section><h2>Our Services</h2><ul><li><a href="%host%/test">Service 1</a></li><li><a href="%host%/demo">Service 2</a></li><li><a href="%host%/foo">Service 3</a></li></ul></section><footer><p>&copy; 2023 My mock website. All rights reserved.</p></footer></body></html>`
		htmlContent := strings.ReplaceAll(body, "%host%", host)
		fmt.Fprintln(w, htmlContent)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Create and configure the test HTTP server.
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(handler))
	testServer.Listener.Close()
	testServer.Listener = l

	// Start the server.
	testServer.Start()
	// Create and configure mock objects
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	logMock := mock_logger.NewMockIlogger(ctrl)
	logMock.EXPECT().Info(gomock.Any()).AnyTimes()
	logMock.EXPECT().Error(gomock.Any()).AnyTimes()

	storeMock := mock_store.NewMockICrawlerStore(ctrl)
	storeMock.EXPECT().WasAlreadyVisited(gomock.Any()).Return(false, nil).AnyTimes()

	ch := &models.CommunitationChans{
		Queue:    make(chan string, 5),
		Workers:  make(chan int, 1),
		Finished: make(chan int, 1),
	}

	// Create and run the crawler.
	b.StartTimer()
	crawler, _ := NewCrawler(1, testServer.URL, ch, logMock, storeMock)
	crawler.SpinUpCrawler()
}
