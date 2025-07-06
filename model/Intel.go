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

func LoadIntel(fileName string) (Intel, error) {
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
