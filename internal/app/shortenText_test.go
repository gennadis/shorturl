package app_test

import (
	"net/http"
	"testing"
)

func TestShortenPlainText(t *testing.T) {
	tests := []test{
		{
			name:    "Shorten valid original URL",
			request: "/",
			method:  http.MethodPost,
			body:    longURL,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				body:        shortURL,
			},
		},
		{
			name:    "Shorten invalid original URL",
			request: "/",
			method:  http.MethodPost,
			body:    "qwertyuiop",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid URL\n",
			},
		},
		{
			name:    "Shorten empty original URL",
			request: "/",
			method:  http.MethodPost,
			body:    "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid URL\n",
			},
		},
	}

	runTests(t, tests)

}
