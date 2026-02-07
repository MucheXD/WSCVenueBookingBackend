package logger_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLog(t *testing.T) {
	// EXAMPLE OF USING THE LOGGER
	slog.Debug("Debug", "key1", "value1")
	slog.Info("Info", "key2", "value2")
	slog.Warn("Warn", "key3", "value3")
	slog.Error("Error", "key4", "value4")
	var ctx gin.Context
	ctx.Error(errors.New("error trace 1"))
	ctx.Error(errors.New("error trace 2"))
	slog.Info("Info with error traces", "errorTraces", ctx.Errors)
}
