package handler

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"webcrawler/config"
	"webcrawler/crawler"
)

// HTTP handler for activating crawl
func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("INFO: Received crawl request for - ", uri)
	cfg := config.Get()

	bot := crawler.New(cfg)

	// Custom fetcher function
	fetcher := func(link string) (io.ReadCloser, error) {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}

		// A HTTP client with timeout
		client := http.Client{
			Transport: transport,
			Timeout:   time.Duration(cfg.FetchTimeout) * time.Second,
		}

		// Validating response
		resp, err := client.Get(link)
		if err != nil {
			log.Println("ERROR: Couldn't fetch - ", err)
			return nil, err
		} else if resp == nil {
			return nil, errors.New("Fetcher: empty response received")
		} else if resp.StatusCode >= 300 {
			return nil, errors.New("Fetcher: http status code is >= 300")
		}

		return resp.Body, nil
	}
	bot.FetcherFunc = fetcher

	// Deploying the crawler
	go bot.Crawl(uri)

	<-bot.Done

	sitemap := struct {
		Index map[string][]string `json:"sitemap"`
	}{
		bot.Sitemap,
	}

	if err := json.NewEncoder(w).Encode(&sitemap); err != nil {
		return
	}
}

// A HTTP ping handler
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "I'am healthy"}`))
}
