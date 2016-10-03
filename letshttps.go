package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"rsc.io/letsencrypt"
)

var (
	backend   = flag.String("backend", "127.0.0.1:80", "backend server (127.0.0.1:80)")
	cache     = flag.String("cache", "letsencrypt.cache", "cache path (default: letsencrypt.cache)")
	httpaddr  = flag.String("http", "", "listen http addr (default: empty)")
	httpsaddr = flag.String("https", ":443", "listen https addr (:443)")
)

func NewReverseProxy(backend string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = req.Host

		ip, port, _ := net.SplitHostPort(req.RemoteAddr)
		req.Header.Set("X-Real-IP", ip)
		req.Header.Set("X-Remote-IP", ip)
		req.Header.Set("X-Remote-Port", port)
		req.Header.Set("X-Forwarded-For", ip)
		req.Header.Set("X-Forwarded-Proto", "https")
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Port", port)
		req.Header.Set("Host", req.Host)
	}

	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", backend)
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func main() {
	flag.Parse()

	var m letsencrypt.Manager
	if err := m.CacheFile(*cache); err != nil {
		log.Fatal(err)
	}

	handler := NewReverseProxy(*backend)

	go func() {
		if httpaddr == nil {
			return
		}

		s := &http.Server{
			Addr:    *httpaddr,
			Handler: http.HandlerFunc(letsencrypt.RedirectHTTP),
		}

		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	s := &http.Server{
		Addr:    *httpsaddr,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}

	if err := s.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}

}
