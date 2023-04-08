package app_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/storage/memstore"
	"github.com/stretchr/testify/assert"
)

const (
	longURL     = "http://example.com"
	hashPattern = "^[a-zA-Z]+$"
)

func TestShortenURL(t *testing.T) {
	storage := memstore.New()
	server := app.New(storage)
	server.MountHandlers()

	// Correct long URL
	shortenReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL))
	shortenResp := httptest.NewRecorder()
	server.Router.ServeHTTP(shortenResp, shortenReq)
	shortURL, _ := url.ParseRequestURI(shortenResp.Body.String())
	urlHash := shortURL.Path[1:]

	assert.Equal(t, http.StatusCreated, shortenResp.Code)
	assert.Equal(t, app.PlaingText, shortenResp.Header().Get(app.ContentType))
	assert.Equal(t, config.Addr, shortURL.Host)
	assert.Equal(t, config.HashLen, len(urlHash))
	assert.Regexp(t, hashPattern, urlHash)

	// Empty string
	emptyStringReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))
	emptyStringResp := httptest.NewRecorder()
	server.Router.ServeHTTP(emptyStringResp, emptyStringReq)

	assert.Equal(t, http.StatusBadRequest, emptyStringResp.Code)
	assert.Equal(t, app.PlaingText, shortenResp.Header().Get(app.ContentType))
	assert.Equal(t, "Invalid URL\n", emptyStringResp.Body.String())

	// Invalid URL
	invalidURLReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("123"))
	invalidURLResp := httptest.NewRecorder()
	server.Router.ServeHTTP(invalidURLResp, invalidURLReq)

	assert.Equal(t, http.StatusBadRequest, invalidURLResp.Code)
	assert.Equal(t, "Invalid URL\n", invalidURLResp.Body.String())

}

func TestExpandURL(t *testing.T) {
	storage := memstore.New()
	server := app.New(storage)
	server.MountHandlers()

	// Correct Hash
	shortenReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL))
	shortenResp := httptest.NewRecorder()
	server.Router.ServeHTTP(shortenResp, shortenReq)
	shortURL, _ := url.ParseRequestURI(shortenResp.Body.String())
	urlHash := shortURL.Path[1:]

	expandReq, _ := http.NewRequest(http.MethodGet, "/"+urlHash, nil)
	expandResp := httptest.NewRecorder()
	server.Router.ServeHTTP(expandResp, expandReq)

	assert.Equal(t, http.StatusTemporaryRedirect, expandResp.Code)
	assert.Equal(t, longURL, expandResp.Header().Get(app.Location))

	// Invalid Hash
	invalidHashReq, _ := http.NewRequest(http.MethodGet, "/"+"invalidHash", nil)
	invalidHashResp := httptest.NewRecorder()
	server.Router.ServeHTTP(invalidHashResp, invalidHashReq)

	assert.Equal(t, http.StatusBadRequest, invalidHashResp.Code)
	assert.Equal(t, "", invalidHashResp.Header().Get(app.Location))

}
