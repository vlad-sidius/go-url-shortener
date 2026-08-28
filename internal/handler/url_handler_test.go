package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
)

func TestShortenURLHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init dependencies
	conf := config.Config{Address: "localhost:8080", BaseURL: `http://localhost:8000`}
	memRepo := repository.NewMemURLRepo()
	urlHandler := NewURLHandler(&conf, memRepo)

	// create router
	router := gin.New()
	urlHandler.RegisterRoutes(router)

	// create server
	server := httptest.NewServer(router)
	defer server.Close()

	// create client
	client := resty.New()

	originalURL := `https://practicum.yandex.ru`
	invalidURL := `djsljfijiwjfiew7289713789723198731`

	testCases := []struct {
		method       string
		expectedCode int
		body         string
	}{
		{method: http.MethodGet, expectedCode: http.StatusNotFound},
		{method: http.MethodPut, expectedCode: http.StatusNotFound},
		{method: http.MethodDelete, expectedCode: http.StatusNotFound},
		{method: http.MethodPost, expectedCode: http.StatusBadRequest, body: invalidURL},
		{method: http.MethodPost, expectedCode: http.StatusCreated, body: originalURL},
	}

	remoteURL := server.URL + "/"

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			var resp *resty.Response
			var err error

			if tc.method == http.MethodPost {
				resp, err = client.R().SetBody(tc.body).Post(remoteURL)
			} else {
				resp, err = client.R().Execute(tc.method, remoteURL)
			}

			assert.NoError(t, err, "Should be success")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Invalid status code")

			if tc.method == http.MethodPost && resp.StatusCode() == http.StatusCreated {
				shortURL := strings.TrimSpace(string(resp.Body()))
				assert.Contains(t, shortURL, conf.BaseURL)

				_, err := url.ParseRequestURI(shortURL)
				assert.NoError(t, err, "Response body is not a valid URL: %q", shortURL)
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init dependencies
	conf := config.Config{Address: "localhost:8080", BaseURL: `http://localhost:8000`}
	memRepo := repository.NewMemURLRepo()
	urlHandler := NewURLHandler(&conf, memRepo)

	// pre-populate the repository
	urlHash := "test123"
	originalURL := "https://practicum.yandex.ru"
	memRepo.Put(urlHash, originalURL)

	// create router
	router := gin.New()
	urlHandler.RegisterRoutes(router)

	// create server
	server := httptest.NewServer(router)
	defer server.Close()

	// create a client that doesn't follow redirects automatically
	client := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(0))

	testCases := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodPost, expectedCode: http.StatusNotFound},
		{method: http.MethodPut, expectedCode: http.StatusNotFound},
		{method: http.MethodDelete, expectedCode: http.StatusNotFound},
		{method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect},
	}

	remoteURL := server.URL + "/" + urlHash

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			var resp *resty.Response
			var err error

			resp, err = client.R().Execute(tc.method, remoteURL)

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Invalid status code")

			if tc.method == http.MethodGet && tc.expectedCode == http.StatusTemporaryRedirect {
				assert.Error(t, err, "Should redirect")
				assert.Contains(t, err.Error(), "stopped after 0 redirects")

				location := resp.Header().Get("Location")
				assert.Equal(t, originalURL, location, "Original URL is invalid")
			} else {
				assert.NoError(t, err, "Should be success")
			}
		})
	}
}
