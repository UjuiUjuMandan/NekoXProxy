package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func relay(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	reqip := strings.Split(u.Host, ":")[0]
	dc := ip2dc(reqip)
	wsurl := dc2wsurl(dc)

	log.Printf("%s %s -> DC%d %s", r.Method, reqip, dc, wsurl)

	if wsurl == "" {
		log.Printf("no relay for %s", reqip)
		w.WriteHeader(502)
		return
	}

	var body io.Reader
	if r.Method == "POST" {
		body = r.Body
	}

	req, _ := http.NewRequest(r.Method, wsurl, body)
	a, err := client.Do(req)
	if err != nil {
		log.Printf("upstream error: %v", err)
		w.WriteHeader(502)
		return
	}

	for k := range a.Header {
		w.Header().Set(k, a.Header.Get(k))
	}

	log.Printf("response %d <- DC%d %s", a.StatusCode, dc, wsurl)
	w.WriteHeader(a.StatusCode)
	io.Copy(w, a.Body)
}
