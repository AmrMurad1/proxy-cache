package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AmrMurad1/proxy-cache/internal"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 0, "Port for the caching proxy server")
	origin := flag.String("origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")

	flag.Parse()

	// Handling --clear-cache flag
	if *clearCache {
		fmt.Println("Clearing cache...")
		proxy := internal.NewProxy("")
		proxy.ClearCache()
		os.Exit(0)
	}

	// Validate required flags
	if *port == 0 {
		fmt.Println("Error: --port is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *origin == "" {
		fmt.Println("Error: --origin is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	proxy := internal.NewProxy(*origin)

	// Register proxy as HTTP handler
	http.Handle("/", proxy)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting caching proxy server on port %d\n", *port)
	fmt.Printf("Forwarding requests to %s\n", *origin)
	fmt.Printf("Server listening on %s\n", addr)

	// Start listening
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
