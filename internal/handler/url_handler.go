package handler

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	repo "github.com/vlad-sidius/go-url-shortener/internal/repository"
)

var memRepo = repo.NewMemURLRepo()

// RegisterRoutes registers all handlers
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", shortenURLHandler)
	mux.HandleFunc("GET /{id}", getURLHandler)
}

// generates hash and store url in local storage
func shortenURLHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Failed to read body", http.StatusBadRequest)
		return
	}

	req.Body.Close()

	originalURL := strings.TrimSpace(string(body))
	shortURL := genURLHash()
	memRepo.Put(shortURL, originalURL)

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("http://localhost:8080/" + shortURL))
}

// resolves an url by hash
func getURLHandler(rw http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")

	originalURL, ok := memRepo.Get(shortURL)
	if !ok {
		http.NotFound(rw, req)
		return
	}

	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func genURLHash() string {
	bytes := make([]byte, 6) // 6 bytes → ~8 symbols in base64
	rand.Read(bytes)         // no error handling is necessary, as Read always succeeds
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}
