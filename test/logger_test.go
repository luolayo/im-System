package logger

import (
	Logger "im-System/logger"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger := Logger.NewLogger(Logger.DebugLevel)
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	defer logger.Sync()
}

func TestLogger_Info(t *testing.T) {
	logger := Logger.NewLogger(Logger.InfoLevel)
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	defer logger.Sync()
}

func TestLogger_Warn(t *testing.T) {
	logger := Logger.NewLogger(Logger.WarnLevel)
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	defer logger.Sync()
}

func TestLogger_Error(t *testing.T) {
	logger := Logger.NewLogger(Logger.ErrorLevel)
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	defer logger.Sync()
}
