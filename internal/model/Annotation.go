package model

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const directoryAnnotations = "data/annotations"

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

func init() {
	if err := os.MkdirAll(directoryAnnotations, 0755); err != nil {
		log.Fatalf("Failed to create data/annotations directory: %v", err)
	}
}

func GetAllAnnotationIDs() ([]string, error) {
	
	files, err := os.ReadDir(directoryIntel)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.HasSuffix(fileName, ".json") {
			id := strings.TrimSuffix(fileName, ".json")
			ids = append(ids, id)
		}
		
	}
	return ids, nil
}

func LoadAllAnnotations(intelID string) ([]Annotation, error) {
	annotationsDir := directoryAnnotations + "/" + intelID

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

		annotations = append(annotations, annotation)
	}

	return annotations, nil
}
func SaveAnnotation(intelID string, annotation Annotation) error {
	if intelID == "" {
		log.Println("Intel ID is empty or invalid")
		return os.ErrInvalid
	}
	intelID = strings.ReplaceAll(intelID, "..", "")
	intelID = strings.TrimSpace(intelID)

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
