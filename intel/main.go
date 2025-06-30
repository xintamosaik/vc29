package intel

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
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

	// Form data
	title := r.FormValue("title")
	log.Println("Title:", title)
	description := r.FormValue("description")
	log.Println("Description:", description)

	content := r.FormValue("content")
	log.Println("Content:", content)
	// Create IntelJSON object
	intelData := IntelJSON{
		Title:       title,
		Description: description,
		Content: make([][]string, 0),
	}
	// We will split the content into lines and the lines into words.
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		
		words := strings.Fields(strings.TrimSpace(line))  // Split line into words
		
		intelData.Content = append(intelData.Content, words)

	}
	log.Println("Content lines:", lines)

	log.Println("Intel data:", intelData)

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

}
