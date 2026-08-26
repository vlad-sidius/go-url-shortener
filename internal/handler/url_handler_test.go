package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	originalURL := `https://practicum.yandex.ru`

	testCases := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, expectedCode: http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/", strings.NewReader(originalURL))
			rw := httptest.NewRecorder()

			shortenURLHandler(rw, req)

			assert.Equal(t, tc.expectedCode, rw.Code, "Invalid status code")

			if rw.Code == http.StatusCreated {
				_, err := url.ParseRequestURI(rw.Body.String())
				assert.NoError(t, err, "Malformed URL")
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	originalURL := `https://practicum.yandex.ru`
	urlHash := `test123`

	memRepo.Put(urlHash, originalURL)

	testCases := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/"+urlHash, nil)
			rw := httptest.NewRecorder()

			getURLHandler(rw, req)

			assert.Equal(t, tc.expectedCode, rw.Code, "Invalid status code")

			if rw.Code == http.StatusTemporaryRedirect {
				location := rw.Header().Get("Location")
				assert.Equal(t, originalURL, location, "Original URL is invalid")
			}
		})
	}
}
