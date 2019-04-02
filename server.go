package main

import (
	"context"
	"flag"
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

	cfg, err := config.Parse(*configFile)
	if err != nil {
		log.Println("Error while parsing config")
		os.Exit(2)
	}

	//Configuring global logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router := mux.NewRouter()
	router.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	router.HandleFunc("/crawl", handler.CrawlHandler).Methods("GET")

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