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

func Test_Logger_Info(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging info")
		}
	}()
	logger := NewLogrusLogger()
	logger.Info("mock-info")
}

func Test_Logger_Warn(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging warn")
		}
	}()
	logger := NewLogrusLogger()
	logger.Warn("mock-warn")
}

func Test_Logger_Error(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic logging error")
		}
	}()
	logger := NewLogrusLogger()
	logger.Error(errors.New("mock-error"))
}
