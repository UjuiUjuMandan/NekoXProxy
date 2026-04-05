package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

var mapper = make(map[string]string) // IP (or prefix) → subdomain
var proxyScheme = "https"
var nekoXProxyBaseDomain string

func parseNekoXString(a string) bool {
	if a == "" {
		return false
	}

	u, err := url.Parse(a)
	if err != nil {
		return false
	}

	nekoXProxyBaseDomain = u.Host

	proxyScheme = u.Scheme

	fmt.Printf("Base domain: %s\nScheme: %s\n", nekoXProxyBaseDomain, proxyScheme)
	return true
}

func loadMapper(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &mapper); err != nil {
		return err
	}
	fmt.Printf("Loaded %d IP mappings from %s\n", len(mapper), path)
	return nil
}

func ip2url(ip string) string {
	sub := ip2subdomain(ip)
	if sub == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s.%s/api", proxyScheme, sub, nekoXProxyBaseDomain)
}

func ip2subdomain(ip string) string {
	parsed := net.ParseIP(ip)
	for k, v := range mapper {
		if ip == k {
			return v
		}
		if parsed != nil {
			if kp := net.ParseIP(k); kp != nil && parsed.Equal(kp) {
				return v
			}
		}
	}
	for k, v := range mapper {
		if strings.HasPrefix(ip, k) {
			return v
		}
	}
	return ""
}
