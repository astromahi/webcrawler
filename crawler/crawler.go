package crawler

import (
	"errors"
	"io"
	"log"
	"net/url"
	"time"

	"webcrawler/config"

	"golang.org/x/net/html"
)

// Crawler object holding the internal state
type Crawler struct {
	seedURL     *url.URL
	concurrency chan struct{}
	queue       chan string
	crawldelay  int
	depth       int
	timeout     int
	visited     map[string]bool

	FetcherFunc func(string) (io.ReadCloser, error)
	Sitemap     map[string][]string
	Quit        chan struct{}
}

// New generates a brand new crawler object
func New(c *config.Config) *Crawler {
	// Default settings
	concurrency := c.Concurrency
	if concurrency < 1 {
		concurrency = 1
	}

	cdelay := c.CrawlDelay
	if cdelay < 1 {
		cdelay = 200
	}

	timeout := c.FetchTimeout
	if timeout < 1 {
		timeout = 10
	}

	queueCap := c.QueueSize
	if queueCap < 1000 {
		queueCap = 1000
	}

	return &Crawler{
		concurrency: make(chan struct{}, int(concurrency)),
		queue:       make(chan string, int(queueCap)),
		crawldelay:  cdelay,
		depth:       c.Depth,
		timeout:     timeout,
		visited:     make(map[string]bool),
		Sitemap:     make(map[string][]string),
		Quit:        make(chan struct{}),
	}
}

// Crawl intiates the link crawling
// Control loop that collecting incoming links
func (c *Crawler) Crawl(seed string) error {
	var err error
	if c.seedURL, err = url.Parse(seed); err != nil {
		log.Println("ERROR: Couldn't parse seed url - ", err)
		c.Quit <- struct{}{}
		return err
	}

	if c.seedURL.Hostname() == "" {
		err = errors.New("Invalid domain")
		log.Println("ERROR: ", err)
		c.Quit <- struct{}{}
		return err
	}

	// Initialising queue with seed URL
	c.queue <- seed

	// Setting crawl delay
	ticker := time.NewTicker(time.Duration(c.crawldelay) * time.Millisecond)

	for link := range c.queue {
		if alreadyVisited := c.visited[link]; alreadyVisited {
			continue
		}
		<-ticker.C

		log.Println("INFO: Crawling - ", link)

		c.visited[link] = true
		c.concurrency <- struct{}{}
		go c.run(link)
	}

	return nil
}

// Connecting fetcher & parser
func (c *Crawler) run(u string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic: ", err)
		}
		<-c.concurrency

		if len(c.concurrency) == 0 && len(c.queue) == 0 {
			c.Quit <- struct{}{}
		}
	}()

	resp, err := c.FetcherFunc(u)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}

	links := c.parseHTML(resp)
	for idx := range links {
		c.queue <- links[idx]
		if links[idx] != u {
			c.Sitemap[u] = append(c.Sitemap[u], links[idx])
		}
	}

	resp.Close()
	return
}

// parseHTML is a parser that parse incoming HTML data
func (c *Crawler) parseHTML(reader io.Reader) []string {
	var (
		links []string
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
						abURL := absURL(link, c.seedURL)
						if abURL == "" {
							continue
						}

						// A smart tracker to handle the duplicate while parsing the
						// HTML DOM
						if isAlreadySent := tracker[abURL]; isAlreadySent {
							continue
						}
						tracker[abURL] = true
						links = append(links, abURL)
					}
				}
			}
		case html.ErrorToken:
			// End of the document
			return links
		}
	}

	return links
}

func absURL(link string, base *url.URL) string {
	u, err := url.Parse(link)
	if err != nil {
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
	// Filtering links
	// Only fetch HTTP/HTTPS links
	/*
		if uri.Scheme != "" && (uri.Scheme != "http" || uri.Scheme != "https") {
			return nil
		}
	*/

	// Cleaning empty fragment & path in URL
	uri.Fragment = ""
	if uri.Path == "/" {
		uri.Path = ""
	}

	if rURL := uri.String(); rURL == "" {
		return nil
	}

	return uri
}
