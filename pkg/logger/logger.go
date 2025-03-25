package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	zerologadapter "logur.dev/adapter/zerolog"
)

type (
	// Logger defines the interface for a logger.
	// Deprecated: Will be replaced with slog.Logger in the future.
	Logger interface {
		Trace(msg string, fields ...map[string]interface{})
		Debug(msg string, fields ...map[string]interface{})
		Info(msg string, fields ...map[string]interface{})
		Warn(msg string, fields ...map[string]interface{})
		Error(msg string, fields ...map[string]interface{})
	}

	// Options struct defines the logger options.
	// LogLevel sets the level of logs to show (trace, debug, info, warn, error).
	// Writer sets the writer to write logs to.
	// Hooks sets the hooks to be used by the logger.
	Options struct {
		LogLevel string
		Writers  []io.Writer
		Hooks    []zerolog.Hook
	}

	// Option defines a function which sets an option on the Options struct.
	Option func(*Options)
)

// WithLogLevel sets the log level on the Options struct.
func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.LogLevel = level
	}
}

// WithWriter sets the writer on the Options struct.
func WithWriters(writers ...io.Writer) Option {
	return func(o *Options) {
		o.Writers = writers
	}
}

// WithHook adds a hook to the logger.
func WithHook(hook zerolog.Hook) Option {
	return func(o *Options) {
		o.Hooks = append(o.Hooks, hook)
	}
}

// New creates a new Logger based on provided Options.
func New(opts ...Option) Logger {
	options := &Options{
		Writers: []io.Writer{os.Stdout},
	}

	for _, opt := range opts {
		opt(options)
	}

	logger := zerolog.New(getWriter(options)).With().Timestamp().Logger()
	for _, hook := range options.Hooks {
		logger = logger.Hook(hook)
	}

	zerolog.SetGlobalLevel(parseLogLevel(options.LogLevel))
	return zerologadapter.New(logger)
}

// NoOp returns a no-operation Logger which doesn't perform any logging operations.
func NoOp() Logger {
	return zerologadapter.New(zerolog.Nop())
}

// PrettyWriter returns a ConsoleWriter which writes logs in a human-readable format.
func PrettyWriter() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFieldFormat}
}

// parseLogLevel parses the log level string and returns the corresponding zerolog.Level.
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// getWriter creates and returns a writer based on the pretty flag and multiple writers.
// If pretty is true, it adds a ConsoleWriter to the writers slice.
// It returns a MultiLevelWriter if multiple writers are provided.
func getWriter(opts *Options) io.Writer {
	if len(opts.Writers) > 1 {
		return zerolog.MultiLevelWriter(opts.Writers...)
	}

	return opts.Writers[0]
}
