package main

import (
	"fmt"
	"gpu-bpf-operator/probe-agent/internal/router"
	"log"
	"os"
)

const version = "v0.0.1"

func main() {
	port := os.Getenv("PROBE_AGENT_PORT")
	if port == "" {
		port = "8080"
	}
	r := router.Setup(version)

	addr := fmt.Sprintf("localhost:%s", port)
	log.Printf("Starting probe-agent %s on %s", version, addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
