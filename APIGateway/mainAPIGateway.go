package main

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/api"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log"
)

func initTracer() *sdktrace.TracerProvider {
	exporter, _ := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			//otlptracegrpc.WithEndpoint(),
		),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("API-Gateway"),
				semconv.TelemetrySDKLanguageGo,
				semconv.EnduserIDKey.String("lolxD"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := gin.Default()
	r.Use(otelgin.Middleware("API-Gateway"))

	r.GET("/api/getuserinfo/:user", api.GetUserInfo)

	err := r.Run() //"localhost:8080"
	if err != nil {
		log.Panicf("Unable to run gateway API: %s", err)
	}
}
