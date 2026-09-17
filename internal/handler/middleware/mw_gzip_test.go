package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupRouter(log *zap.Logger, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GzipMiddleware(log))
	router.GET("/test", handler)
	return router
}

func TestGzipMiddleware_CompressesJSONResponse(t *testing.T) {
	log := zap.NewNop()
	router := setupRouter(log, func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, `{"message":"hello world"}`)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "gzip", rw.Header().Get("Content-Encoding"))

	gz, err := gzip.NewReader(rw.Body)
	require.NoError(t, err)
	defer gz.Close()

	body, err := io.ReadAll(gz)
	require.NoError(t, err)
	assert.Contains(t, string(body), "hello world")
}

func TestGzipMiddleware_NoCompressionWithoutHeader(t *testing.T) {
	log := zap.NewNop()
	router := setupRouter(log, func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, `{"message":"hello"}`)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Empty(t, rw.Header().Get("Content-Encoding"))
	assert.Contains(t, rw.Body.String(), "hello")
}

func TestGzipMiddleware_NoCompressionForErrorResponse(t *testing.T) {
	log := zap.NewNop()
	router := setupRouter(log, func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusInternalServerError, `{"error":"fail"}`)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.Empty(t, rw.Header().Get("Content-Encoding"))
}

func TestGzipMiddleware_DecompressesGzipRequestBody(t *testing.T) {
	log := zap.NewNop()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GzipMiddleware(log))
	router.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		require.NoError(t, err)
		c.String(http.StatusOK, string(body))
	})

	// Create gzipped body
	var buf strings.Builder
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte("hello decompressed"))
	require.NoError(t, err)
	require.NoError(t, gz.Close())

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(buf.String()))
	req.Header.Set("Content-Encoding", "gzip")
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "hello decompressed", rw.Body.String())
}
