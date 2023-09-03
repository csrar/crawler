package store

import (
	"encoding/json"
	"fmt"

	"github.com/csrar/crawler/internal/models"
	"github.com/dsnet/golib/memfile"
)

//go:generate mockgen -source=store.go -destination=mocks/store_mock.go
type ICrawlerStore interface {
	WasAlreadyVisited(site string) (bool, error)
}

type store struct {
	visited ImFile
}

func NewMemfileStore() ICrawlerStore {
	return &store{
		visited: memfile.New([]byte{}),
	}
}

func (s store) WasAlreadyVisited(site string) (bool, error) {
	visited, err := s.readData()
	if err != nil {
		return false, err
	}
	if _, ok := visited.Sites[site]; !ok {
		if err := s.storeData(site, *visited); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s store) storeData(site string, data models.SiteStore) error {
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
