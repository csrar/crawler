package models

type Config struct {
	Workers   int
	WepPage   string
	QueueSize int
}

type SiteStore struct {
	Sites map[string]bool `json:"sites"`
}

type CommunitationChans struct {
	Queue    chan string
	Workers  chan int
	Finished chan int
}
