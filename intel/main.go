package intel

import (
	"encoding/json"
	"log"
	"time"
	"strconv"
	"net/http"
	"os"
	"strings"
	"context"

)


type IntelJSON struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Content     [][]string`json:"content"`
}

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

	intelData := IntelJSON{
		Title:       title,
		Description: description,
		Content: make([][]string, 0),
	}

	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		words := strings.Fields(strings.TrimSpace(line)) 
		intelData.Content = append(intelData.Content, words)
	}

	// Add the data/intel directory
	directory := "data/intel"
	if err := os.MkdirAll(directory, 0755); err != nil {
		log.Fatalf("Failed to create data/intel directory: %v", err)
	}

	// save as JSON file with timestamp converted to string
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	fileName := directory + "/" + timestamp + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(intelData); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
	log.Println("Intel data saved to", fileName)
	
	// Respond 
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = Index().Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}

}
