package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/html"

	"crawler/config"
)

func main() {
	var configFile = flag.String("config", "config/config.json", "configuration file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if nCPU := runtime.NumCPU(); nCPU > 1 {
		runtime.GOMAXPROCS(nCPU)
	}

	cfg := config.Parse(*configFile)
	if cfg == nil {
		log.Println("Error while parsing config")
		os.Exit(1)
	}

	//Configuring global logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router := mux.NewRouter()
	router.HandleFunc("/crawl", crawl).Methods("GET")

	server := &http.Server{
		Addr:           cfg.HTTPAddress + ":" + strconv.Itoa(cfg.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 10,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		log.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	log.Println("Listening server on ", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	log.Println("Server gracefully stopped.")
}

var (
	sitemap map[string][]string
	siteMu  sync.RWMutex
)

type Crawler struct {
	SeedURL  *url.URL
	Visited  map[string]bool
	Climiter chan struct{}
	Queue    chan string
	Response chan *http.Response
	Sitemap  map[string][]string
	Quit     chan struct{}
}

func crawl(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	log.Println("Received URI: ", uri)

	cfg := config.Get()

	crawler := Crawler{
		Visited:  make(map[string]bool),
		Climiter: make(chan struct{}, cfg.Concurrency),
		Queue:    make(chan string, cfg.QueueSize),
		Response: make(chan *http.Response),
		Sitemap:  make(map[string][]string),
		Quit:     make(chan struct{}),
	}

	seed, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	crawler.SeedURL = seed

	// Initialising queue with given seed url
	go func() {
		crawler.Queue <- uri
	}()

	go crawler.Crawl()

	<-crawler.Quit

	for k := range crawler.Sitemap {
		fmt.Println(k)
		for _, v := range sitemap[k] {
			fmt.Printf("  - %s\n", v)
		}
	}
	log.Println("Killing the crawler...")
}

func (c *Crawler) Fetch(u string) {
	time.Sleep(2 * time.Second)
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic: ", err)
		}
		<-c.Climiter
	}()

	log.Println(u)

	client := http.Client{
		Timeout:  10 * time.Second,
	}

	resp, err := client.Get(u)
	if err != nil || resp == nil {
		log.Println("ERROR: Error occured", err)
		return
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.parseHTML(resp.Body)
	}

	resp.Body.Close()
}

func (c *Crawler) Crawl() {
	for link := range c.Queue {
		if alreadyVisited := c.Visited[link]; alreadyVisited {
			continue
		}
		c.Visited[link] = true
		c.Climiter <- struct{}{}
		log.Println("No of concurrency: ", len(c.Climiter))
		log.Println("Length of Queue: ", len(c.Queue))
		go c.Fetch(link)

		if len(c.Climiter) == 0 && len(c.Queue) == 0 {
			log.Println("Queue is empty")
			c.Quit <- struct{}{}
		}
	}
}

func (c *Crawler) parseHTML(reader io.Reader) {
	var (
		token html.TokenType
	)

	tracker := make(map[string]bool)
	tokenizer := html.NewTokenizer(reader)
	for {
		token = tokenizer.Next()
		switch token {
		case html.StartTagToken:
			t := tokenizer.Token()

			// Check if the token is an <a> tag
			if isAnchor := t.Data == "a"; !isAnchor {
				continue
			}

			for idx := range t.Attr {
				if t.Attr[idx].Key == "href" {
					// Making sure the href data is not empty
					if link := t.Attr[idx].Val; link != "" && link != "#" {
						abURL := absURL(link, c.SeedURL)
						if abURL == "" {
							continue
						}

						if isAlreadySent := tracker[abURL]; isAlreadySent {
							continue
						}
						tracker[abURL] = true

						c.Queue <- abURL
					}
				}
			}
		case html.ErrorToken:
			// End of the document
			return
		}
	}
}

func absURL(link string, base *url.URL) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Println("ERROR: ", err)
		return ""
	}

	uri := sanitiseURL(u)
	if uri == nil {
		return ""
	}

	// Restricting external links
	if uri.Host != "" && uri.Host != base.Host {
		return ""
	}
	return base.ResolveReference(uri).String()
}

// Basic sanitation for cleaning the URL
func sanitiseURL(uri *url.URL) *url.URL {
	// Only fetch HTTP/HTTPS links
	// Filtering links
	if uri.Scheme != "" && (uri.Scheme != "http" || uri.Scheme != "https") {
		return nil

	}

	uri.Fragment = ""
	if uri.Path == "/" {
		uri.Path = ""
	}

	if rURL := uri.String(); rURL == "" {
		return nil
	}

	return uri
}
