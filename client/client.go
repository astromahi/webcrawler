// A Client utility for contacting server
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	var uri string
	flag.StringVar(&uri, "url", "", "URL to crawl")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(*uri) < 1 {
		fmt.Println("ERROR: Given URL is empty")
		return
	}

	if rawURL, err := url.Parse(uri); err != nil || rawURL.Hostname() == "" {
		fmt.Printf("ERROR: Hostname(%s) is not valid - %s\n", rawURL.Hostname(), err)
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

	serverURL := fmt.Sprintf("http://127.0.0.1:8080/crawl?uri=%s", uri)

	// Validating response
	resp, err := client.Get(serverURL)
	defer resp.Body.Close()
	if err != nil || resp == nil {
		fmt.Println("ERROR: Couldn't fetch - ", err)
		return
	}

	if resp.StatusCode >= 300 {
		fmt.Println("ERROR: HTTP status code is not success")
		return
	}

	var sitemap = struct {
		Data map[string][]string `json:"sitemap"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		fmt.Println("ERROR: Error while parsing response - ", err)
	}

	// Printing sitemap
	for key, val := range sitemap.Data {
		fmt.Println(key)
		for idx := range val {
			fmt.Printf("  - %s\n", val[idx])
		}
	}

	return
}
