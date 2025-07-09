package main
import (
	"context"
	"log"
	"path"
	"os"

	"vc29/internal/pages"
)
func main() {
	staticPath := "static"
	
	name := path.Join(staticPath, "home.html")
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	// Write it out.
	err = pages.Home().Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}

}
