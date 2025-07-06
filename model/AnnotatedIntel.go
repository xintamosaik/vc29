package model
type AnnotatedIntel struct {
	CreatedAt   string                  `json:"created_at"` // This is a timestamp in string format, e.g., "1633072800"
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Content     [][]AnnotatedWord `json:"content"` // This is a slice of slices of AnnotatedWord, where each AnnotatedWord contains the word and its annotations
}
