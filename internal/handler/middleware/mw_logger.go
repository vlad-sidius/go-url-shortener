package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.URL.RequestURI()
		method := c.Request.Method

		// Process request
		c.Next()

		duration := time.Since(start)

		log.Info("Request",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
		)

		log.Info("Response",
			zap.Int("status", c.Writer.Status()),
			zap.Int("contentLength", c.Writer.Size()),
		)
	}
}
