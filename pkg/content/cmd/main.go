package main

import (
	"log"
	"os"

	"events-consumers.psmarcin.dev/content"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	funcframework.RegisterEventFunction("/", content.Get)
	//funcframework.RegisterEventFunction("/", content.Process)
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
