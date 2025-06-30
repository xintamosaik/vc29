package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)
const port = ":3000"

type PageData struct {
	Title string
	Content string
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Hello World",
			Content: "Welcome to the world of web development. Enjoy coding!",
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("Starting server on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
