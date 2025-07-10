package main

import (

	"log"
	"net/http"
	"os"

	"vc29/internal/routes"
)

const port = ":3000"

func init() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
}

func main() {
	
	// just serve the static folder (above the fold)
	static := http.FileServer(http.Dir("static"))
	http.Handle("/", static)

	// And the dist folder (under the fold)
	dist := http.FileServer(http.Dir("dist"))
	http.Handle("/dist/", http.StripPrefix("/dist/", dist))

	routes.Register()
	// Start the HTTP server
	log.Println("Starting server on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
