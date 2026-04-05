package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func relay(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	reqip := host
	relayURL := ip2url(reqip)

	log.Printf("%s %s -> %s", r.Method, reqip, relayURL)

	if relayURL == "" {
		log.Printf("no relay for %s", reqip)
		w.WriteHeader(502)
		return
	}

	var body io.Reader
	if r.Method == "POST" {
		body = r.Body
	}

	req, _ := http.NewRequest(r.Method, relayURL, body)
	a, err := client.Do(req)
	if err != nil {
		log.Printf("upstream error: %v", err)
		w.WriteHeader(502)
		return
	}

	for k := range a.Header {
		w.Header().Set(k, a.Header.Get(k))
	}

	log.Printf("response %d <- %s", a.StatusCode, relayURL)
	w.WriteHeader(a.StatusCode)
	io.Copy(w, a.Body)
}
