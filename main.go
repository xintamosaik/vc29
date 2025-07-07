package main

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/a-h/templ"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/xintamosaik/vc29/model"
	"github.com/xintamosaik/vc29/pages"
)

const port = ":3000"

func init() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Bundle the JavaScript and CSS files using esbuild
	result := api.Build(api.BuildOptions{

		EntryPoints:       []string{"src.js"},
		Outfile:           "dist/under_the_fold.js",
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
}

func main() {

	// just serve the static folder (above the fold)
	static:= http.FileServer(http.Dir("static"))
	http.Handle("/", static)

	// And the dist folder (under the fold)
	dist := http.FileServer(http.Dir("dist"))
	http.Handle("/dist/", http.StripPrefix("/dist/", dist))
	http.Handle("GET /under_the_fold", templ.Handler(under_the_fold()))

	// Endpoints for features
	http.Handle("GET /intel/new", templ.Handler(pages.New()))
	http.HandleFunc("GET /intel/list", HandleIntelIndex)
	http.HandleFunc("POST /intel/create", HandleNewIntel)
	http.HandleFunc("GET /intel/annotate/{id}", HandleAnnotate)
	http.HandleFunc("POST /intel/annotate/{id}", HandleNewAnnotation)

	// Start the HTTP server
	fmt.Println("Starting server on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// This function handles the submission of new intel data.
// It processes the form data, creates a new IntelJSON object,
// saves it as a JSON file in the data/intel directory, and then renders the Intel index page.
// If the request method is not POST, it responds with a "Method not allowed" error.
//
// It also handles errors related to file creation, encoding, and rendering the page.
// If any error occurs, it logs the error and responds with an appropriate HTTP status code.
//
// The function expects the form data to contain "title", "description", and "content"
// fields, where "content" is a multiline string that will be split into words and stored
// as a slice of strings in the IntelJSON object.
func HandleNewIntel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("New intel submitted")

	// Process form data
	title := r.FormValue("title")
	description := r.FormValue("description")
	content := r.FormValue("content")

	err := model.SaveIntel(title, description, content)
	if err != nil {
		http.Error(w, "Failed to save intel data", http.StatusInternalServerError)
		log.Println("Error saving intel data:", err)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	intelShorts, err := model.GetAllIntelShorts()
	if err != nil {
		http.Error(w, "Failed to read intel file", http.StatusInternalServerError)
		log.Println("Error reading intel file:", err)
		return
	}

	err = pages.IntelList(intelShorts).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}

}

// HandleIntelIndex handles the request for the Intel index page.
// It reads all intel files, creates a list of IntelShort objects,
// and renders the index template with the list.
//
// If an error occurs during reading or rendering, it responds with an error message.
func HandleIntelIndex(w http.ResponseWriter, r *http.Request) {

	log.Println("Handling Intel index page")

	intelShorts, err := model.GetAllIntelShorts()
	if err != nil {
		http.Error(w, "Failed to r ead intel files", http.StatusInternalServerError)
		log.Println("Error reading intel files:", err)
		return
	}

	err = pages.IntelList(intelShorts).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}
}

// handleAnnotate is a view that gets an intel data and then allows users to add annotations to it.
func HandleAnnotate(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Intel annotation")

	intelID := r.PathValue("id")

	if intelID == "" {
		http.Error(w, "Intel ID is required", http.StatusBadRequest)
		return
	}

	ann, err := model.LoadAllAnnotations(intelID)
	if err != nil {
		http.Error(w, "Failed to read annotations", http.StatusInternalServerError)
		log.Println("Error reading annotations:", err)
		return
	}
	log.Println("Annotations loaded successfully for Intel ID:", intelID)
	log.Println("Number of annotations:", len(ann))
	annotatedIntel, err := model.GetAnnotatedIntel(intelID)
	if err != nil {
		http.Error(w, "Failed to get annotated intel", http.StatusInternalServerError)
		log.Println("Error getting annotated intel:", err)
		return
	}
	err = pages.Annotate(ann, annotatedIntel).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}

	log.Println("Intel annotation page rendered successfully")
}

func HandleNewAnnotation(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling new annotation submission")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	intelID := r.PathValue("id")
	if intelID == "" {
		http.Error(w, "Intel ID is required", http.StatusBadRequest)
		return
	}

	annotation := model.Annotation{
		StartParagraph: r.FormValue("start_paragraph"),
		StartWord:      r.FormValue("start_word"),
		EndParagraph:   r.FormValue("end_paragraph"),
		EndWord:        r.FormValue("end_word"),
		Keyword:        r.FormValue("keyword"),
		Description:    r.FormValue("description"),
		UpdatedAt:      strconv.FormatInt(time.Now().Unix(), 10),
	}

	err := model.SaveAnnotation(intelID, annotation)
	if err != nil {
		http.Error(w, "Failed to create annotation", http.StatusInternalServerError)
		log.Println("Error creating annotation:", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	ann := make([]model.Annotation, 0)

	ann, err = model.LoadAllAnnotations(intelID)
	if err != nil {
		http.Error(w, "Failed to read annotations", http.StatusInternalServerError)
		log.Println("Error reading annotations:", err)
		return
	}

	annotatedIntel, err := model.GetAnnotatedIntel(intelID)
	if err != nil {
		http.Error(w, "Failed to get annotated intel", http.StatusInternalServerError)
		log.Println("Error getting annotated intel:", err)
		return
	}

	err = pages.Annotate(ann, annotatedIntel).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}
}
