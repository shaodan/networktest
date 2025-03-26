package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Define command line flags
	serverMode := flag.Bool("server", false, "Run in server mode")
	clientMode := flag.Bool("client", false, "Run in client mode")
	host := flag.String("host", "localhost", "Host to connect to in client mode or bind to in server mode")
	port := flag.Int("port", 8080, "Port to use")
	count := flag.Int("count", 10, "Number of packets to send in client mode")
	interval := flag.Int("interval", 1000, "Interval between packets in milliseconds")

	// Parse command line flags
	flag.Parse()

	// Validate flags
	if *serverMode == *clientMode {
		fmt.Println("Error: Must specify either -server or -client mode")
		flag.Usage()
		os.Exit(1)
	}

	address := fmt.Sprintf("%s:%d", *host, *port)

	// Run in server or client mode
	if *serverMode {
		fmt.Printf("Starting server on %s\n", address)
		if err := runServer(address); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		fmt.Printf("Starting client, connecting to %s\n", address)
		if err := runClient(address, *count, *interval); err != nil {
			log.Fatalf("Client error: %v", err)
		}
	}
}