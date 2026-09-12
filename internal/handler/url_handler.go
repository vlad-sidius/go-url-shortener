package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

type urlService interface {
	CreateShortCode(originalURL string) (string, error)
	ResolveOriginalURL(hash string) (string, bool)
}

type URLHandler struct {
	log     *zap.Logger
	service urlService
}

func NewURLHandler(log *zap.Logger, service urlService) *URLHandler {
	return &URLHandler{log, service}
}

// RegisterRoutes registers all handlers
func (h *URLHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/", h.shortenURLHandler)
	router.POST("/api/shorten", h.apiShortenURLHandler)
	router.GET("/:id", h.getURLHandler)
}

// generates hash and store url in local storage
func (h *URLHandler) shortenURLHandler(ctx *gin.Context) {
	body, err := ctx.GetRawData()
	if err != nil {
		h.log.Error("Failed to read body")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// validate url
	originalURL := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// generate short url
	shortURL, err := h.service.CreateShortCode(originalURL)
	if err != nil {
		h.log.Error("Failed to create short URL", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.String(http.StatusCreated, shortURL)
}

// generates hash and store url in local storage, use JSON serialization
func (h *URLHandler) apiShortenURLHandler(ctx *gin.Context) {
	// read body
	var request model.ShortenRequest

	body, err := ctx.GetRawData()
	if err != nil {
		h.log.Error("Failed to read body")
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse("Failed to read body"))
		return
	}

	// parse body JSON
	if err := json.Unmarshal(body, &request); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse("Failed to parse body JSON"))
		return
	}

	// validate url
	originalURL := strings.TrimSpace(request.URL)
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse("Provided URL is invalid"))
		return
	}

	// generate short url
	shortURL, err := h.service.CreateShortCode(originalURL)
	if err != nil {
		h.log.Error("Failed to create short URL", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, model.NewShortenResponse(shortURL))
}

// resolves an url by hash
func (h *URLHandler) getURLHandler(ctx *gin.Context) {
	hash := ctx.Param("id")

	originalURL, ok := h.service.ResolveOriginalURL(hash)
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}
