package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
)

type urlRepo interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type URLHandler struct {
	conf *config.Config
	repo urlRepo
}

func NewURLHandler(conf *config.Config, repo urlRepo) *URLHandler {
	return &URLHandler{
		conf: conf,
		repo: repo,
	}
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
		ctx.String(http.StatusBadRequest, "Failed to read body")
		return
	}

	// validate url
	originalURL := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		ctx.String(http.StatusBadRequest, "Provided URL is invalid")
		return
	}

	hashCode := h.genURLHash()
	h.repo.Put(hashCode, originalURL)

	ctx.String(http.StatusCreated, h.conf.BaseURL+"/"+hashCode)
}

// resolves an url by hash
func (h *URLHandler) getURLHandler(ctx *gin.Context) {
	hash := ctx.Param("id")

	originalURL, ok := h.repo.Get(hash)
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (h *URLHandler) genURLHash() string {
	bytes := make([]byte, 6) // 6 bytes → ~8 symbols in base64
	rand.Read(bytes)         // no error handling is necessary, as Read always succeeds
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}
