package routes

import (
	"context"

	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"

	"vc29/internal/model"
	"vc29/internal/components"
)


func init() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
}

func Register() {

	// Endpoints for features
	http.Handle("GET /intel/new", templ.Handler(components.New()))
	http.HandleFunc("GET /intel/list", HandleIntelIndex)
	http.HandleFunc("POST /intel/create", HandleNewIntel)
	http.HandleFunc("GET /intel/annotate/{id}", HandleAnnotate)
	http.HandleFunc("POST /intel/annotate/{id}", HandleNewAnnotation)

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

	err = components.IntelList(intelShorts).Render(context.Background(), w)
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
	intelShorts, err := model.GetAllIntelShorts()
	if err != nil {
		http.Error(w, "Failed to r ead intel files", http.StatusInternalServerError)
		log.Println("Error reading intel files:", err)
		return
	}

	err = components.IntelList(intelShorts).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}
}

// handleAnnotate is a view that gets an intel data and then allows users to add annotations to it.
func HandleAnnotate(w http.ResponseWriter, r *http.Request) {

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

	annotatedIntel, err := model.GetAnnotatedIntel(intelID)
	if err != nil {
		http.Error(w, "Failed to get annotated intel", http.StatusInternalServerError)
		log.Println("Error getting annotated intel:", err)
		return
	}

	err = components.Annotate(ann, annotatedIntel).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}
}

func HandleNewAnnotation(w http.ResponseWriter, r *http.Request) {
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

	err = components.Annotate(ann, annotatedIntel).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}
}
