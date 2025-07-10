package main

import (
	"context"
	"log"
	"os"
	"path"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"vc29/internal/components"
	"vc29/internal/layouts"
	"vc29/internal/pages"

	"github.com/a-h/templ"
)

type mainPage struct {
	slug string
	filename       string
	templComponent templ.Component
}

var mainPages = []mainPage{
	{"home", "index.html", pages.Home()},
	{"intel", "intel.html", pages.Intel()},
	{"signals", "signals.html", pages.Signals()},
	{"drafts", "drafts.html", pages.Drafts()},
}

func main() {
	staticPath := "static"
	
	for _, page := range mainPages {
		fileHome := path.Join(staticPath, page.filename)
		file, err := os.Create(fileHome)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		file.WriteString("<!DOCTYPE html>")
		// Write it out.
		navigation := components.Navigation(page.slug)
		body := layouts.Body( navigation, page.templComponent)
		err = layouts.Frame("VC29 | " + page.slug, body).Render(context.Background(), file)
		if err != nil {
			log.Fatalf("failed to write index page: %v", err)
		}
	}


	// Bundle the JavaScript and CSS files using esbuild
	result := esbuild.Build(esbuild.BuildOptions{

		EntryPoints:       []string{"src.js"},
		Outfile:           "dist/under_the_fold.js",
		Bundle:            true,
		Write:             true,
		LogLevel:          esbuild.LogLevelInfo,
		Format:            esbuild.FormatIIFE,
		Platform:          esbuild.PlatformBrowser,
		MinifyWhitespace:  false, // for dev builds - change later
		MinifyIdentifiers: false, // for dev builds - change later
		MinifySyntax:      false, // for dev builds - change later
		Loader: map[string]esbuild.Loader{
			".css": esbuild.LoaderCSS,
		},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}
