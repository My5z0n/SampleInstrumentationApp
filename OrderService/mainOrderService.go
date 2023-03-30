package main

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/OrderService/api"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
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
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("OrderService"),
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

	var run chan struct{}

	cfg := Utils.InitConfig()
	msgHdlFactory := MessageHandler.GetFactory(cfg)

	//Handlers
	go MessageHandler.MsgRcv(api.CreateOrderHandler, Utils.CreateOrderQueueName, &msgHdlFactory, Utils.All)
	go MessageHandler.MsgRcv(api.ProcessOrderHandler, Utils.ProcessConfirmedOrderQueueName, &msgHdlFactory, Utils.All)
	go MessageHandler.MsgRcv(api.ProcessReturnedPaymentHandler, Utils.ProcessReturnedPaymentQueueName, &msgHdlFactory, Utils.All)

	<-run

}
