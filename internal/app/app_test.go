package app_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gennadis/shorturl/config"
	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	longURL = "http://amazon.com.tr"
)

type want struct {
	statusCode  int
	contentType string
	body        string
}

type test struct {
	name    string
	request string
	method  string
	body    string
	want    want
}

func TestMisc(t *testing.T) {
	tests := []test{
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
	}

	runTests(t, tests)
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

func runTests(t *testing.T, tests []test) {
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

func shortenURL(url string) string {
	hash := app.GenerateHash(url, config.HashLen)
	return fmt.Sprintf("http://%s/%s", config.Addr, hash)

}
