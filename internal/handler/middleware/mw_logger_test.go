package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerMiddleware(t *testing.T) {
	core, observedLogs := observer.New(zapcore.InfoLevel)
	testLogger := zap.New(core)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggerMiddleware(testLogger))

	router.GET("/test-endpoint", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test-endpoint?foo=bar", nil)
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "ok", rw.Body.String())

	require.Equal(t, 2, observedLogs.Len(), "Expected exactly 2 log entries")

	reqLog := observedLogs.All()[0]
	respLog := observedLogs.All()[1]

	// assert request log
	assert.Equal(t, "Request", reqLog.Message)
	assert.Equal(t, "/test-endpoint?foo=bar", reqLog.ContextMap()["uri"])
	assert.EqualValues(t, http.MethodGet, reqLog.ContextMap()["method"])
	assert.Contains(t, reqLog.ContextMap(), "duration")

	// assert response log
	assert.Equal(t, "Response", respLog.Message)
	assert.Equal(t, int64(http.StatusOK), respLog.ContextMap()["status"])
	assert.Equal(t, int64(2), respLog.ContextMap()["contentLength"])
}
