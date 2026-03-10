// Takes an incoming request from a user
package proxy

import (
	"net/http"
	"log"
	"github.com/sukiecodes/sentinel-proxy/internal/lb"
)

// Proxy handles the incoming request and dispatches it to an alive/healthy backend
func Proxy(serverPool *lb.ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := serverPool.GetNextPeer() // ask balancer for next available backend

		if peer != nil {
			log.Printf("Proxying request from %s to %s", r.RemoteAddr, peer.URL)
			
			// use built-in reverse proxy to forward the request
			// handles copying the body, headers, and returning the response
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		// error handling, if no available backend (maybe none alive), error out
		log.Printf("No healthy backends available for request from %s", r.RemoteAddr)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}