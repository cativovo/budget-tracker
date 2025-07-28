package ctxvalue_test

import (
	"context"
	"testing"

	"github.com/cativovo/budget-tracker/internal/ctxvalue"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetAndGetLoggerFromContext(t *testing.T) {
	l := zap.NewNop().Sugar()
	ctxWithLogger := ctxvalue.ContextWithLogger(context.Background(), l)
	gotLogger := ctxvalue.LoggerFromContext(ctxWithLogger)
	assert.Equal(t, l, gotLogger)
}
