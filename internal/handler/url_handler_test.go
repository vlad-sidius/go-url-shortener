package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
)

func TestShortenURLHandler(t *testing.T) {
	memRepo := repository.NewMemURLRepo()
	urlHandler := NewURLHandler(memRepo)

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
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Handle(tc.method, "/", urlHandler.shortenURLHandler)

			var req *http.Request

			if tc.method == http.MethodPost {
				req = httptest.NewRequest(tc.method, "/", strings.NewReader(originalURL))
				req.Header.Set("Content-Type", "text/plain")
			} else {
				req = httptest.NewRequest(tc.method, "/", nil)
			}

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)

			assert.Equal(t, tc.expectedCode, rw.Code, "Invalid status code")

			if rw.Code == http.StatusCreated {
				_, err := url.ParseRequestURI(rw.Body.String())
				assert.NoError(t, err, "Malformed URL")
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	memRepo := repository.NewMemURLRepo()
	urlHandler := NewURLHandler(memRepo)

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
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Handle(tc.method, "/:id", urlHandler.getURLHandler)

			req := httptest.NewRequest(tc.method, "/"+urlHash, nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)

			assert.Equal(t, tc.expectedCode, rw.Code, "Invalid status code")

			if rw.Code == http.StatusTemporaryRedirect {
				location := rw.Header().Get("Location")
				assert.Equal(t, originalURL, location, "Original URL is invalid")
			}
		})
	}
}
