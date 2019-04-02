// A Client utility for contacting server
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var uri = flag.String("url", "", "URL to crawl")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(*uri) < 1 {
		log.Println("ERROR: Given URL is empty")
		return
	}

	// A HTTP client with timeout
	client := http.Client{
		Timeout: 180 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	serverURL := `http://127.0.0.1:8080/crawl?uri="` + *uri + `"`

	// Validating response
	resp, err := client.Get(serverURL)
	defer resp.Body.Close()
	if err != nil || resp == nil {
		log.Println("ERROR: Couldn't fetch - ", err)
		return
	}

	if resp.StatusCode >= 300 {
		log.Println("ERROR: HTTP status code is not success")
		return
	}

	var sitemap = struct {
		Data map[string][]string `json:"sitemap"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		log.Println("ERROR: Error while parsing response - ", err)
	}

	// Printing sitemap
	for key, val := range sitemap.Data {
		log.Println(key)
		for idx := range val {
			log.Printf("  - %s\n", val[idx])
		}
	}

	return
}
