package logger

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/jonesrussell/dashboard/internal/logger/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// requestIDKey is the context key for request IDs
	requestIDKey = contextKey("request_id")
)

type zapLogger struct {
	logger    *zap.Logger
	writer    *lumberjack.Logger
	sync.Once // For safe cleanup
}

// NewZapLogger creates a new logger instance
func NewZapLogger(cfg types.Config) (types.Logger, error) {
	// Create the log directory if it doesn't exist
	if cfg.OutputPath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.OutputPath), 0o755); err != nil {
			return nil, err
		}
	}

	// Configure log rotation
	writer := &lumberjack.Logger{
		Filename:   cfg.OutputPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		getZapLevel(cfg.Level),
	)

	// Create logger
	logger := &zapLogger{
		logger: zap.New(core),
		writer: writer,
	}

	return logger, nil
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug implements Logger
func (l *zapLogger) Debug(msg string, fields ...types.Field) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

// Info implements Logger
func (l *zapLogger) Info(msg string, fields ...types.Field) {
	l.logger.Info(msg, convertFields(fields...)...)
}

// Warn implements Logger
func (l *zapLogger) Warn(msg string, fields ...types.Field) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

// Error implements Logger
func (l *zapLogger) Error(msg string, fields ...types.Field) {
	l.logger.Error(msg, convertFields(fields...)...)
}

// Fatal implements Logger
func (l *zapLogger) Fatal(msg string, fields ...types.Field) {
	l.logger.Fatal(msg, convertFields(fields...)...)
}

// WithFields implements Logger
func (l *zapLogger) WithFields(fields ...types.Field) types.Logger {
	return &zapLogger{
		logger: l.logger.With(convertFields(fields...)...),
		writer: l.writer,
	}
}

// WithContext implements Logger
func (l *zapLogger) WithContext(ctx context.Context) types.Logger {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return l.WithFields(types.NewField("request_id", reqID))
	}
	return l
}

// Close implements Logger
func (l *zapLogger) Close() error {
	var err error
	l.Do(func() {
		err = l.writer.Close()
	})
	return err
}

// convertFields converts our Field type to zap.Field
func convertFields(fields ...types.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
