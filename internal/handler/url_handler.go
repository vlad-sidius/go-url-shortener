package handler

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type urlService interface {
	CreateShortCode(originalURL string) (string, error)
	ResolveOriginalURL(hash string) (string, bool)
}

type URLHandler struct {
	service urlService
}

func NewURLHandler(service urlService) *URLHandler {
	return &URLHandler{service}
}

// RegisterRoutes registers all handlers
func (h *URLHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/", h.shortenURLHandler)
	router.GET("/:id", h.getURLHandler)
}

// generates hash and store url in local storage
func (h *URLHandler) shortenURLHandler(ctx *gin.Context) {
	body, err := ctx.GetRawData()
	if err != nil {
		log.Println("Failed to read body")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// validate url
	originalURL := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		log.Println("Provided URL is invalid")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	shortURL, err := h.service.CreateShortCode(originalURL)
	if err != nil {
		log.Printf("Failed to create short URL: %v\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.String(http.StatusCreated, shortURL)
}

// resolves an url by hash
func (h *URLHandler) getURLHandler(ctx *gin.Context) {
	hash := ctx.Param("id")

	originalURL, ok := h.service.ResolveOriginalURL(hash)
	if !ok {
		log.Printf("Original URL for hash '%s' was not found\n", hash)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}
