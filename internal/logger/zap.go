package logger

import (
	"context"
	"os"
	"path/filepath"
	"sync"

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

// Config holds the configuration for the logger
type Config struct {
	Level      string
	OutputPath string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool // compress old files
	Debug      bool // development mode
}

// NewZapLogger creates a new logger instance
func NewZapLogger(cfg Config) (Logger, error) {
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

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		parseLevel(cfg.Level),
	)

	// Create logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	if cfg.Debug {
		logger = logger.WithOptions(zap.Development())
	}

	return &zapLogger{
		logger: logger,
		writer: writer,
	}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, convertFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, convertFields(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, convertFields(fields...)...)
}

func (l *zapLogger) WithFields(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(convertFields(fields...)...),
		writer: l.writer,
	}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	// Extract request ID from context if available
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return l.WithFields(NewField("request_id", reqID))
	}
	return l
}

// Close ensures all log files are properly closed
func (l *zapLogger) Close() error {
	var err error
	l.Do(func() {
		// First sync the zap logger
		if err = l.logger.Sync(); err != nil {
			return
		}
		// Then close the underlying writer
		if l.writer != nil {
			err = l.writer.Close()
		}
	})
	return err
}

// Helper functions

func convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func parseLevel(level string) zapcore.Level {
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
