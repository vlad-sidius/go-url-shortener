package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentType     = "Content-Type"

	ContentTypeJSON = "application/json"
	ContentTypeText = "text/html"

	EncodingGzip = "gzip"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	if w.writer == nil {
		if w.shouldCompress() {
			w.Header().Set(HeaderContentEncoding, EncodingGzip)
			w.writer = gzip.NewWriter(w.ResponseWriter)
		} else {
			return w.ResponseWriter.Write(data)
		}
	}

	return w.writer.Write(data)
}

func (w *gzipResponseWriter) shouldCompress() bool {
	status := w.Status()
	if status >= 400 {
		return false
	}

	contentType := w.Header().Get(HeaderContentType)
	if contentType == "" {
		return false
	}

	compressibleTypes := []string{ContentTypeText, ContentTypeJSON}

	for _, typ := range compressibleTypes {
		if strings.HasPrefix(contentType, typ) {
			return true
		}
	}

	return false
}

func GzipMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info(HeaderContentEncoding + ctx.GetHeader(HeaderContentEncoding))

		// try uncompress gzipped request
		if strings.Contains(ctx.GetHeader(HeaderContentEncoding), EncodingGzip) {
			gz, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				ctx.AbortWithStatusJSON(
					http.StatusBadRequest,
					model.NewErrorResponse("Invalid gzip request body"),
				)
				return
			}

			defer gz.Close()
			ctx.Request.Body = gz
		}

		// set gzip writer if needed
		if strings.Contains(ctx.GetHeader(HeaderAcceptEncoding), EncodingGzip) {
			gzipWriter := &gzipResponseWriter{
				ResponseWriter: ctx.Writer,
				writer:         nil,
			}

			ctx.Writer = gzipWriter
		}

		// Process request
		ctx.Next()

		if gw, ok := ctx.Writer.(*gzipResponseWriter); ok && gw.writer != nil {
			_ = gw.writer.Close()
		}
	}
}
