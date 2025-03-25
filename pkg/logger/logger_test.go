package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rs/zerolog"
)

type testHook struct {
	called bool
}

func (h *testHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	h.called = true
}

func TestLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := logger.New(logger.WithWriters(buf), logger.WithLogLevel("debug"))
	logger.Debug("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message', got %v", output)
	}
}

func TestLoggerWithLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		logFn func(l logger.Logger)
	}{
		{"debug", "debug", func(l logger.Logger) { l.Debug("test message") }},
		{"info", "info", func(l logger.Logger) { l.Info("test message") }},
		{"warn", "warn", func(l logger.Logger) { l.Warn("test message") }},
		{"error", "error", func(l logger.Logger) { l.Error("test message") }},
		{"unknown", "unknown", func(l logger.Logger) { l.Info("test message") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			log := logger.New(logger.WithWriters(buf), logger.WithLogLevel(tt.level))
			tt.logFn(log)
			output := buf.String()

			expectedLevel := tt.level
			if expectedLevel == "unknown" {
				expectedLevel = "info" // default level is "info"
			}

			if !strings.Contains(output, expectedLevel) {
				t.Errorf("Expected log level %v, but it was not present in output: %v", expectedLevel, output)
			}
		})
	}
}

func TestWithLogLevel(t *testing.T) {
	level := "debug"
	opts := &logger.Options{}
	logger.WithLogLevel(level)(opts)
	if opts.LogLevel != level {
		t.Errorf("WithLogLevel() didn't set LogLevel to %s", level)
	}
}

func TestLoggerWithHook(t *testing.T) {
	buf := &bytes.Buffer{}
	th := &testHook{}
	log := logger.New(logger.WithWriters(buf), logger.WithLogLevel("debug"), logger.WithHook(th))
	log.Debug("test message")

	if !th.called {
		t.Errorf("Expected the hook to be called, but it wasn't")
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message', got %v", output)
	}
}

func TestNew(t *testing.T) {
	logger := logger.New(logger.WithLogLevel("debug"))
	if logger == nil {
		t.Errorf("New() returned nil")
	}
}
