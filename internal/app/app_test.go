package app_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/storage/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	longURL = "http://amazon.com.tr"
)

func TestHandlers(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		want    want
	}{{
		name:    "Shorten valid original URL",
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
		{
			name:    "Use wrong method",
			request: "/",
			method:  http.MethodGet,
			body:    "",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				body:        "",
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

	storage := memstore.New()
	server := app.New(storage)
	server.MountHandlers()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.Router
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := makeRequest(t, ts, tt.method, tt.request, tt.body)

			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func shortenURL(url string) string {
	hash := app.GenerateHash(url, config.HashLen)
	return fmt.Sprintf("http://%s/%s", config.Addr, hash)

}
