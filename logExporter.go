package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/sdk/trace"
)

type LogExporter struct{}

func NewLogExporter() *LogExporter {
	return &LogExporter{}
}

func (e *LogExporter) ExportSpans(ctx context.Context, ss []trace.ReadOnlySpan) error {
	for _, s := range ss {
		log.Printf("trace:  %+v\n\n", s)
	}
	return nil
}

func (e *LogExporter) MarshalLog() interface{} {
	return nil
}

func (e *LogExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (e *LogExporter) Start(ctx context.Context) error {
	return nil
}
