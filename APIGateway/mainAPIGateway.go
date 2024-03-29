package main

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/api"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	_ "github.com/My5z0n/SampleInstrumentationApp/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
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
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Use(otelgin.Middleware("API-Gateway"))

	cfg := Utils.InitConfig()
	msgHdlFactory := MessageHandler.GetFactory(cfg)

	api.SetSetting(cfg, &msgHdlFactory)
	r.GET("/api/user/:user", api.GetUserInfo)
	r.POST("/api/order", api.CreateOrder)
	r.GET("/api/product/:productname", api.GetProductDetails)
	r.GET("/api/ping", api.Ping)

	go MessageHandler.MsgRcv(nil, Utils.GetProductDetailsResponseQueueName, &msgHdlFactory, Utils.IDOnly)
	go MessageHandler.MsgRcv(nil, Utils.GetUserInfoResponseQueueName, &msgHdlFactory, Utils.IDOnly)

	err := r.Run() //"localhost:8080"
	if err != nil {
		log.Panicf("Unable to run gateway API: %s", err)
	}
}
