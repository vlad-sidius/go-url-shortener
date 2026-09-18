package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Initialize(t *testing.T) {
	t.Run("successfully creates a logger with info level", func(t *testing.T) {
		log, err := NewLogger()

		require.NoError(t, err, "NewLogger should not return an error for valid LogLevel")
		require.NotNil(t, log, "Logger instance should not be nil")

		assert.True(t, log.Core().Enabled(zapcore.InfoLevel), "Logger should be enabled for Info level")
		assert.False(t, log.Core().Enabled(zapcore.DebugLevel), "Logger should NOT be enabled for Debug level")

		// Clean up
		defer func() {
			_ = log.Sync()
		}()
	})
}
