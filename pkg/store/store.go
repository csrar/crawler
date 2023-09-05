package store

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/csrar/crawler/internal/models"
	"github.com/gosimple/slug"
)

//go:generate mockgen -source=store.go -destination=mocks/store_mock.go
type ICrawlerStore interface {
	WasAlreadyVisited(site string) (bool, error)
	StoreData(site string, data models.SiteStore) error
}

type store struct {
	visited ImFile
	mu      *sync.Mutex
}

func NewMemfileStore(mu *sync.Mutex, file ImFile) ICrawlerStore {
	return &store{
		visited: file,
		mu:      mu,
	}
}

func (s store) WasAlreadyVisited(site string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slug_site := slug.Make(site)
	visited, err := s.readData()
	if err != nil {
		return false, err
	}
	if _, ok := visited.Sites[slug_site]; !ok {
		if err := s.StoreData(slug_site, *visited); err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s store) StoreData(site string, data models.SiteStore) error {
	data.Sites[site] = true
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling visited data %+v", err)
	}
	err = s.visited.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating visited data %+v", err)
	}
	_, err = s.visited.WriteAt(payload, 0)
	if err != nil {
		return fmt.Errorf("error writing visited data %+v", err)
	}
	return nil
}

func (s store) readData() (*models.SiteStore, error) {
	message := &models.SiteStore{}
	err := json.Unmarshal(s.visited.Bytes(), message)
	if err != nil {
		return nil, fmt.Errorf("error decoding unmarshaling data: %+v", err)
	}
	return message, nil
}
