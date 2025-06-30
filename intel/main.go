package intel

import (
	"encoding/json"
	"log"
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
	log.Println("Content:", content)

	intelData := IntelJSON{
		Title:       title,
		Description: description,
		Content: make([][]string, 0),
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		words := strings.Fields(strings.TrimSpace(line)) 
		intelData.Content = append(intelData.Content, words)
	}

	// save as JSON file
	fileName := "intel.json"
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


	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = Index().Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}

}
