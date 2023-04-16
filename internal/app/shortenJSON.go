package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func (s *Server) shortenJSON(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqJSON := RequestJSON{}
	if err := json.Unmarshal(b, &reqJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	newURL, err := url.ParseRequestURI(reqJSON.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	newHash := GenerateHash(newURL.String(), s.Config.HashLen)
	if err := s.Store.Write(newHash, newURL.String()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respJSON := ResponseJSON{Result: fmt.Sprintf("%s/%s", s.Config.BaseUrl, newHash)}

	resp, err := json.Marshal(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))
}
