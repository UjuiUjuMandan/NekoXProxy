package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var nekoXProxyBaseDomain string
var nekoXProxyDomains []string

var client *http.Client

func main() {
	listen := flag.String("l", "127.0.0.1:26641", "HttpProxy listen port")
	_nekoXProxyString := flag.String("p", "", "NekoX Proxy URL")
	flag.Parse()

	if *_nekoXProxyString == "" {
		log.Fatalln("Relay address is required. Use -p to specify the NekoX Proxy URL.")
	}

	if !parseNekoXString(*_nekoXProxyString) {
		log.Fatalln("Failed to parse NekoX proxy.")
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
