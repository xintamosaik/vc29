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

	"github.com/a-h/templ"
)

const directory = "data/intel"
const directoryAnnotations = "data/annotations"

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

	intelShort.CreatedAt = fileNameOnly
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

func getAnnotations(intelID string) ([]Annotation, error) {
	annotationsDir := directoryAnnotations + "/" + intelID

	files, err := os.ReadDir(annotationsDir)
	if err != nil {
		log.Printf("Error reading annotations directory %s: %v", annotationsDir, err)
		return nil, err
	}

	
	annotations := make([]Annotation, 0, len(files))

	for _, file := range files {
		
		if !file.IsDir() {
			filePath := annotationsDir + "/" + file.Name()
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading annotation file %s: %v", file.Name(), err)
				continue
			}

			var annotation Annotation
			if err := json.Unmarshal(fileContent, &annotation); err != nil {
				log.Printf("Error unmarshalling annotation file %s: %v", file.Name(), err)
				continue
			}

			annotation.UpdatedAt = strings.TrimSuffix(file.Name(), ".json") // Use the file name as the updated_at field
			annotations = append(annotations, annotation)
		}
	}

	return annotations, nil
}

// handleAnnotate is a view that gets an intel data and then allows users to add annotations to it.
func HandleAnnotate(w http.ResponseWriter, r *http.Request) {
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

	annotations,err := getAnnotations(intelID)
	if err != nil {
		http.Error(w, "Failed to read annotations", http.StatusInternalServerError)
		log.Println("Error reading annotations:", err)
		return
	}
	log.Println("Annotations loaded successfully for Intel ID:", intelID)
	log.Println("Number of annotations:", len(annotations))
	err = Annotate(intelFull, annotations).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}

	log.Println("Intel annotation page rendered successfully")
}

type Annotation struct {
	StartParagraph string `json:"start_paragraph"`
	StartWord      string `json:"start_word"`
	EndParagraph   string `json:"end_paragraph"`
	EndWord        string `json:"end_word"`
	Keyword        string `json:"keyword"`
	Description    string `json:"description"`
	UpdatedAt      string `json:"updated_at"` // We do not need created_at, because we use the file name as a timestamp
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

	log.Printf("New annotation for Intel ID %s: start(%s, %s), end(%s, %s), description: %s\n",
		intelID, startParagraph, startWord, endParagraph, endWord, description) // And yes, that worked

	// Create the annotations directory if it doesn't exist
	if err := os.MkdirAll(directoryAnnotations+"/"+intelID, 0755); err != nil {
		http.Error(w, "Failed to create annotations directory", http.StatusInternalServerError)
		log.Println("Error creating annotations directory:", err)
		return
	}

	// We will use unix timestamps again for the annotation file name
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	annotationFileName := directoryAnnotations + "/" + intelID + "/" + timestamp + ".json"
	annotationFile, err := os.Create(annotationFileName)
	if err != nil {
		http.Error(w, "Failed to create annotation file", http.StatusInternalServerError)
		log.Println("Error creating annotation file:", err)
		return
	}
	defer annotationFile.Close()

	annotation := Annotation{
		StartParagraph: startParagraph,
		StartWord:      startWord,
		EndParagraph:   endParagraph,
		EndWord:        endWord,
		Keyword:        keyword,
		Description:    description,
		UpdatedAt:      timestamp, // We use the timestamp as the updated_at field
	}
	encoder := json.NewEncoder(annotationFile)
	if err := encoder.Encode(annotation); err != nil {
		http.Error(w, "Failed to encode annotation JSON", http.StatusInternalServerError)
		log.Println("Error encoding annotation JSON:", err)
		return
	}

	// We need to actually save the annotation to a file
	log.Printf("Annotation for Intel ID %s saved to %s\n", intelID, annotationFileName)
	// Log the success
	log.Printf("Annotation saved to %s\n", annotationFileName)
	// Respond with success
	w.Header().Set("Content-Type", "text/plain")

	// We show the use the annotation page for the intel again, so they can add more annotations
	w.WriteHeader(http.StatusOK)

	log.Println("Annotation page rendered successfully")
	log.Println("New annotation submission handled successfully")
	// Most web devs do success messages but we will do something clever.
	// So we will actually load the annotations that already exist. And if one of them is only some minutes old, we will show something that
	// gives a hint about that. "Latest add: some minutes ago - keyword: something"
	// Later though. Which I guess is now..


	annotations := make([]Annotation, 0)
	log.Println("Loading annotations for Intel ID:", intelID)
	log.Println("Annotations directory:", directoryAnnotations+"/"+intelID)
	log.Println("Annotations directory exists:", directoryAnnotations+"/"+intelID)

	annotations, err = getAnnotations(intelID)
	if err != nil {
		http.Error(w, "Failed to read annotations", http.StatusInternalServerError)
		log.Println("Error reading annotations:", err)
		return
	}

	log.Println("Annotations loaded successfully for Intel ID:", intelID)
	log.Println("Number of annotations:", len(annotations))

	intelFull, err := getIntelFull(directory + "/" + intelID + ".json")
	if err != nil {
		http.Error(w, "Failed to read intel file", http.StatusInternalServerError)
		log.Println("Error reading intel file:", err)
		return
	}

	err = Annotate(intelFull, annotations).Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render annotation page", http.StatusInternalServerError)
		log.Println("Error rendering annotation page:", err)
		return
	}
}

// Register initializes the Intel handlers for the HTMX routes.
func Register() {
	http.HandleFunc("GET /intel", HandleIntelIndex)
	http.Handle("GET /intel/new", templ.Handler(New()))
	http.HandleFunc("POST /intel/create", HandleNewIntel)
	http.HandleFunc("GET /intel/annotate/{id}", HandleAnnotate)
	http.HandleFunc("POST /intel/annotate/{id}", HandleNewAnnotation)
}
