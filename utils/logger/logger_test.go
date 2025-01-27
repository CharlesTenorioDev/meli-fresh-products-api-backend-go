package logger

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	Info("This is an info message")
	Error("This is an error message", nil)
}

func TestGetOutputLogs(t *testing.T) {
	t.Run("should return stdout when LOG_OUTPUT is not set", func(t *testing.T) {
		os.Unsetenv(LogOutput)
		got := getOutputLogs()
		want := "stdout"
		if got != want {
			t.Errorf("getOutputLogs() = %v, want %v", got, want)
		}
	})

	t.Run("should return the value of LOG_OUTPUT when it is set", func(t *testing.T) {
		os.Setenv(LogOutput, "test_output")
		got := getOutputLogs()
		want := "test_output"
		if got != want {
			t.Errorf("getOutputLogs() = %v, want %v", got, want)
		}
	})

	t.Run("should return lowercase value of LOG_OUTPUT", func(t *testing.T) {
		os.Setenv(LogOutput, "TeSt_OutPuT")
		got := getOutputLogs()
		want := "test_output"
		if got != want {
			t.Errorf("getOutputLogs() = %v, want %v", got, want)
		}
	})
}

func TestGetLevelLogs(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		want     zapcore.Level
	}{
		{"should return InfoLevel when LOG_LEVEL is not set", "", zapcore.InfoLevel},
		{"should return InfoLevel when LOG_LEVEL is invalid", "invalid", zapcore.InfoLevel},
		{"should return InfoLevel when LOG_LEVEL is info", "info", zapcore.InfoLevel},
		{"should return ErrorLevel when LOG_LEVEL is error", "error", zapcore.ErrorLevel},
		{"should return DebugLevel when LOG_LEVEL is debug", "debug", zapcore.DebugLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(LogLevel, tt.logLevel)
			if got := getLevelLogs(); got != tt.want {
				t.Errorf("getLevelLogs() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("should return lowercase value of LOG_LEVEL", func(t *testing.T) {
		os.Setenv(LogLevel, "DeBuG")
		got := getLevelLogs()
		want := zapcore.DebugLevel
		if got != want {
			t.Errorf("getLevelLogs() = %v, want %v", got, want)
		}
	})
}
