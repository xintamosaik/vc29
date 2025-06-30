package main

import (
	"fmt"

	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/xintamosaik/vc29/home"
		"github.com/xintamosaik/vc29/about"
			"github.com/xintamosaik/vc29/contact"
)

const port = ":3000"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// HTMX handlers:
	// HTMX handler for GET /home
	http.Handle("/home", templ.Handler(home.Index()))
	// HTMX handler for GET /about
	http.Handle("/about", templ.Handler(about.Index()))
	
	// HTMX handler for GET /contact
	http.Handle("/contact", templ.Handler(contact.Index()))

	fmt.Println("Starting server on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
