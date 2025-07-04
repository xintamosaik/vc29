package annotations

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// Annotation represents an annotation on an intel data.
type Annotation struct {
	StartParagraph string `json:"start_paragraph"`
	StartWord      string `json:"start_word"`
	EndParagraph   string `json:"end_paragraph"`
	EndWord        string `json:"end_word"`
	Keyword        string `json:"keyword"`
	Description    string `json:"description"`
	UpdatedAt      string `json:"updated_at"` // We do not need created_at, because we use the file name as a timestamp
}

type AnnotatedWord struct {
	Word          string   `json:"word"`
	AnnotationIDs []string `json:"annotation_id"`     // These are the IDs of the annotations that apply to this word
	Keywords      []string `json:"keyword,omitempty"` // Optional keyword for the word, if applicable
}

type AnnotatedIntel struct {
	CreatedAt   string            `json:"created_at"` // This is a timestamp in string format, e.g., "1633072800"
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Content     [][]AnnotatedWord `json:"content"` // This is a slice of slices of AnnotatedWord, where each AnnotatedWord contains the word and its annotations
}

const directoryAnnotations = "data/annotations"

func Save(intelID string, annotation Annotation) error {

	// Create the annotations directory if it doesn't exist
	if err := os.MkdirAll(directoryAnnotations+"/"+intelID, 0755); err != nil {

		log.Println("Error creating annotations directory:", err)
		return err
	}

	// We will use unix timestamps again for the annotation file name

	annotationFileName := directoryAnnotations + "/" + intelID + "/" + annotation.UpdatedAt + ".json"
	annotationFile, err := os.Create(annotationFileName)
	if err != nil {

		log.Println("Error creating annotation file:", err)
		return err
	}
	defer annotationFile.Close()

	encoder := json.NewEncoder(annotationFile)
	if err := encoder.Encode(annotation); err != nil {

		log.Println("Error encoding annotation JSON:", err)
		return err
	}

	return nil
}

func GetAll(intelID string) ([]Annotation, error) {
	annotationsDir := directoryAnnotations + "/" + intelID

	if err := os.MkdirAll(annotationsDir, 0755); err != nil {
		log.Printf("Error creating annotations directory %s: %v", annotationsDir, err)
		return nil, err
	}

	files, err := os.ReadDir(annotationsDir)
	if err != nil { // This is very improbable, but we handle it anyway
		log.Printf("Error reading annotations directory %s: %v", annotationsDir, err)
		return nil, err
	}

	annotations := make([]Annotation, 0, len(files))

	for _, file := range files {

		if file.IsDir() {
			continue
		}

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

	return annotations, nil
}
