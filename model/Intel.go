package model

import (
	"encoding/json"
	"log"

	"os"
	"strconv"
	"strings"
	"time"
)


const directoryIntel = "data/intel"

func init() {
	if err := os.MkdirAll(directoryIntel, 0755); err != nil {
		log.Fatalf("Failed to create data/intel directory: %v", err)
	}
}

type Intel struct {
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

// IntelFull is a more detailed version of IntelJSON, including the content.
type IntelFull struct {
	CreatedAt   string
	Title       string
	Description string
	Content     [][]string
}

func SaveIntel(title, description, content string) error {
	intelData := Intel{
		Title:       title,
		Description: description,
		Content:     make([][]string, 0),
	}

	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		words := strings.Fields(strings.TrimSpace(line))
		intelData.Content = append(intelData.Content, words)
	}

	// save as JSON file with timestamp converted to string
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	fileName := directoryIntel + "/" + timestamp + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(intelData); err != nil {

		log.Println("Error encoding JSON:", err)

	}
	log.Println("Intel data saved to", fileName)

	return nil
}

func LoadIntel(id string) (Intel, error) {
	fileName := directoryIntel + "/" + id + ".json"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return Intel{}, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return Intel{}, err
	}
	defer file.Close()

	var intel Intel

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&intel); err != nil {
		return Intel{}, err
	}

	return intel, nil
}


// getIntelShort reads a JSON file from the data/intel directory,
// parses it into an IntelJSON struct, and returns an IntelShort struct
// containing the title, description, and createdAt fields.
//
// It returns an error if the file cannot be read or parsed.
func GetIntelShort(id string) (IntelShort, error) {

	intel, err := LoadIntel(id)
	if err != nil {
		return IntelShort{}, err
	}
	
	var intelShort IntelShort

	intelShort.CreatedAt = id
	intelShort.Description = intel.Description
	intelShort.Title = intel.Title

	return intelShort, nil
}


// Reads all intel files from the data/intel directory,
// parses them into IntelShort objects, and returns a slice of these objects.
// If an error occurs during reading or parsing, it logs the error and continues with the next file.
//
// It returns a slice of IntelShort objects and an error if any.
func GetAllIntelShorts() ([]IntelShort, error) {
	intelIDs, err := GetAllAnnotationIDs()
	if err != nil {
		log.Println("Error getting intel IDs:", err)
		return nil, err
	}

	intels := make([]IntelShort, 0, len(intelIDs))
	for _, id := range intelIDs {

		intel, err := GetIntelShort(id)
		if err != nil {
			log.Printf("Error reading file %s: %v", id, err)
			continue
		}
		intels = append(intels, intel)

	}

	return intels, nil
}
