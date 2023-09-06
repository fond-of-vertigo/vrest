package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/fond-of-vertigo/vrest"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func main() {
	tp, err := newTracerProvider("otel_via_slog", "http://localhost:4318/v1/traces")
	if err != nil {
		panic(err)
	}

	logger := slog.New(&OtelSlogHandler{parent: slog.Default().Handler(), tp: tp})

	restClient := vrest.New().
		SetBaseURL("http://www.google.de").
		SetLogger(logger)

	err = restClient.NewRequest().
		DoGet("/")
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
}

type OtelSlogHandler struct {
	parent slog.Handler
	tp     *sdktrace.TracerProvider
}

func (h *OtelSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level == slog.LevelDebug {
		return true
	}
	return h.parent.Enabled(ctx, level)
}

func (h *OtelSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "otel.action" {
			switch a.Value.String() {
			case "trace.start":
				ctxPtr := ctx.Value("otel.context.ptr")
				if p, ok := ctxPtr.(*context.Context); ok {
					*p = context.WithValue(ctx, "Hallo", 1)
				}
			}
			return false
		}
		return true
	})

	return h.parent.Handle(ctx, r)
}

func (h *OtelSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.parent.WithAttrs(attrs)
}

func (h *OtelSlogHandler) WithGroup(name string) slog.Handler {
	return h.parent.WithGroup(name)
}

func newTracerProvider(serviceName, traceURL string) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	traceClient, err := newTraceClient(traceURL)
	if err != nil {
		return nil, err
	}

	exp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)

	// Create the tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)

	// Extract the span context from the headers of the NATS message
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func newTraceClient(traceURL string) (otlptrace.Client, error) {
	turl, err := url.Parse(traceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL \"%s\": %w", traceURL, err)
	}

	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(turl.Host),
		otlptracehttp.WithURLPath(turl.Path),
	}
	if turl.Scheme == "http" {
		options = append(options, otlptracehttp.WithInsecure())
	}
	return otlptracehttp.NewClient(options...), nil
}
