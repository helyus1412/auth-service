package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/helyus1412/auth-service/cmd/routes"
	"github.com/helyus1412/auth-service/config"
	"github.com/helyus1412/auth-service/pkg/databases"
	"github.com/helyus1412/auth-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.uber.org/zap/zapcore"
)

var (
	serviceName         = config.GlobalEnv.AppName
	serviceVersion      = config.GlobalEnv.AppVersion
	serviceIsProduction = config.GlobalEnv.IsProduction
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tp, mp, err := newOtelProviders(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize tracer provider: %v", err))
	}
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		// Configures the global TextMapPropagator to use the
		// W3C Trace Context format for propagating trace context.
		// This configuration ensures that spans have the correct
		// parent-child relationship within a trace.
		autoprop.NewTextMapPropagator(),
	))
	tracer := otel.Tracer("main-tracer")

	e := echo.New()
	e.HideBanner = true
	e.Use(otelecho.Middleware(serviceName))

	// Initialize our custom Logger
	logger, err := logger.New(logger.Config{
		ServiceName:  serviceName,
		LogLevel:     zapcore.DebugLevel,
		IsProduction: serviceIsProduction,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	db, err := databases.InitPostgre()
	if err != nil {
		log.Fatalf("failed init postgre: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db.DB, "./pkg/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrated successfully")

	routes.InitRoutes(e, db, tracer, logger)

	go func() {
		port := fmt.Sprintf(":%d", config.GlobalEnv.HTTPPort)
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx,
				"main",
				"Setup",
				"HTTP server shutdown unexpectedly",
				err,
				zapcore.Field{Key: "listenerPort", Type: zapcore.StringType, String: port},
			)
		}
	}()

	// --- Wait for Shutdown Signal ---
	// Block here until the context is cancelled (e.g., by Ctrl+C).
	<-ctx.Done()
	logger.Info(ctx, "main", "Shutdown", "Shutdown signal received. Starting graceful shutdown...")

	// Create a new context for shutdown with a timeout.
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()

	// Shutdown the HTTP server.
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx,
			"main",
			"Shutdown",
			"HTTP server shutdown error",
			err,
		)
	}

	//flush the log buffer on exit
	logger.Sync()

	if err := tp.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx,
			"main",
			"Shutdown",
			"Error shutting down tracer provider",
			err,
			zapcore.Field{},
		)
	}

	if err := mp.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx,
			"main",
			"Shutdown",
			"Error shutting down metrics provider",
			err,
			zapcore.Field{},
		)
	}
}

func newOtelProviders(ctx context.Context) (*sdktrace.TracerProvider, *metric.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
	)

	// ---- Traces ----
	var traceExporter sdktrace.SpanExporter
	var err error
	if serviceIsProduction {
		traceExporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", config.GlobalEnv.OtelTraceHost, config.GlobalEnv.OtelTracePort)),
			otlptracehttp.WithInsecure(),
		)
	} else {
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
	if err != nil {
		return nil, nil, fmt.Errorf("trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// ---- Metrics ----
	var metricExporter metric.Exporter
	if serviceIsProduction {
		metricExporter, err = otlpmetrichttp.New(
			ctx,
			otlpmetrichttp.WithEndpoint(fmt.Sprintf("%s:%d", config.GlobalEnv.OtelMetricsHost, config.GlobalEnv.OtelMetricsPort)),
			otlpmetrichttp.WithInsecure(),
		)
	} else {
		metricExporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}
	if err != nil {
		return nil, nil, fmt.Errorf("metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)

	return tp, mp, nil
}
