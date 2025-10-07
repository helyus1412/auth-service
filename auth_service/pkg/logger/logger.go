package logger

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a custom logger that wraps zap.Logger.
type Logger struct {
	zapLogger *zap.Logger
}

// Config holds the configuration for the logger.
type Config struct {
	ServiceName  string
	LogLevel     zapcore.Level // e.g., zapcore.InfoLevel, zapcore.DebugLevel
	IsProduction bool
}

// New initializes a new Logger.
// This is the corrected implementation.
func New(config Config) (*Logger, error) {
	var zapConfig zap.Config
	var encoder zapcore.Encoder

	// Step 1: Configure the encoder
	if config.IsProduction {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		encoder = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Pretty colors for dev
		encoder = zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
	}

	zapConfig.Level = zap.NewAtomicLevelAt(config.LogLevel)

	// Step 2: Configure the writer (make it asynchronous)
	// We use os.Stdout, but in a real app, you might use a file.
	writer := zapcore.AddSync(os.Stdout)
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS:            writer,
		Size:          4096,             // 4KB buffer
		FlushInterval: 30 * time.Second, // Flush every 30 seconds
	}

	// Step 3: Create the Core that ties the encoder, writer, and log level together.
	core := zapcore.NewCore(encoder, asyncWriter, zapConfig.Level)

	// Step 4: Create the final logger, adding options like caller and stacktrace.
	// AddCallerSkip(1) is used because we've wrapped the logger in our own methods.
	finalLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).With(
		zap.String("service.name", config.ServiceName),
	)

	return &Logger{zapLogger: finalLogger}, nil
}

// Sync flushes any buffered log entries.
// It's crucial to call this before the application exits.
func (l *Logger) Sync() {
	_ = l.zapLogger.Sync()
}

// withTraceContext adds trace_id and span_id from the context to the log fields.
func (l *Logger) withTraceContext(ctx context.Context, fields []zap.Field) []zap.Field {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		fields = append(fields,
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return fields
}

// log is the internal logging method.
func (l *Logger) log(ctx context.Context, level zapcore.Level, function string, condition string, message string, fields ...zap.Field) {
	// Add our standard structured fields
	standardFields := []zap.Field{
		zap.String("function", function),
		zap.String("condition", condition),
	}

	// Add trace context if available
	allFields := l.withTraceContext(ctx, append(standardFields, fields...))

	// Write the log entry
	if ce := l.zapLogger.Check(level, message); ce != nil {
		ce.Write(allFields...)
	}
}

// Info logs a message at the info level.
func (l *Logger) Info(ctx context.Context, function, condition, message string, fields ...zap.Field) {
	l.log(ctx, zapcore.InfoLevel, function, condition, message, fields...)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(ctx context.Context, function, condition, message string, fields ...zap.Field) {
	l.log(ctx, zapcore.WarnLevel, function, condition, message, fields...)
}

// Error logs a message at the error level, including the error details.
func (l *Logger) Error(ctx context.Context, function, condition, message string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	l.log(ctx, zapcore.ErrorLevel, function, condition, message, allFields...)
}
