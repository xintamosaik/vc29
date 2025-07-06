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
}
func main() {

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

	http.Handle("GET /home", templ.Handler(pages.Home()))
	http.Handle("GET /intel", templ.Handler(pages.Intel()))

	http.HandleFunc("GET /intel/list", HandleIntelIndex)

	http.Handle("GET /intel/new", templ.Handler(New()))

	http.HandleFunc("POST /intel/create", HandleNewIntel)
	http.HandleFunc("GET /intel/annotate/{id}", HandleAnnotate)
	http.HandleFunc("POST /intel/annotate/{id}", HandleNewAnnotation)
	http.Handle("GET /drafts", templ.Handler(pages.Drafts()))
	http.Handle("GET /signals", templ.Handler(pages.Signals()))

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
	annotatedIntel, err := GetAnnotatedIntel(intelID)
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

	startParagraph := r.FormValue("start_paragraph")
	startWord := r.FormValue("start_word")
	endParagraph := r.FormValue("end_paragraph")
	endWord := r.FormValue("end_word")
	keyword := r.FormValue("keyword")
	description := r.FormValue("description")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	annotation := model.Annotation{
		StartParagraph: startParagraph,
		StartWord:      startWord,
		EndParagraph:   endParagraph,
		EndWord:        endWord,
		Keyword:        keyword,
		Description:    description,
		UpdatedAt:      timestamp,
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

	annotatedIntel, err := GetAnnotatedIntel(intelID)
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


func GetAnnotatedIntel(id string) (model.AnnotatedIntel, error) {

	if id == "" {
		log.Println("No Intel ID provided")
		return model.AnnotatedIntel{}, nil
	}

	intel, err := model.LoadIntel(id)
	if err != nil {
		return model.AnnotatedIntel{}, err
	}

	var intelFull model.IntelFull

	intelFull.CreatedAt = id
	intelFull.Description = intel.Description
	intelFull.Title = intel.Title
	intelFull.Content = intel.Content

	ann, err := model.LoadAllAnnotations(id)
	if err != nil {
		log.Println("Error getting annotations for Intel ID:", id, err)
		return model.AnnotatedIntel{}, err
	}

	annotatedContent := make([][]model.AnnotatedWord, len(intelFull.Content))
	for i, paragraph := range intelFull.Content {
		annotatedContent[i] = make([]model.AnnotatedWord, len(paragraph))
		for j, word := range paragraph {
			annotatedContent[i][j] = model.AnnotatedWord{
				Word:          word,
				AnnotationIDs: []string{}, // Initialize with an empty slice
			}
			// Check if there are annotations for this sequence of words
			for _, annotation := range ann {
				// Convert string indices to integers for proper comparison
				startParagraph, err := strconv.Atoi(annotation.StartParagraph)
				if err != nil {
					continue
				}
				endParagraph, err := strconv.Atoi(annotation.EndParagraph)
				if err != nil {
					continue
				}
				startWord, err := strconv.Atoi(annotation.StartWord)
				if err != nil {
					continue
				}
				endWord, err := strconv.Atoi(annotation.EndWord)
				if err != nil {
					continue
				}

				// Check if current position is within annotation range
				isWithinAnnotation := false

				if i < startParagraph {
					continue
				}
				if i > endParagraph {
					continue
				}

				if i == startParagraph && i == endParagraph {
					// Annotation is within the same paragraph
					if j >= startWord && j <= endWord {
						isWithinAnnotation = true
					}
				} else if i == startParagraph {
					// Current paragraph is the start paragraph
					if j >= startWord {
						isWithinAnnotation = true
					}
				} else if i == endParagraph {
					// Current paragraph is the end paragraph
					if j <= endWord {
						isWithinAnnotation = true
					}
				} else {
					// Current paragraph is between start and end paragraphs
					isWithinAnnotation = true
				}
				if !isWithinAnnotation {
					continue
				}

				// If the annotation is within the range, add the ID and keyword
				annotatedContent[i][j].AnnotationIDs = append(annotatedContent[i][j].AnnotationIDs, annotation.UpdatedAt)
				annotatedContent[i][j].Keywords = append(annotatedContent[i][j].Keywords, annotation.Keyword)

			}
		}
	}

	return model.AnnotatedIntel{
		CreatedAt:   intelFull.CreatedAt,
		Title:       intelFull.Title,
		Description: intelFull.Description,
		Content:     annotatedContent,
	}, nil

}
