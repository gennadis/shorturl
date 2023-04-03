package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gennadis/shorturl/internal/app/hash"
)

const HashLen = 6

var urls = make(map[string]string)

// Accepts POST requests with `url` to shorten in request body.
// Returns 201 and short url in request body if successful.
func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", 400)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	id := hash.GenerateHash(HashLen)
	urls[id] = string(b)

	response := fmt.Sprintf("http://localhost:8080/%s", id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/", shortenURL)
	http.ListenAndServe(":8080", nil)
}
