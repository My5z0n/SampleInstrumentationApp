package api

import (
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
)

var tracer = otel.Tracer("OrderService")

func MsgRcv(handler func(trace.Span, map[string]any), queueName string) {
	msgHandler := MessageHandler.MessageHandler{}
	msgHandler.Create(queueName)
	inputChan := msgHandler.RegisterConsumer()

	for msg := range inputChan {
		if span, ok := msg["OTELSPAN"].(trace.Span); ok {
			if msgBody, ok := msg["msg"].(map[string]any); ok {
				go handler(span, msgBody)
			}

		}
	}
}
func CreateOrderHandler(span trace.Span, msg map[string]any) {
	//spanCtx := span.SpanContext()
	defer span.End()

	//Logging
	log.Printf("Received a message: %s", msg)
}
