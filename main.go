package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/xintamosaik/vc29/about"
	"github.com/xintamosaik/vc29/contact"
	"github.com/xintamosaik/vc29/drafts"
	"github.com/xintamosaik/vc29/home"
	"github.com/xintamosaik/vc29/intel"
	"github.com/xintamosaik/vc29/signals"
)

const port = ":3000"

func main() {

	// Housekeeping: Create a data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Bundle the JavaScript and CSS files using esbuild
	result := api.Build(api.BuildOptions{
		
		EntryPoints:       []string{"src.js"},
		Outfile:           "dist.js",
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		Format:            api.FormatIIFE,
		Platform:          api.PlatformBrowser,
		MinifyWhitespace:  false, // for dev builds - change later
		MinifyIdentifiers: false, // for dev builds - change later
		MinifySyntax:      false, // for dev builds - change later
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
		},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	// Static files: html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Static files: js bundle
	http.HandleFunc("GET /dist.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist.js")
	})

	// Static files: css bundle
	http.HandleFunc("GET /dist.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist.css")
	})

	// HTMX handlers: home
	http.Handle("/home", templ.Handler(home.Index()))

	// HTMX handlers: intel
	http.HandleFunc("/intel", intel.HandleIntelIndex)
	http.Handle("/intel/new", templ.Handler(intel.New()))
	http.HandleFunc("POST /intel/create", intel.HandleNewIntel)
	http.HandleFunc("GET /intel/annotate/{id}", intel.HandleAnnotate)
	
	// HTMX handlers: drafts
	http.Handle("/drafts", templ.Handler(drafts.Index()))

	// HTMX handlers: signals
	http.Handle("/signals", templ.Handler(signals.Index()))

	// HTMX handlers: about
	http.Handle("/about", templ.Handler(about.Index()))

	// HTMX handlers: contact
	http.Handle("/contact", templ.Handler(contact.Index()))

	// Start the HTTP server
	fmt.Println("Starting server on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
