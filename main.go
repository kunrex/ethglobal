package main

import (
	"log"
	"net/http"
	"os"

	"git-server/pkg/routes"
	"git-server/pkg/server"
)

func main() {
	// Create the Git server
	gitServer := server.NewInMemoryGitServer()

	// Create some sample repositories
	err := gitServer.CreateRepository("test-repo")
	if err != nil {
		log.Printf("Warning: Failed to create test repository: %v", err)
	}

	err = gitServer.CreateRepository("demo")
	if err != nil {
		log.Printf("Warning: Failed to create demo repository: %v", err)
	}

	// Set up routes
	router := routes.SetupRoutes(gitServer)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Git server on port %s", port)
	log.Printf("Visit http://localhost:%s to see the web interface", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}