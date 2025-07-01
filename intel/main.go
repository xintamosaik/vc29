package intel

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const directory = "data/intel"

// IntelJSON represents the structure of the intel data stored in JSON files.
type IntelJSON struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Content     [][]string `json:"content"`
}

// IntelShort is a simplified version of IntelJSON, where the content is not included.
type IntelShort struct {
	CreatedAt   string
	Title       string
	Description string
}

type IntelFull struct {
	CreatedAt   string
	Title       string
	Description string
	Content     [][]string
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

	intelData := IntelJSON{
		Title:       title,
		Description: description,
		Content:     make([][]string, 0),
	}

	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		words := strings.Fields(strings.TrimSpace(line))
		intelData.Content = append(intelData.Content, words)
	}

	// Add the data/intel directory
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

	intelShorts, err := getAllIntelShorts()
	if err != nil {
		http.Error(w, "Failed to read intel file", http.StatusInternalServerError)
		log.Println("Error reading intel file:", err)
		return
	}

	err = Index(intelShorts).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}

}

// getIntelFull reads a JSON file from the data/intel directory,
// parses it into an IntelJSON struct, and returns an IntelFull struct
// containing the createdAt, title, description, and content fields.
//
// It returns an error if the file cannot be read or parsed.
func getIntelFull(fileName string) (IntelFull, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return IntelFull{}, err
	}
	defer file.Close()

	var intel IntelJSON
	var intelFull IntelFull
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&intel); err != nil {
		return IntelFull{}, err
	}

	trimmedFileName := strings.TrimSuffix(fileName, ".json")
	fileNameOnly := strings.TrimPrefix(trimmedFileName, directory+"/") // Whhich is a unix time stamp

	intelFull.CreatedAt = fileNameOnly
	intelFull.Description = intel.Description
	intelFull.Title = intel.Title
	intelFull.Content = intel.Content

	return intelFull, nil
}

// getIntelShort reads a JSON file from the data/intel directory,
// parses it into an IntelJSON struct, and returns an IntelShort struct
// containing the title, description, and createdAt fields.
//
// It returns an error if the file cannot be read or parsed.
func getIntelShort(fileName string) (IntelShort, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return IntelShort{}, err
	}
	defer file.Close()

	var intel IntelJSON
	var intelShort IntelShort
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&intel); err != nil {

		return IntelShort{}, err
	}

	trimmedFileName := strings.TrimSuffix(fileName, ".json")
	fileNameOnly := strings.TrimPrefix(trimmedFileName, directory+"/") // Whhich is a unix time stamp
	
	intelShort.CreatedAt =fileNameOnly 
	intelShort.Description = intel.Description
	intelShort.Title = intel.Title

	return intelShort, nil
}

// getAllIntelShorts reads all intel files from the data/intel directory,
// parses them into IntelShort objects, and returns a slice of these objects.
// If an error occurs during reading or parsing, it logs the error and continues with the next file.
//
// It returns a slice of IntelShort objects and an error if any.
func getAllIntelShorts() ([]IntelShort, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	intels := make([]IntelShort, 0, len(files))

	for _, file := range files {
		if !file.IsDir() {
			intel, err := getIntelShort(directory + "/" + file.Name())
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}
			intels = append(intels, intel)
		}
	}

	return intels, nil
}


// HandleIntelIndex handles the request for the Intel index page.
// It reads all intel files, creates a list of IntelShort objects,
// and renders the index template with the list.
//
// If an error occurs during reading or rendering, it responds with an error message.
func HandleIntelIndex(w http.ResponseWriter, r *http.Request) {

	log.Println("Handling Intel index page")

	intelShorts, err := getAllIntelShorts()
	if err != nil {
		http.Error(w, "Failed to read intel files", http.StatusInternalServerError)
		log.Println("Error reading intel files:", err)
		return
	}

	err = Index(intelShorts).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render intel page", http.StatusInternalServerError)
		log.Println("Error rendering intel page:", err)
		return
	}
}

// stampToDate converts a timestamp string to a formatted date string.
// It expects the timestamp to be in seconds since the Unix epoch.
// The returned date is formatted as "2006-01-02 15:04:05".
//
// Example: "1633072800" -> "2021-10-01 00:00:00"
//
// If the input is not a valid timestamp, it returns an error.
func stampToDate(fileNameOnly string) (string, error) {
	timestamp, err := strconv.ParseInt(fileNameOnly, 10, 64)
		if err != nil {
		return "", err
	}

	date := time.Unix(timestamp, 0)

	return date.Format("2006-01-02 15:04:05"), nil
}

// handleAnnotate is a view that gets an intel data and then allows users to add annotations to it.
func HandleAnnotate( w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Intel annotation")

	intelID := r.PathValue("id")
	
	if intelID == "" {
		http.Error(w, "Intel ID is required", http.StatusBadRequest)
		return
	}
	intelFull, err := getIntelFull(directory + "/" + intelID + ".json")
	if err != nil {
		http.Error(w, "Failed to read intel file", http.StatusInternalServerError)
		log.Println("Error reading intel file:", err)
		return
	}
	err = Annotate(intelFull).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}

	log.Println("Intel annotation page rendered successfully")

}