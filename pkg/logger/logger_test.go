package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger()
	assert.NotNil(t, logger)
}

func TestLoggerInfo(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging info")
		}
	}()
	logger := NewLogrusLogger()
	logger.Info("mock-info")
}

func TestLoggerWarn(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging warn")
		}
	}()
	logger := NewLogrusLogger()
	logger.Warn("mock-warn")
}

func TestLoggerError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging error")
		}
	}()
	logger := NewLogrusLogger()
	logger.Error(errors.New("mock-error"))
}
