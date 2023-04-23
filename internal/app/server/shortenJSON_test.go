package server

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestShortenJSON(t *testing.T) {
	validURLReqJSON, _ := json.Marshal(RequestJSON{URL: longURL})
	validURLRespJSON, _ := json.Marshal(ResponseJSON{Result: shortURL})

	invalidURLReqJSON, _ := json.Marshal(RequestJSON{URL: "qwertyuiop"})

	emptyURLReqJSON, _ := json.Marshal(RequestJSON{URL: ""})

	tests := []test{
		{
			name:    "Shorten valid original URL",
			request: "/api/shorten",
			method:  http.MethodPost,
			body:    string(validURLReqJSON),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				body:        string(validURLRespJSON),
			},
		},
		{
			name:    "Shorten invalid original URL",
			request: "/api/shorten",
			method:  http.MethodPost,
			body:    string(invalidURLReqJSON),
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid URL\n",
			},
		},
		{
			name:    "Shorten empty original URL",
			request: "/api/shorten",
			method:  http.MethodPost,
			body:    string(emptyURLReqJSON),
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid URL\n",
			},
		},
	}

	runTests(t, tests)

}
