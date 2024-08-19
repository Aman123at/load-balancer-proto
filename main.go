package main

import (
	"io"
	"log"
	"net/http"
)

var Servers = []string{"http://localhost:8001", "http://localhost:8002", "http://localhost:8003", "http://localhost:8004", "http://localhost:8005"}
var requests int = 0

func getServerIdx() int {
	log.Println(requests)
	return requests % len(Servers)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	requests = requests + 1
	serverIdx := getServerIdx()
	server := Servers[serverIdx]

	// Create a new request to forward
	proxyReq, err := http.NewRequest(r.Method, server, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Copy headers
	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	// Forward the request
	proxyClient := &http.Client{}
	proxyResp, err := proxyClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	// Copy response headers
	for header, values := range proxyResp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// Copy response body
	io.Copy(w, proxyResp.Body)

}

func main() {
	log.Println("Load balancer")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
