// Simple HTTP server that listens on a port provided through command line arguments and returns a message for a request
// Will be used as a dummy backend service to test our sentinel proxy

package backend

import (
	"fmt" // general I/O formatting
	"log" // application logging
	"net/http"
	"os"
)

func main() {
	// passing the port as a command line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <port>")
	}
	port := os.Args[1]

	// registers a handler for all requests to "/"
	// w http.ResponseWriter used to send the HTTP response, r *http.Request contains request information
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from backend running on port %s\n", port)
		log.Printf("Received request from %s", r.RemoteAddr)
	}) 

	log.Printf("Starting dummy backend on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}