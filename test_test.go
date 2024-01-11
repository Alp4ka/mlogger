package mlogger

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	cfg := Config{
		Source: "grizha",
	}

	logger, err := NewProduction(ctx, cfg)
	if err != nil {
		panic(err)
	}

	logger.Info("pisya-popa")
	//assert.Equal(t, )
}
