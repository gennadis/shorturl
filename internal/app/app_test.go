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

	shortenReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL))
	shortenResp := httptest.NewRecorder()
	server.Router.ServeHTTP(shortenResp, shortenReq)

	shortURL, _ := url.ParseRequestURI(shortenResp.Body.String())
	urlHash := shortURL.Path[1:]

	assert.Equal(t, shortenResp.Code, http.StatusCreated)
	assert.Equal(t, shortenResp.Header().Get(app.ContentType), app.PlaingText)
	assert.Equal(t, shortURL.Host, config.Addr)
	assert.Equal(t, len(urlHash), config.HashLen)
	assert.Regexp(t, hashPattern, urlHash)
}

func TestExpandURL(t *testing.T) {
	storage := memstore.New()
	server := app.New(storage)
	server.MountHandlers()

	shortenReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL))
	shortenResp := httptest.NewRecorder()
	server.Router.ServeHTTP(shortenResp, shortenReq)

	shortURL, _ := url.ParseRequestURI(shortenResp.Body.String())
	urlHash := shortURL.Path[1:]

	expandReq, _ := http.NewRequest(http.MethodGet, "/"+urlHash, nil)
	expandResp := httptest.NewRecorder()
	server.Router.ServeHTTP(expandResp, expandReq)

	assert.Equal(t, expandResp.Code, http.StatusTemporaryRedirect)
	assert.Equal(t, expandResp.Header().Get(app.Location), longURL)
}
