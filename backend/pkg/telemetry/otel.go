package telemetry

import (
	"context"
	"log"
	"strings"

	"rifa/backend/pkg/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Shutdown func(ctx context.Context) error

// InitCollector sets TracerProvider + MeterProvider with OTLP/HTTP so we can
// send data to a Collector.
func InitCollector(
	ctx context.Context,
	opts config.CollectorOpts,
) (Shutdown, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("rifa-backend"),
			attribute.String("deployment.environment", opts.CollectorEnv),
		),
	)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{}
	if hv := opts.CollectorExporterHeaders; hv != "" {
		headers = ParseKeyValueList(hv)
	}

	traceExp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(opts.CollectorExporter),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExp),
	)

	metricExp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(opts.CollectorExporter),
		otlpmetrichttp.WithHeaders(headers),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)),
	)

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) {
		log.Printf("[otel] %v", e)
	}))

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	return func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return err
		}
		return mp.Shutdown(ctx)
	}, nil
}

// ParseKeyValueList takes a string like "key1=value1,key2=value2"
// and returns a map[string]string{ "key1":"value1", "key2":"value2" }.
// Invalid pairs (missing '=' or empty key/value) are skipped.
func ParseKeyValueList(input string) map[string]string {
	result := make(map[string]string)
	if input == "" {
		return result
	}

	pairs := strings.SplitSeq(input, ",")
	for pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			// skip malformed pair
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key == "" || val == "" {
			continue
		}

		result[key] = val
	}
	return result
}
