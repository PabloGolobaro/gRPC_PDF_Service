package main

import (
	"controller-service-grpc/internal/routes"
	"syscall"
)

const DefaultPort = "8080"

func main() {
	port := DefaultPort
	if value, ok := syscall.Getenv("CS_PORT"); ok {
		port = value
	}

	routes.Start(port)
}
