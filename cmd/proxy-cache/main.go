package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO: clear cache
	port := flag.Int("port", 0, "Port for the caching proxy server")
	origin := flag.String("origin", "", "Origin server URL")

	flag.Parse()

	if *port == 0 || *origin == "" {
		log.Fatal("Error: --port and --origin are required")
	}

	fmt.Printf("Starting proxy server on port %d\n", *port)
	fmt.Printf("Forwarding to %s\n", *origin)

	// Create http handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello from proxy!")
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
