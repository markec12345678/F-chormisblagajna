package logger

import (
	"testing"
)

func TestNewZeroLog(t *testing.T) {
	l := NewZeroLog()
	_ = l
}

func TestZeroLog_Info(t *testing.T) {
	l := NewZeroLog()
	l.Info("test info message")
}

func TestZeroLog_Warning(t *testing.T) {
	l := NewZeroLog()
	l.Warning("test warning message")
}

func TestZeroLog_Error(t *testing.T) {
	l := NewZeroLog()
	l.Error("test error message")
}

func TestZeroLog_ImplementsILogger(t *testing.T) {
	var _ ILogger = &ZeroLog{}
}

func TestZeroLog_MethodsWithFields(t *testing.T) {
	l := NewZeroLog()
	l.Info("info with fields", "key1", "value1", "key2", 42)
	l.Warning("warning with fields", "key", "value")
	l.Error("error with fields", "count", 3)
}
