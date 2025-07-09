package model

import (
	"log"
	"strconv"
)

type AnnotatedIntel struct {
	CreatedAt   string                  `json:"created_at"` // This is a timestamp in string format, e.g., "1633072800"
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Content     [][]AnnotatedWord `json:"content"` // This is a slice of slices of AnnotatedWord, where each AnnotatedWord contains the word and its annotations
}

func GetAnnotatedIntel(id string) (AnnotatedIntel, error) {

	if id == "" {
		log.Println("No Intel ID provided")
		return AnnotatedIntel{}, nil
	}

	intel, err := LoadIntel(id)
	if err != nil {
		return AnnotatedIntel{}, err
	}

	var intelFull IntelFull

	intelFull.CreatedAt = id
	intelFull.Description = intel.Description
	intelFull.Title = intel.Title
	intelFull.Content = intel.Content

	ann, err := LoadAllAnnotations(id)
	if err != nil {
		log.Println("Error getting annotations for Intel ID:", id, err)
		return AnnotatedIntel{}, err
	}

	annotatedContent := make([][]AnnotatedWord, len(intelFull.Content))
	for i, paragraph := range intelFull.Content {
		annotatedContent[i] = make([]AnnotatedWord, len(paragraph))
		for j, word := range paragraph {
			annotatedContent[i][j] = AnnotatedWord{
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

	return AnnotatedIntel{
		CreatedAt:   intelFull.CreatedAt,
		Title:       intelFull.Title,
		Description: intelFull.Description,
		Content:     annotatedContent,
	}, nil

}
