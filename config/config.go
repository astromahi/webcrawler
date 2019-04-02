package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

// Config holds configuration options
type Config struct {
	HTTPAddress  string `json:"http_address"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	Concurrency  int    `json:"concurrency"`
	FetchTimeout int    `json:"fetch_timeout"`
	QueueSize    int    `json:"queue_size"`
	CrawlDelay   int    `json:"crawl_delay"`
	Depth        int    `json:"depth"`
}

var (
	cfg  Config
	once sync.Once
)

// Parse parses the json configuration file
// And converting it into native type
func Parse(file string) (*Config, error) {
	if file == "" {
		return nil, errors.New("config: Given file name is empty")
	}

	var errMain error

	once.Do(func() {
		// Reading the flags
		data, err := ioutil.ReadFile(file)
		if err != nil {
			errMain = err
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			errMain = err
		}
	})

	if errMain != nil {
		return nil, errMain
	}

	return &cfg, nil
}

// Getter for configuration object
func Get() *Config {
	return &cfg
}
