package pkg

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"

	"go.opentelemetry.io/otel/trace"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Tracer struct {
}

type TracerProviderConfig struct {
	JaegerEndpoint string
	ServiceName    string
	ServiceVersion string
	Environment    string
	Disabled       bool
}

type TracerProvider struct {
	provider trace.TracerProvider
}

func NewProvider(config TracerProviderConfig) (TracerProvider, error) {
	if config.Disabled {
		return TracerProvider{provider: trace.NewNoopTracerProvider()}, nil
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerEndpoint)),
	)

	if err != nil {
		return TracerProvider{}, err
	}

	prv := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			sdkresource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ServiceName),
				semconv.ServiceVersionKey.String(config.ServiceVersion),
				semconv.DeploymentEnvironmentKey.String(config.Environment),
			)),
	)

	otel.SetTracerProvider(prv)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return TracerProvider{provider: prv}, nil
}

func (p *TracerProvider) Close(ctx context.Context) error {
	if prv, ok := p.provider.(*sdktrace.TracerProvider); ok {
		return prv.Shutdown(ctx)
	}

	return nil
}

func NewRequestFromHeader(headers http.Header, ctx context.Context) context.Context {
	var (
		requestID = headers.Get("request-id")
	)

	if requestID == "" {
		requestID = uuid.NewString()
	}

	ctx = context.WithValue(ctx, "request-id", requestID)
	return ctx
}

func NewSpan(ctx context.Context, name string, cus SpanCustomiser) (context.Context, trace.Span) {
	if ctx.Value("request-id") == nil {
		ctx = NewRequestFromHeader(http.Header{}, ctx)
	}

	if cus == nil {
		ctx, span := otel.Tracer("").Start(ctx, name)

		AddSpanTags(span, map[string]string{
			"request-id": ctx.Value("request-id").(string),
		})

		return ctx, span
	}

	custom := cus.customise()

	ctx, span := otel.Tracer("").Start(ctx, name, custom)

	AddSpanTags(span, map[string]string{
		"request-id": ctx.Value("request-id").(string),
	})

	return ctx, span
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func AddSpanTags(span trace.Span, tags map[string]string) {
	list := make([]attribute.KeyValue, len(tags))

	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	span.SetAttributes(list...)
}

func AddSpanEvents(span trace.Span, name string, events map[string]string) {
	list := make([]trace.EventOption, len(events))

	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}

	span.AddEvent(name, list...)
}

func AddSpanError(span trace.Span, err error) {
	span.RecordError(err)
}

func FailSpan(span trace.Span, msg string) {
	span.SetStatus(codes.Error, msg)
}

type SpanCustomiser interface {
	customise() trace.SpanOption
}
