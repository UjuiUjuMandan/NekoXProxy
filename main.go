package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var client *http.Client

func main() {
	listen := flag.String("l", "127.0.0.1:26641", "HttpProxy listen port")
	_nekoXProxyString := flag.String("p", "", "Relay base URL (e.g. ws://mtproto.example.com)")
	configFile := flag.String("c", "dcmap.json", "IP-to-subdomain mapping config file")
	flag.Parse()

	if *_nekoXProxyString == "" {
		log.Fatalln("Relay address is required. Use -p to specify the relay base URL.")
	}

	if err := loadMapper(*configFile); err != nil {
		log.Fatalf("Failed to load config %s: %v", *configFile, err)
	}

	if !parseNekoXString(*_nekoXProxyString) {
		log.Fatalln("Failed to parse relay URL.")
	}

	client = &http.Client{}

	http.HandleFunc("/", relay)
	server := &http.Server{
		Addr:         *listen,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Telegram HTTP Proxy started at", *listen)
	server.ListenAndServe()
}
