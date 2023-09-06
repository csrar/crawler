package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/csrar/crawler/internal/models"
)

type config struct {
	cfg models.Config
}

//go:generate mockgen -source=config.go -destination=mocks/config_mock.go
type IConfig interface {
	GetConfig() models.Config
}

func NewConfig() IConfig {
	cfg := models.Config{}
	cfg.WepPage = getStringVal(keyWebPage, defaultWebPage)
	cfg.Workers = getIntValue(keyWorkers, defaultWorkers)
	cfg.QueueSize = detaultQueueSize
	return &config{
		cfg: cfg,
	}
}

func (c *config) GetConfig() models.Config {
	return c.cfg
}

func getStringVal(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getIntValue(key string, def int) int {
	returnValue := 0
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	returnValue, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("invalid integer value for %s, app will will use default value: %d\n", key, def)
		return def
	}
	return returnValue
}
