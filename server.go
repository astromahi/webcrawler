// Webcrawler server that takes requests from the clients and
// responding with the sitemap as response
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"webcrawler/config"
	"webcrawler/handler"

	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "config/config.json", "configuration file")
	flag.Parse()

	if nCPU := runtime.NumCPU(); nCPU > 1 {
		runtime.GOMAXPROCS(nCPU)
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Println("Error while parsing config")
		os.Exit(2)
	}

	//Configuring global logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	router.HandleFunc("/crawl", handler.CrawlHandler).Methods("GET")

	// Creating server object
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
	if err = server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	log.Println("Server gracefully stopped.")
}
