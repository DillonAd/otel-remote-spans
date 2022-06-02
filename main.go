package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	ctx := context.Background()
	shutdownTracing := initTracing(ctx)
	defer shutdownTracing()

	localSpan(ctx)
	remoteSpan(ctx)
}

func localSpan(ctx context.Context) {
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "local-span")
	defer span.End()
	log.Printf("localSpan - IsRecording(): %v", span.IsRecording())
}

func remoteSpan(ctx context.Context) {
	spanContextConfig := trace.SpanContextConfig{}

	var err error
	if spanContextConfig.TraceID, err = trace.TraceIDFromHex("96495ea803401c716368207eb4344f43"); err != nil {
		log.Printf("error parsing traceId: %v", err)
		return
	}

	if spanContextConfig.SpanID, err = trace.SpanIDFromHex("d101d692a980bc78"); err != nil {
		log.Printf("error parsing spanId: %v", err)
		return
	}

	spanContext := trace.NewSpanContext(spanContextConfig)

	if !spanContext.IsValid() {
		log.Printf("span context not valid: %+v", spanContext)
	}

	remoteCtx := trace.ContextWithRemoteSpanContext(ctx, spanContext)
	_, remoteSpan := otel.GetTracerProvider().Tracer("default").Start(remoteCtx, "remote-span")
	defer remoteSpan.End()
	log.Printf("remoteSpan - IsRecording(): %v", remoteSpan.IsRecording())
}

func initTracing(ctx context.Context) func() {

	exporter := NewLogExporter()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		log.Println("shutdown tracing")
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("error stopping tracer: %v", err)
		}
	}
}
