package main

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/api"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
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
			otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
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

	cfg := Utils.InitConfig()
	msgHdlFactory := MessageHandler.GetFactory(cfg)

	api.SetSetting(cfg, msgHdlFactory)
	r.GET("/api/getuserinfo/:user", api.GetUserInfo)
	r.POST("/api/createorder", api.CreateOrder)
	r.GET("/api/productdetail/:productname", api.GetProductDetails)
	r.GET("/api/ping", api.Ping)

	err := r.Run() //"localhost:8080"
	if err != nil {
		log.Panicf("Unable to run gateway API: %s", err)
	}
}
