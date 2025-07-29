package log_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetAndGetLoggerFromContext(t *testing.T) {
	l := zap.NewNop().Sugar()
	ctxWithLogger := log.WithContext(context.Background(), l)
	gotLogger := log.FromContext(ctxWithLogger)
	assert.Equal(t, l, gotLogger)
}
