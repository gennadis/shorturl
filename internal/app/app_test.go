package app_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

const (
	longUrl     = "http://example.com"
	hashPattern = "^[a-zA-Z]+$"
)

func TestShortenerApp(t *testing.T) {
	store := memstore.New()
	a := app.New(store)

	// Test URL Shortener
	shortenReq, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longUrl))
	shortenResp := httptest.NewRecorder()
	a.Multiplex(shortenResp, shortenReq)

	shortUrl, _ := url.ParseRequestURI(shortenResp.Body.String())
	urlHash := shortUrl.Path[1:]

	assertResponseCode(t, shortenResp.Code, http.StatusCreated)
	assertResponseHeader(t, "Content-Type", shortenResp.Header().Get("Content-Type"), "text/plain")
	assertHost(t, shortUrl.Host, config.Addr)
	assertHashLen(t, len(urlHash), config.HashLen)
	assertHashLettersOnly(t, hashPattern, urlHash)

	// Test URL Expander
	expandReq, _ := http.NewRequest(http.MethodGet, "/"+urlHash, nil)
	expandResp := httptest.NewRecorder()
	a.Multiplex(expandResp, expandReq)

	assertResponseCode(t, expandResp.Code, http.StatusTemporaryRedirect)
	assertResponseHeader(t, "Location", expandResp.Header().Get("Location"), longUrl)

}

func assertResponseCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response code is wrong, got %d, want %d", got, want)
	}
}

func assertResponseHeader(t testing.TB, header, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response %s header is wrong, got %s, want %s", header, got, want)
	}
}

func assertHost(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("host address is wrong, got %s, want %s", got, want)
	}
}

func assertHashLen(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("hash length is wrong, got %d, want %d", got, want)
	}
}

func assertHashLettersOnly(t testing.TB, pattern, hash string) {
	t.Helper()
	if ok, _ := regexp.MatchString(pattern, hash); !ok {
		t.Errorf("hash must contain letters only, hash %s, pattern %s", hash, pattern)
	}
}
