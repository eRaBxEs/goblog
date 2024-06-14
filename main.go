package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type SlugReader interface {
	Reader(slug string) (string, error)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		fmt.Fprintf(w, "Post: %s", slug)
	})

	err := http.ListenAndServe(":7070", mux)
	if err != nil {
		log.Fatal(err)
	}

}

type FileReader struct{}

func (fsr FileReader) Reader(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
