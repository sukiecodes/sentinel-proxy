// Defines a flag, initializes pool of backends, launch goroutines, starts main server and pass all requests
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"fmt"

	"github.com/sukiecodes/sentinel-proxy/internal/lb"
	"github.com/sukiecodes/sentinel-proxy/internal/proxy"
)

func main() {
	// define command line flags
	var serverList string 
	var port int 
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, separated by comma")
	flag.IntVar(&port, "port", 8080, "Port to serve proxy on")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide at least one backend using --backends flag")
	}

	// parse the comma-separated string into a list of URLs
	servers := strings.Split(serverList, ",")
	var backends []*lb.Backend 

	for _, s := range servers {
		serverUrl, err := url.Parse(s)
		if err != nil {
			log.Fatalf("Invalid backend URL: %s", s)
		}

		// create the ReverseProxy for this specific backend, will copy request to backend 
		proxyHandler := httputil.NewSingleHostReverseProxy(serverUrl)

		backends = append(backends, &lb.Backend{
			URL: serverUrl,
			Alive: true, // since it was just created, assume alive until first health check
			ReverseProxy: proxyHandler,
		})

		log.Printf("Configured backend: %s\n", serverUrl)
	}

	// initialize ServerPool 
	pool := &lb.ServerPool{
		Backends: backends,
	}

	// start health check in the background (essentially heartbeats)
	go pool.HealthCheck()

	// define the server and the handler 
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: proxy.Proxy(pool),
	}

	log.Printf("Sentinel Proxy started at: %d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}