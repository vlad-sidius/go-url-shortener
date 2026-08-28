package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type urlRepo interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type URLHandler struct {
	repo urlRepo
}

func NewURLHandler(repo urlRepo) *URLHandler {
	return &URLHandler{
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
	if ctx.Request.Method != http.MethodPost {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed to read body")
		return
	}

	hashCode := h.genURLHash()
	originalURL := strings.TrimSpace(string(body))
	h.repo.Put(hashCode, originalURL)

	ctx.String(http.StatusCreated, "http://localhost:8080/"+hashCode)
}

// resolves an url by hash
func (h *URLHandler) getURLHandler(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		// only GET requests allowed
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

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
