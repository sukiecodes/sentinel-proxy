// Will verify that a backend is capable of handling user's traffic before proxy tries sending user there
package lb

import (
	"net/http"
	"fmt"
	"time"
	"log"
)

// Ping function takes in a backend server and returns true if able to connect
func Ping(b *Backend) bool {
	// defining a timeout so no deadlock
	timeout := 2 * time.Second
	client := http.Client {
		Timeout: timeout,
	}

	resp, err := client.Get(b.URL.String())

	// checking if there's a network error
	if err != nil {
		fmt.Println("Error trying to connect: ", err)
		return false
	} 

	// clean up response body immediately
	// placing this here so that we know resp exists and that there are no errors
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ping failed for %s: Status %d", b.URL, resp.StatusCode)
		return false
	}

	return true 
}

func (s *ServerPool) HealthCheck() {
	// ticker that fires every 10 seconds
	ticker := time.NewTicker(10 * time.Second)

	// looping forever, waiting for the tick
	for range ticker.C {
		log.Println("Starting active health check...")
		s.checkBackends()
	}
}

func (s *ServerPool) checkBackends() {
	for _, b := range s.Backends {
		status := Ping(b)
		b.SetAlive(status)
		log.Printf("Backend %s is alive: %v", b.URL, status)
	}
}