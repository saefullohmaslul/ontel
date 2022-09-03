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

// SpanFromContext returns the current span from a context. If you wish to avoid
// creating child spans for each operation and just rely on the parent span, use
// this function throughout the application. With such practise you will get
// flatter span tree as opposed to deeper version. You can always mix and match
// both functions.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanTags adds a new tags to the span. It will appear under "Tags" section
// of the selected span. Use this if you think the tag and its value could be
// useful while debugging.
func AddSpanTags(span trace.Span, tags map[string]string) {
	list := make([]attribute.KeyValue, len(tags))

	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	span.SetAttributes(list...)
}

// AddSpanEvents adds a new events to the span. It will appear under the "Logs"
// section of the selected span. Use this if the event could mean anything
// valuable while debugging.
func AddSpanEvents(span trace.Span, name string, events map[string]string) {
	list := make([]trace.EventOption, len(events))

	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}

	span.AddEvent(name, list...)
}

// AddSpanError adds a new event to the span. It will appear under the "Logs"
// section of the selected span. This is not going to flag the span as "failed".
// Use this if you think you should log any exceptions such as critical, error,
// warning, caution etc. Avoid logging sensitive data!
func AddSpanError(span trace.Span, err error) {
	span.RecordError(err)
}

// FailSpan flags the span as "failed" and adds "error" label on listed trace.
// Use this after calling the `AddSpanError` function so that there is some sort
// of relevant exception logged against it.
func FailSpan(span trace.Span, msg string) {
	span.SetStatus(codes.Error, msg)
}

// SpanCustomiser is used to enforce custom span options. Any custom concrete
// span customiser type must implement this interface.
type SpanCustomiser interface {
	customise() trace.SpanOption
}
