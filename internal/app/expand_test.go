package app_test

import (
	"net/http"
	"testing"
)

func TestExpand(t *testing.T) {
	tests := []test{
		{
			name:    "Prepare valid original URL",
			request: "/",
			method:  http.MethodPost,
			body:    longURL,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				body:        shortenURL(longURL),
			},
		},
		{
			name:    "Expand valid hash",
			request: "/eebc6e3",
			method:  http.MethodGet,
			body:    "",
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "",
				body:        "",
			},
		},
		{
			name:    "Expand invalid hash",
			request: "/qwertyuiop",
			method:  http.MethodGet,
			body:    "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "wrong hash provided\n",
			},
		},
	}

	runTests(t, tests)

}
