package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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
			http.Error(w, "Error reading from file", http.StatusInternalServerError)
			return
		}
		mdRenderer := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
			),
		)
		var buf bytes.Buffer
		err = mdRenderer.Convert([]byte(postMarkdown), &buf)
		if err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
			return
		}
		if err != nil {
			// TODO: To handle different errors in the future
			http.Error(w, "Post not found!", http.StatusNotFound)
			return
		}

		io.Copy(w, &buf)

	}
}
