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
	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}))

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

func PostHandler(sl SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		postMarkdown, err := sl.Reader(slug)
		if err != nil {
			// TODO: To handle different errors in the future
			http.Error(w, "Post not found!", http.StatusNotFound)
		}

		fmt.Fprint(w, postMarkdown)

	}
}
