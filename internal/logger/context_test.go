package logger_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetAndGetLoggerFromContext(t *testing.T) {
	l := zap.NewNop().Sugar()

	ctxWithLogger := logger.ContextWithLogger(context.Background(), l)
	gotLogger := logger.FromContext(ctxWithLogger)
	assert.Equal(t, l, gotLogger)
}
