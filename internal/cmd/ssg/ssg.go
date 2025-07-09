package main

import (
	"context"
	"log"
	"os"
	"path"

	"vc29/internal/layouts"
	"vc29/internal/pages"

	"github.com/a-h/templ"
)

type mainPage struct {
	filename       string
	templComponent templ.Component
}

var mainPages = []mainPage{
	{"home.html", pages.Home()},
	{"intel.html", pages.Intel()},
	{"signals.html", pages.Signals()},
	{"drafts.html", pages.Drafts()},
}

func main() {
	staticPath := "static"
	
	fileHome := path.Join(staticPath, "home.html")
	f, err := os.Create(fileHome)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	// Write it out.
	home := pages.Home()
	err = layouts.Frame("VC29 | home", home).Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}

}
